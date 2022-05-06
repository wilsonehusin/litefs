package litefs

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"go.husin.dev/litefs/config"
	"go.husin.dev/litefs/internal/log"
)

type Server struct {
	ctx  context.Context
	s    *http.Server
	done chan error
}

func NewServer(ctx context.Context, cfg *config.Config) (*Server, error) {
	lfs, err := NewFS(cfg)
	if err != nil {
		return nil, fmt.Errorf("initializing FS: %w", err)
	}

	s := &http.Server{
		Addr:         cfg.Address,
		Handler:      lfs,
		ReadTimeout:  cfg.RequestTimeout,
		WriteTimeout: cfg.RequestTimeout,
		ErrorLog:     log.Stdlib("http"),
		BaseContext: func(net.Listener) context.Context {
			reqctx, _ := context.WithTimeout(ctx, cfg.RequestTimeout) //nolint:govet
			return reqctx
		},
	}

	return &Server{
		ctx:  ctx,
		s:    s,
		done: make(chan error, 1),
	}, nil
}

func (s *Server) Run() error {
	// TODO: ListenAndServe
	log.Info().Str("addr", s.s.Addr).Msg("starting server")
	<-s.ctx.Done()
	err := s.ctx.Err()
	s.done <- err

	return err
}

func (s *Server) Done() error {
	log.Debug().Msg("waiting for server shutdown")
	return <-s.done
}
