package config

import (
	"log"
	"sync"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App      AppConfig
	Upstream UpstreamConfig
}

type AppConfig struct {
	AppName         string `envconfig:"APP_NAME" default:"portifyx-api-gateway"`
	AppEnv          string `envconfig:"APP_ENV" default:"local"`
	AppVersion      string `envconfig:"APP_VERSION" default:"1.0.0"`
	AppPort         string `envconfig:"APP_PORT" default:"8080"`
	HTTPTimeoutInMs int64  `envconfig:"HTTP_TIMEOUT_IN_MS" default:"30000"`
	TimeZone        string `envconfig:"TIMEZONE" default:"Asia/Bangkok"`
}

type UpstreamConfig struct {
	PortfolioServiceURL string `envconfig:"UPSTREAM_PORTFOLIO_URL" default:"http://localhost:8081"`
	ExchangeServiceURL  string `envconfig:"UPSTREAM_EXCHANGE_URL" default:"http://localhost:8082"`
	MarketDataServiceURL string `envconfig:"UPSTREAM_MARKETDATA_URL" default:"http://localhost:8083"`
}

var (
	cfg  Config
	once sync.Once
)

func Init() {
	once.Do(func() {
		_ = godotenv.Load() // optional .env load, falls back to environment variables
		if err := envconfig.Process("", &cfg); err != nil {
			log.Fatalf("failed to parse environment config: %v", err)
		}
	})
}

func GetConfig() *Config {
	return &cfg
}
