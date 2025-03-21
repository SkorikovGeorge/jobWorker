package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/SkorikovGeorge/jobWorker/internal/router"
	"github.com/SkorikovGeorge/jobWorker/internal/workerpool"
	"github.com/rs/zerolog/log"
)

type Server struct {
	Config     *ServerConfig
	httpServer *http.Server
}

func New() *Server {
	log.Info().Msg("Server: configuring new server...")

	mx := router.NewRouter()
	log.Info().Msg("Server: setting up routes...")
	router.SetupRoutes(mx)

	return &Server{
		Config: &cfg,
		httpServer: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", cfg.Address, cfg.Port),
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
			Handler:      mx,
		},
	}
}

func (srv *Server) Run() error {
	log.Info().Msgf("Server: listening on port %d", srv.Config.Port)
	return srv.httpServer.ListenAndServe()
}

func (srv *Server) Shutdown(ctx context.Context) error {
	log.Info().Msg("Server: server shutting down server")
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := workerpool.Workers.Shutdown(timeoutCtx); err != nil {
		log.Error().Err(fmt.Errorf("Server shutdown: %w", err)).Msgf("Server shutdown: %v", err)
	}
	return srv.httpServer.Shutdown(ctx)
}
