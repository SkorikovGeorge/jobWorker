package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/SkorikovGeorge/jobWorker/internal/server"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var cfg = server.ServerConfig{
	Address:         "localhost",
	Port:            8080,
	ReadTimeout:     5,
	WriteTimeout:    5,
	ShutdownTimeout: 30,
	IdleTimeout:     60,
}

func main() {
	var err error
	srv := server.New(&cfg)
	log.Info().Msg("Starting server")

	go func() {
		if err = server.Run(srv); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Msg("cant run")
			// log.Fatal().Err(errors.Wrap(err, errs.ErrStartServer)).Msg(errors.Wrap(err, errs.ErrStartServer).Error())
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Info().Msg("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		log.Fatal().Msg("wrong wrong....")
		// log.Fatal().Err(errors.Wrap(err, errs.ErrShutdown)).Msg(errors.Wrap(err, errs.ErrShutdown).Error())
	}
	log.Info().Msg("Server is shut down gracefully")
}
