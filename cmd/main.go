package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	internalErr "github.com/SkorikovGeorge/jobWorker/internal/errors"
	"github.com/SkorikovGeorge/jobWorker/internal/server"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func main() {
	var err error
	log.Info().Msg("starting server...")
	srv := server.New()

	go func() {
		if err = server.Run(srv); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(errors.Wrap(err, internalErr.StartingServer)).Msg(errors.Wrap(err, internalErr.StartingServer).Error())
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Info().Msg("server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), server.Cfg.ShutdownTimeout)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(errors.Wrap(err, internalErr.ShutdownServer)).Msg(errors.Wrap(err, internalErr.ShutdownServer).Error())
	}
	log.Info().Msg("server has been shut down gracefully")
}
