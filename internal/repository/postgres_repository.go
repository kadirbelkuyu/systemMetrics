package repository

import (
	"database/sql"
	"log"
	"systemMetric/internal/domain"
)

type PostgresRepository struct {
	DB *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{DB: db}
}

func (r *PostgresRepository) SaveLog(cpuUsage, memoryUsage, diskUsage float64) error {
	query := `INSERT INTO logs (cpu_usage, memory_usage, disk_usage) VALUES ($1, $2, $3)`
	_, err := r.DB.Exec(query, cpuUsage, memoryUsage, diskUsage)
	if err != nil {
		log.Printf("Error saving log: %v", err)
		return err
	}
	return nil
}

func (r *PostgresRepository) GetAllLogs() ([]domain.SystemMetrics, error) {
	query := `SELECT created_at, cpu_usage, memory_usage, disk_usage FROM logs ORDER BY created_at DESC`
	rows, err := r.DB.Query(query)
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
			return nil, err
		}
		logs = append(logs, metric)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, err
	}

	return logs, nil
}
