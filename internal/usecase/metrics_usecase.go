package usecase

import (
	stdlog "log"
	"systemMetric/internal/service"
	"systemMetric/pkg/logger"
	"time"
)

type MetricsUseCase interface {
	StartMetricsLogging() error
}

type metricsUseCase struct {
	metricsService service.MetricsService
	logger         logger.Logger
}

func NewMetricsUseCase(ms service.MetricsService, logger logger.Logger) MetricsUseCase {
	return &metricsUseCase{
		metricsService: ms,
		logger:         logger,
	}
}

func (uc *metricsUseCase) StartMetricsLogging() error {
	for {
		metrics, err := uc.metricsService.GetMetrics()
		if err != nil {
			return err
		}

		stdlog.Printf("CPU: %.2f%%, Memory: %.2f%%, Disk: %.2f%%\n", metrics.CPUUsage, metrics.MemoryUsage, metrics.DiskUsage)

		if err := uc.logger.LogMetrics(metrics); err != nil {
			return err
		}

		time.Sleep(10 * time.Second)
	}
}
