package agent

import (
	"fmt"
	"log"
	"log/slog"
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