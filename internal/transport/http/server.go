package http

import (
	"context"
	"github.com/zhayt/clean-arch-tmp-forum/config"
	"github.com/zhayt/clean-arch-tmp-forum/internal/transport/http/handler"
	"github.com/zhayt/clean-arch-tmp-forum/internal/transport/http/middleware"
	"github.com/zhayt/clean-arch-tmp-forum/logger"
	"net"
	"net/http"
	"time"
)

const (
	_defaultShutdownTimeout = 3 * time.Second
)

type Server struct {
	httpServer      *http.Server
	handler         *handler.Handler
	mid             *middleware.Middleware
	Notify          chan error
	shutdownTimeOut time.Duration
}

func NewServer(cfg *config.Config, l *logger.Logger, handler *handler.Handler, mid *middleware.Middleware) *Server {
	srv := &http.Server{
		Addr:              net.JoinHostPort("", cfg.Server.Port),
		MaxHeaderBytes:    1 << 20,
		ErrorLog:          l.Error,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
	}

	return &Server{
		httpServer:      srv,
		handler:         handler,
		mid:             mid,
		Notify:          make(chan error, 1),
		shutdownTimeOut: _defaultShutdownTimeout,
	}
}

func (s *Server) Start() {
	s.httpServer.Handler = s.InitRoutes()
	go func() {
		s.Notify <- s.httpServer.ListenAndServe()
	}()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeOut)
	defer cancel()

	return s.httpServer.Shutdown(ctx)
}
