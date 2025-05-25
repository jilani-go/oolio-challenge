package config

import "time"

// ServerConfig holds http server config.
type ServerConfig struct {
	Port         string        `json:"port"`
	ReadTimeout  time.Duration `json:"readTimeout"`
	WriteTimeout time.Duration `json:"writeTimeout"`
	IdleTimeout  time.Duration `json:"idleTimeout"`
}

// Config holds the application's config.
type Config struct {
	Server ServerConfig `json:"server"`
}

// Load creates and returns a new config with default values.
func Load() *Config {
	cfg := &Config{
		Server: ServerConfig{
			Port:         "8080", // Default HTTP port
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
	}
	return cfg
}
