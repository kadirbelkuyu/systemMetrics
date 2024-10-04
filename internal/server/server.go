package server

import (
	"context"
	"database/sql"
	"fmt"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	mwLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"gopkg.in/gomail.v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"systemMetric/config"
	"systemMetric/internal/repository"
	"systemMetric/internal/service"
	"systemMetric/internal/usecase"
	"systemMetric/pkg/logger"
	"time"
)

const (
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

// Server struct
type Server struct {
	app        *fiber.App
	cfg        *config.Config
	db         *sql.DB
	mailDialer *gomail.Dialer
}

// NewServer New Server constructor
func NewServer(cfg *config.Config, db *sql.DB, mailDialer *gomail.Dialer) *Server {
	return &Server{
		app:        fiber.New(),
		cfg:        cfg,
		db:         db,
		mailDialer: mailDialer,
	}
}

func (s *Server) Run() error {
	// Initialize standard Go html template engine
	engine := html.New("./templates", ".html")

	// Initialize Fiber app with the view engine
	s.app = fiber.New(fiber.Config{
		Views: engine,
	})

	ctx, _ := context.WithCancel(context.Background())

	// Init Repositories
	postgresRepo := repository.NewPostgresRepository(s.db)

	// Initialize logger
	logger := logger.NewLogger(postgresRepo)

	// Initialize service and use case
	metricsService := service.NewMetricsService()
	useCase := usecase.NewMetricsUseCase(metricsService, logger)

	// Start logging metrics in a separate goroutine
	go func() {
		if err := useCase.StartMetricsLogging(); err != nil {
			log.Fatalf("Metrikler loglanırken hata oluştu: %v", err)
		}
	}()

	// Use Fiber logger middleware
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
		logs, err := logger.GetAllLogs()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(logs)
	})

	// Swagger route
	s.app.Get("/swagger/*", swagger.HandlerDefault)

	// Define Fiber settings
	server := &http.Server{
		Addr:           s.cfg.Server.Port,
		ReadTimeout:    time.Second * s.cfg.Server.ReadTimeout,
		WriteTimeout:   time.Second * s.cfg.Server.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Server is listening on PORT: %s", s.cfg.Server.Port)
		if err := s.app.Listen(server.Addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error starting Server: %v", err)
		}
	}()

	// Graceful shutdown
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
