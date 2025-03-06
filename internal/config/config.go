package config

import (
	"github.com/kelseyhightower/envconfig"

	"github.com/AddMile/backend/internal/shared"
)

type Config struct {
	Debug           bool               `envconfig:"DEBUG" required:"true"`
	Environment     shared.Environment `envconfig:"ENVIRONMENT" required:"true"`
	APIPort         int                `envconfig:"API_PORT" required:"true"`
	WorkerPort      int                `envconfig:"WORKER_PORT" required:"true"`
	CORSOrigin      string             `envconfig:"CORS_ORIGIN" required:"true"`
	FrontendBaseURL string             `envconfig:"FRONTEND_BASE_URL" required:"true"`
	GoogleProjectID string             `envconfig:"GOOGLE_PROJECT_ID" required:"true"`
	APIKey          string             `envconfig:"API_KEY" required:"true"`
	PostgresDSN     string             `envconfig:"POSTGRES_DSN" required:"true"`

	CustomerIOAPIKey        string `envconfig:"CUSTOMER_IO_API_KEY" required:"true"`
	CustomerIOEndpoint      string `envconfig:"CUSTOMER_IO_ENDPOINT" required:"true"`
	CustomerIOBatchSize     int    `envconfig:"CUSTOMER_IO_BATCH_SIZE" required:"true"`
	CustomerIOFlushInterval int    `envconfig:"CUSTOMER_IO_FLUSH_INTERVAL_MS" required:"true"`
	CustomerIOVerbose       bool   `envconfig:"CUSTOMER_IO_VERBOSE" required:"true"`

	TopicUserCreated string `envconfig:"TOPIC_USER_CREATED" required:"true"`
}

func Load() *Config {
	var cfg Config
	envconfig.MustProcess("", &cfg)

	return &cfg
}
