package logger

import (
	"fmt"
	"log"
	"systemMetric/internal/domain"
	"systemMetric/internal/repository"
)

type Logger interface {
	LogMetrics(metrics *domain.SystemMetrics) error
	GetAllLogs() ([]domain.SystemMetrics, error)
	GetLogs(offset, limit int) ([]domain.SystemMetrics, error)
}

type logger struct {
	repo *repository.PostgresRepository
}

func LogSystemMetrics(metrics domain.SystemMetrics) {
	fmt.Printf("System Metrics: %+v\n", metrics)
}

func NewLogger(repo *repository.PostgresRepository) Logger {
	return &logger{repo: repo}
}

func (l *logger) LogMetrics(metrics *domain.SystemMetrics) error {
	return l.repo.SaveLog(metrics.CPUUsage, metrics.MemoryUsage, metrics.DiskUsage)
}

func (l *logger) GetAllLogs() ([]domain.SystemMetrics, error) {
	return l.repo.GetAllLogs()
}

func (l *logger) GetLogs(offset, limit int) ([]domain.SystemMetrics, error) {
	query := `SELECT created_at, cpu_usage, memory_usage, disk_usage FROM logs ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := l.repo.DB.Query(query, limit, offset)
	if err != nil {
		log.Printf("Error querying logs: %v", err)
		return nil, err
	}
	defer rows.Close()

	var logs []domain.SystemMetrics
	for rows.Next() {
		var metric domain.SystemMetrics
		if err := rows.Scan(&metric.CreatedAt, &metric.CPUUsage, &metric.MemoryUsage, &metric.DiskUsage); err != nil {
			log.Printf("Error scanning log: %v", err)
			fmt.Printf("Scanned metric: %+v\n", metric)
			return nil, err
		}
		logs = append(logs, metric)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, err
	}

	log.Printf("Fetched logs from DB: %v", logs)
	return logs, nil
}
