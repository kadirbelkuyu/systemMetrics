package service

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"systemMetric/internal/domain"
)

type MetricsService interface {
	GetMetrics() (*domain.SystemMetrics, error)
}

type metricsService struct{}

func NewMetricsService() MetricsService {
	return &metricsService{}
}

func (s *metricsService) GetMetrics() (*domain.SystemMetrics, error) {
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}

	memStats, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	diskStats, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}

	return &domain.SystemMetrics{
		CPUUsage:    cpuPercent[0],
		MemoryUsage: memStats.UsedPercent,
		DiskUsage:   diskStats.UsedPercent,
	}, nil
}
