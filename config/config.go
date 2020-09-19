package config

import (
	"time"

	"github.com/joeshaw/envdecode"
)

// Config stores all configuration
type Config struct {
	HTTP  HTTPConfig
	Mongo MongoConfig
	App   AppConfig
}

// HTTPConfig stores HTTP configuration
type HTTPConfig struct {
	Hostname string `env:"HTTP_HOSTNAME"`
	Port     uint16 `env:"HTTP_PORT,default=5000"`
	Auth     HTTPAuthConfig
}

// HTTPAuthConfig stores HTTP authentication configuration
type HTTPAuthConfig struct {
	User     string `env:"HTTP_AUTH_USER"`
	Password string `env:"HTTP_AUTH_PASSWORD"`
}

// MongoConfig stores Mongo DB configuration
type MongoConfig struct {
	URL     string        `env:"MONGO_URL,default=mongodb://localhost:27017"`
	DB      string        `env:"MONGO_DB,default=scanner"`
	Timeout time.Duration `env:"MONGO_TIMEOUT,default=5s"`
}

// AppConfig stores application configuration
type AppConfig struct {
	Name string `env:"APP_NAME,default=OTC Scanner"`
}

// GetConfig loads and returns configuration
func GetConfig() (Config, error) {
	var cfg Config
	err := envdecode.StrictDecode(&cfg)
	return cfg, err
}
