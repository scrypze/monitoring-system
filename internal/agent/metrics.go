package agent

import (
	"log/slog"
	"monitoring-system/proto/agent"
	"time"

	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type MetricsCollector struct {
	logger *slog.Logger
	config *Config
}

func NewMetricsCollector(logger *slog.Logger, config *Config) *MetricsCollector {
	return &MetricsCollector{
		logger: logger,
		config: config,
	}
}

func (mc *MetricsCollector) Collect() (*agent.MetricsData, error) {
	now := time.Now().Unix()

	metrics := &agent.MetricsData{
		HostId:    mc.config.HostID,
		Timestamp: now,
	}

	if contains(mc.config.EnabledMetrics, "cpu") {
		cpuUsage, err := mc.collectCPU()

		if err != nil {
			mc.logger.Error("failed to collect cpu metrics", "error", err)
		} else {
			metrics.CpuUsage = cpuUsage
		}	
	}

	if contains(mc.config.EnabledMetrics, "memory") {
		memUsage, err := mc.collectMemory()

		if err != nil {
			mc.logger.Error("failed to collect memory metrics", "error", err)
		} else {
			metrics.MemoryUsage = memUsage
		}	
	}
	
	if contains(mc.config.EnabledMetrics, "disk") {
		diskUsage, err := mc.collectDisk()

		if err != nil {
			mc.logger.Error("failed to collect disk metrics", "error", err)
		} else {
			metrics.DiskUsage = diskUsage
		}
	}
	

	netIn, netOut, err := mc.collectNetwork()

	if err != nil {
		mc.logger.Error("failed to collect network metrics", "error", err)
	} else {
		metrics.NetworkIn = netIn
		metrics.NetworkOut = netOut
	}
	
	if contains(mc.config.EnabledMetrics, "uptime") {
		uptime, err := mc.collectUptime()

		if err != nil {
			mc.logger.Error("failed to collect uptime", "error", err)
		} else {
			metrics.Uptime = uptime
		}
	}
	
	return metrics, nil
}

func (mc *MetricsCollector) collectCPU() (float64, error) {
	percentages, err := cpu.Percent(0, false)

	if err != nil {
		return 0, err
	}

	if len(percentages) > 0 {
		return percentages[0], nil
	}

	return 0, nil
}

func (mc *MetricsCollector) collectMemory() (float64, error) {
	memInfo, err := mem.VirtualMemory()

	if err != nil {
		return 0, err
	}

	return memInfo.UsedPercent, nil
}

func (mc *MetricsCollector) collectDisk() (float64, error) {
	diskInfo, err := disk.Usage("/")

	if err != nil {
		return 0, err
	}

	return diskInfo.UsedPercent, nil
}

func (mc *MetricsCollector) collectNetwork() (float64, float64, error) {
	netStats, err := net.IOCounters(false)

	if err != nil {
		return 0, 0, err
	}

	var totalBytesIn, totalBytesOut uint64
	for _, stat := range netStats {
		totalBytesIn += stat.BytesRecv
		totalBytesOut += stat.BytesSent
	}

	return float64(totalBytesIn), float64(totalBytesOut), nil
}

func (mc *MetricsCollector) collectUptime() (int64, error) {
	uptime, err := host.Uptime()

	if err != nil {
		return 0, err
	}

	return int64(uptime), nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}