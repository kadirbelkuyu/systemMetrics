package main

import (
	"log"
	"systemMetric/config"
	"systemMetric/internal/repository"
	"systemMetric/internal/service"
	"systemMetric/internal/usecase"
	"systemMetric/pkg/logger"
)

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
	useCase := usecase.NewMetricsUseCase(service, logger)

	if err := useCase.StartMetricsLogging(); err != nil {
		log.Fatalf("Metrikler loglanırken hata oluştu: %v", err)
	}
}
