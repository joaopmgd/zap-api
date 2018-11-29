package config

import (
	"os"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

// Config will setup the Endpoints, the sources that will be requested, Log and Memory Cache configuration
type Config struct {
	Endpoints   *Endpoint
	Datasources *map[string]bool
	Cache       *cache.Cache
	Logger      *logrus.Logger
}

// Endpoints for the future Requests
type Endpoint struct {
	ZapProperties string
}

func GetConfig() *Config {
	return &Config{
		Endpoints: &Endpoint{
			ZapProperties: os.Getenv("ZAP_PROPERTIES_ENDPOINT"),
		},
		Datasources: &map[string]bool{
			"zap":      true,
			"vivareal": true,
		},
		Cache:  cache.New(10*time.Minute, 60*time.Minute),
		Logger: logrus.New(),
	}
}
