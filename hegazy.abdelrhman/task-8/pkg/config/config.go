package config

// Config structure for YAML configuration
type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

var cfg Config

// Get returns the loaded configuration
func Get() Config {
	return cfg
}