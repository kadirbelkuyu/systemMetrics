package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"systemMetric/config"
	"systemMetric/internal/domain"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(cfg config.DatabaseConfig) (*PostgresRepository, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Dbname)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) SaveLog(cpu, memory, disk float64) error {
	query := `INSERT INTO system_logs (cpu_usage, memory_usage, disk_usage) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, cpu, memory, disk)
	return err
}

func (r *PostgresRepository) GetAllLogs() ([]domain.SystemMetrics, error) {
	query := `SELECT cpu_usage, memory_usage, disk_usage FROM system_logs`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []domain.SystemMetrics
	for rows.Next() {
		var log domain.SystemMetrics
		if err := rows.Scan(&log.CPUUsage, &log.MemoryUsage, &log.DiskUsage); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}
