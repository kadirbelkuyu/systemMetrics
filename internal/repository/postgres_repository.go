package repository

import (
	"database/sql"
	_ "github.com/lib/pq"
	"systemMetric/internal/domain"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) SaveLog(cpu, memory, disk float64) error {
	query := `INSERT INTO system_logs (cpu_usage, memory_usage, disk_usage) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, cpu, memory, disk)
	return err
}

func (r *PostgresRepository) GetAllLogs() ([]domain.SystemMetrics, error) {
	query := `SELECT cpu_usage, memory_usage, disk_usage, created_at FROM system_logs`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []domain.SystemMetrics
	for rows.Next() {
		var log domain.SystemMetrics
		if err := rows.Scan(&log.CPUUsage, &log.MemoryUsage, &log.DiskUsage, &log.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}
