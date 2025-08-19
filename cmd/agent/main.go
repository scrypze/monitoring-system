package main

import (
	"log"
	"log/slog"
	"os"
	"time"

	"monitoring-system/internal/agent"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	
	configManager := agent.NewConfigManager(logger)
	
	err := configManager.LoadConfig()

	if err != nil {
		log.Fatal("Failed load config: ", err)
	}
	
	config := configManager.GetConfig()
	
	collector := agent.NewMetricsCollector(logger, config)
	
	logger.Info("starting metrics collection", 
		"host_id", config.HostID,
		"interval", config.MetricsInterval,
	)
	
	for {
		metrics, err := collector.Collect()
		if err != nil {
			logger.Error("failed to collect metrics", "error", err)
		} else {
			logger.Info("collected metrics",
				"cpu", metrics.CpuUsage,
				"memory", metrics.MemoryUsage,
				"disk", metrics.DiskUsage,
				"uptime", metrics.Uptime,
			)
		}
		
		time.Sleep(time.Duration(config.MetricsInterval) * time.Second)
	}
}