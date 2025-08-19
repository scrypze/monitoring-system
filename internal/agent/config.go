package agent

import (
	"fmt"
	"log"
	"log/slog"
	"monitoring-system/proto/agent"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)


type Config struct {
	HostID string
	MetricsInterval int
	EnabledMetrics []string
	CollectorEndpoint string
}

type ConfigManager struct {
	config *Config
	logger *slog.Logger
}

func NewConfigManager(logger *slog.Logger) *ConfigManager {
	return &ConfigManager{
		logger: logger,
	}
}

func (cm *ConfigManager) LoadConfig() error {
	cm.config = &Config{}

	err := godotenv.Load()

	if err != nil {
		return fmt.Errorf("failed to load .env file: %w", err)
	}

	hostID := mustGetEnv("AGENT_HOST_ID")
	metricsInterval := mustGetEnv("AGENT_METRICS_INTERVAL")
	enabledMetrics := mustGetEnv("AGENT_ENABLED_METRICS")
	collectorEndpoint := mustGetEnv("AGENT_COLLECTOR_ENDPOINT")

	partsEnabledMetrics := strings.Split(enabledMetrics, ",")
	
	metricsIntervalInt, err := strconv.Atoi(metricsInterval)
	if err != nil {
		return fmt.Errorf("invalid AGENT_METRICS_INTERVAL: %w", err)
	}
	
	if metricsIntervalInt <= 0 {
		return fmt.Errorf("AGENT_METRICS_INTERVAL must positive")
	}

	cm.config.HostID = hostID
	cm.config.MetricsInterval = metricsIntervalInt
	cm.config.EnabledMetrics = partsEnabledMetrics
	cm.config.CollectorEndpoint = collectorEndpoint

	err = cm.validateConfig(cm.config)

	if err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	return nil
}

func (cm *ConfigManager) GetConfig() *Config {
	return cm.config
}

func getEnv(key, defualtValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defualtValue
	}

	return value
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)

	if value == "" {
		log.Fatal("Variable %s not set", key)
	}

	return value
}

func (cm *ConfigManager) UpdateConfig(newConfig *agent.AgentConfig) error {
	if cm.config == nil {
		return fmt.Errorf("config not loaded") 
	}

	updatedConfig := &Config{
		HostID: cm.config.HostID,
		MetricsInterval: cm.config.MetricsInterval,
		EnabledMetrics:    make([]string, len(cm.config.EnabledMetrics)),
        CollectorEndpoint: cm.config.CollectorEndpoint,
	}
	copy(updatedConfig.EnabledMetrics, cm.config.EnabledMetrics)

	if newConfig.HostId != "" {
        updatedConfig.HostID = newConfig.HostId
    }
    if newConfig.MetricsInterval > 0 {
        updatedConfig.MetricsInterval = int(newConfig.MetricsInterval)
    }
    if len(newConfig.EnabledMetrics) > 0 {
        updatedConfig.EnabledMetrics = newConfig.EnabledMetrics
    }
    if newConfig.CollectorEndpoint != "" {
        updatedConfig.CollectorEndpoint = newConfig.CollectorEndpoint
    }	

	err := cm.validateConfig(updatedConfig)

	if err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	cm.config = updatedConfig

	cm.logger.Info("config updated",
		"host_id", cm.config.HostID,
		"metrics_interval", cm.config.MetricsInterval,
        "enabled_metrics", cm.config.EnabledMetrics,
	)

	return nil
}

func (cm *ConfigManager) validateConfig(config *Config) error {
	if config.HostID == "" {
		return fmt.Errorf("AGENT_HOST_ID must not be empty")
	}

	if config.MetricsInterval <= 0 {
		return fmt.Errorf("AGENT_METRICS_INTERVAL must be positive")
	}

	if len(config.EnabledMetrics) == 0 {
		return fmt.Errorf("AGENT_ENABLED_METRICS must not be empty")
	} else {
		if !isValidEnableMetrics(config.EnabledMetrics) {
			return fmt.Errorf("In AGENT_ENABLED_METRICS searched invalid elements")
		}
	}

	if config.CollectorEndpoint == "" {
		return fmt.Errorf("AGENT_COLLECTOR_ENDPOINT must not be empty")
	}

	return nil
}

func isValidEnableMetrics(elements []string) bool {
	allowed := map[string]bool{
		"cpu":     true,
		"memory":  true,
		"disk":    true,
		"network": true,
		"uptime":  true,
	}

	for _, element := range elements {
		if !allowed[element] {
			return false
		}
	}
	return true
}