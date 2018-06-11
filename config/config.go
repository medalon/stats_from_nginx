package config

import "github.com/kelseyhightower/envconfig"

// StatsNginxConfig ...
type StatsNginxConfig struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

// GetConfig ...
func GetConfig() (*StatsNginxConfig, error) {
	var c StatsNginxConfig
	err := envconfig.Process("statsnginx", &c)
	return &c, err
}
