package server

import (
	"context"
	"github.com/ChristinaFomenko/gophermart/config"
	"net/http"
	"time"
)

const (
	idleTimeout  = 60 * time.Second
	readTimeout  = 60 * time.Second
	writeTimeout = 60 * time.Second
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         cfg.RunAddress,
			Handler:      handler,
			IdleTimeout:  idleTimeout,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
