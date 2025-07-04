package server

import (
	"log/slog"
	"net/http"

	"github.com/sig-0/boring-avatars-go/server/config"
)

type (
	Middleware = func(http.Handler) http.Handler
	Option     func(s *Server)
)

// WithLogger specifies the logger for the server
func WithLogger(l *slog.Logger) Option {
	return func(s *Server) {
		s.logger = l
	}
}

// WithConfig specifies the config for the server
func WithConfig(c *config.Config) Option {
	return func(s *Server) {
		s.config = c
	}
}

func WithMiddlewares(mw ...Middleware) Option {
	return func(s *Server) {
		s.middlewares = append(s.middlewares, mw...)
	}
}
