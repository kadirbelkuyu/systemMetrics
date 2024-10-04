package domain

import "time"

type SystemMetrics struct {
	CreatedAt   time.Time `json:"created_at"`
	CPUUsage    float64   `json:"cpu_usage"`
	MemoryUsage float64   `json:"memory_usage"`
	DiskUsage   float64   `json:"disk_usage"`
}
