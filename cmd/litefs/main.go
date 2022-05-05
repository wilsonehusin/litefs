package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"go.husin.dev/litefs"
	"go.husin.dev/litefs/config"
	"go.husin.dev/litefs/internal/log"
)

var stderr = os.Stderr

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(stderr, "\nerror: %s\n", err.Error())
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg, err := config.ParseEnv()
	if err != nil {
		return fmt.Errorf("parsing config from environment: %w", err)
	}

	if cfg.Dev {
		log.Human(stderr)
	} else {
		log.Output(stderr)
	}

	s, err := litefs.NewServer(ctx, cfg)
	if err != nil {
		return fmt.Errorf("booting server: %w", err)
	}

	defer func() {
		log.Err(s.Done()).Msg("shutdown")
	}()
	return s.Run()
}
