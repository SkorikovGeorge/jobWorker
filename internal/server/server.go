package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/SkorikovGeorge/jobWorker/internal/router"
	"github.com/rs/zerolog/log"
)

func New() *http.Server {
	log.Info().Msg("configuring new server...")

	mx := router.NewRouter()
	log.Info().Msg("setting up routes...")
	router.SetupRoutes(mx)

	return &http.Server{
		Addr:         fmt.Sprintf("%s:%d", Cfg.Address, Cfg.Port),
		ReadTimeout:  Cfg.ReadTimeout,
		WriteTimeout: Cfg.WriteTimeout,
		IdleTimeout:  Cfg.IdleTimeout,
		Handler:      mx,
	}
}

func Run(server *http.Server) error {
	log.Info().Msgf("listening on port %s", strings.Split(server.Addr, ":")[1])
	return server.ListenAndServe()
}
