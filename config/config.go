package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"

	"go.husin.dev/litefs/internal/log"
)

type Config struct {
	// Upstream contains address of instance to replicate from.
	Upstream string

	// BlobPath contains prefix of where to store blobs.
	BlobPath string `required:"true" default:"."`

	// Database contains path to SQLite database.
	Database string `required:"true"`

	// Address contains HTTP address (i.e., host:port) where server listens to.
	Address string `default:":8080"`

	// RequestTimeout contains timeout duration for requests.
	RequestTimeout time.Duration `envconfig:"reqtimeout" default:"10s"`

	// DBTimeout contains timeout diration for database operations.
	DBTimeout time.Duration `default:"5s"`

	// Debug configures debug logging.
	Debug bool

	// Silent silences logging, takes precedence over Debug.
	Silent bool

	// Dev configures dev mode.
	Dev bool
}

const appname = "litefs"

func ParseEnv() (*Config, error) {
	return ParseEnvWithPrefix(appname)
}

func ParseEnvWithPrefix(prefix string) (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process(prefix, cfg); err != nil {
		return nil, fmt.Errorf("parsing environment variables '%s_*': %w", prefix, err)
	}

	if cfg.Silent {
		log.Disable()
	} else if cfg.Debug {
		log.SetDebug()
	}
	return cfg, nil
}
