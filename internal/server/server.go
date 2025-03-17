package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type ServerConfig struct {
	Address         string        `yaml:"address" mapstructure:"address"`
	Port            int           `yaml:"port" mapstructure:"port"`
	ReadTimeout     time.Duration `yaml:"read_timeout" mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout" mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" mapstructure:"shutdown_timeout"`
	IdleTimeout     time.Duration `yaml:"idle_timeout" mapstructure:"idle_timeout"`
}

func New(cfg *ServerConfig) *http.Server {
	log.Info().Msg("creating new instance...")
	return &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Address, cfg.Port),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		// Handler: ,
	}
}

func Run(server *http.Server) error {
	log.Info().Msg("starting server...")
	log.Info().Msgf("listening on port %s...", strings.Split(server.Addr, ":")[1])
	return server.ListenAndServe()
}
