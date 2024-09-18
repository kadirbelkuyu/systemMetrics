package logger

import (
	"systemMetric/internal/domain"
	"systemMetric/internal/repository"
)

type Logger interface {
	LogMetrics(metrics *domain.SystemMetrics) error
	GetAllLogs() ([]domain.SystemMetrics, error)
}

func (l *logger) GetAllLogs() ([]domain.SystemMetrics, error) {
	return l.repo.GetAllLogs()
}

type logger struct {
	repo *repository.PostgresRepository
}

func NewLogger(repo *repository.PostgresRepository) Logger {
	return &logger{repo: repo}
}

func (l *logger) LogMetrics(metrics *domain.SystemMetrics) error {
	return l.repo.SaveLog(metrics.CPUUsage, metrics.MemoryUsage, metrics.DiskUsage)
}
