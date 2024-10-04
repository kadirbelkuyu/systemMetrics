package server

import (
	"context"
	"database/sql"
	_ "encoding/json"
	"fmt"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	_ "strconv"
	"syscall"
	"systemMetric/config"
	"systemMetric/internal/repository"
	"systemMetric/internal/service"
	"systemMetric/internal/usecase"
	"systemMetric/pkg/logger"
	"time"

	"github.com/gofiber/fiber/v2"
	mwLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"gopkg.in/gomail.v2"
)

const (
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

type Server struct {
	app        *fiber.App
	cfg        *config.Config
	db         *sql.DB
	mailDialer *gomail.Dialer
}

func NewServer(cfg *config.Config, db *sql.DB, mailDialer *gomail.Dialer) *Server {
	return &Server{
		app:        fiber.New(),
		cfg:        cfg,
		db:         db,
		mailDialer: mailDialer,
	}
}

func (s *Server) Run() error {
	engine := html.New("./templates", ".html")
	s.app = fiber.New(fiber.Config{
		Views: engine,
	})

	ctx, _ := context.WithCancel(context.Background())

	postgresRepo := repository.NewPostgresRepository(s.db)
	logger := logger.NewLogger(postgresRepo)
	metricsService := service.NewMetricsService()
	useCase := usecase.NewMetricsUseCase(metricsService, logger)

	go func() {
		if err := useCase.StartMetricsLogging(); err != nil {
			log.Fatalf("Metrikler loglanırken hata oluştu: %v", err)
		}
	}()

	s.app.Use(mwLogger.New())

	s.app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", nil)
	})

	s.app.Get("/metrics", func(c *fiber.Ctx) error {
		metrics, err := metricsService.GetMetrics()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(metrics)
	})

	s.app.Get("/logs", func(c *fiber.Ctx) error {
		page := c.QueryInt("page", 1)
		if page < 1 {
			page = 1
		}
		limit := c.QueryInt("limit", 50)
		if limit < 1 {
			limit = 50
		}
		offset := (page - 1) * limit

		logs, err := logger.GetLogs(offset, limit)
		if err != nil {
			log.Printf("Error fetching logs: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		log.Printf("Fetched logs: %v", logs)
		return c.JSON(logs)
	})

	s.app.Get("/swagger/*", swagger.HandlerDefault)

	server := &http.Server{
		Addr:           s.cfg.Server.Port,
		ReadTimeout:    time.Second * s.cfg.Server.ReadTimeout,
		WriteTimeout:   time.Second * s.cfg.Server.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	go func() {
		log.Printf("Server is listening on PORT: %s", s.cfg.Server.Port)
		if err := s.app.Listen(server.Addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error starting Server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		fmt.Errorf("signal.Notify: %v", v)
	case done := <-ctx.Done():
		fmt.Errorf("ctx.Done: %v", done)
	}

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	fmt.Println("Server Exited Properly")
	return s.app.ShutdownWithContext(ctx)
}
