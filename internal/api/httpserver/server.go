package httpserver

import (
	"context"
	"intellisys-backend/internal/domain/log"
	"intellisys-backend/internal/infra/config"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	log        log.Logger
}

func NewServer(config *config.Config, handler http.Handler, log log.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         config.Server.Host + ":" + config.Server.Port,
			Handler:      handler,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		log: log,
	}
}

func (s *Server) Start() error {
	s.log.Info("Servidor rodando", "Endereço", s.httpServer.Addr)

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info("Servidor pausado")

	return s.httpServer.Shutdown(ctx)
}
