package main

import (
	"fmt"
	"log"
	"os"
	"systemMetric/config"
	"systemMetric/infra/postgresql"
	"systemMetric/internal/server"
	"systemMetric/pkg/mailer"
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
	log.Println("Starting API Server...")

	configPath := config.GetConfigPath(os.Getenv("config"))
	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	psqlDB, err := postgresql.ConnectPostgres(cfg)
	if err != nil {
		log.Fatalf("postgresql init: %s", err)
	} else {
		fmt.Println("Postgres connected")
	}

	mailDialer := mailer.NewMailDialer(cfg)
	fmt.Println("Mail dialer connected")

	s := server.NewServer(cfg, psqlDB, mailDialer)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
