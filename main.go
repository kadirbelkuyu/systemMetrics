package main

import (
	"log"
	"systemMetric/config"
	"systemMetric/internal/repository"
	"systemMetric/internal/service"
	"systemMetric/pkg/logger"

	swagger "github.com/arsmn/fiber-swagger/v2" // fiber-swagger middleware
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	_ "systemMetric/docs" // Import the generated docs
)

// @title System Metrics API
// @version 1.0
// @description This is a sample server for system metrics.
// @host localhost:3000
// @BasePath /

// @Summary Get system metrics
// @Description Get the current system metrics
// @Tags metrics
// @Accept json
// @Produce json
// @Success 200 {object} domain.SystemMetrics
// @Failure 500 {object} fiber.Map
// @Router /metrics [get]

// @Summary Get all system logs
// @Description Get all system logs from the database
// @Tags logs
// @Accept json
// @Produce json
// @Success 200 {array} domain.SystemMetrics
// @Failure 500 {object} fiber.Map
// @Router /logs [get]
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Config yüklenemedi: %v", err)
	}

	db, err := repository.NewPostgresRepository(cfg.Database)
	if err != nil {
		log.Fatalf("Veritabanı bağlantısı kurulamadı: %v", err)
	}

	logger := logger.NewLogger(db)
	service := service.NewMetricsService()

	// Initialize standard Go html template engine
	engine := html.New("./templates", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", nil)
	})

	app.Get("/metrics", func(c *fiber.Ctx) error {
		metrics, err := service.GetMetrics()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(metrics)
	})

	app.Get("/logs", func(c *fiber.Ctx) error {
		logs, err := logger.GetAllLogs()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(logs)
	})

	// Swagger route
	app.Get("/swagger/*", swagger.HandlerDefault)

	log.Fatal(app.Listen(":3000"))
}
