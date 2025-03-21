package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	internalErr "github.com/SkorikovGeorge/jobWorker/internal/consts"
	"github.com/SkorikovGeorge/jobWorker/internal/server"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func main() {
	var err error
	log.Info().Msg("Server: starting server...")
	srv := server.New()

	go func() {
		if err = srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(errors.Wrap(err, internalErr.ErrStartingServer)).Msg(errors.Wrap(err, internalErr.ErrStartingServer).Error())
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), srv.Config.ShutdownTimeout)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(errors.Wrap(err, internalErr.ErrShutdownServer)).Msg(errors.Wrap(err, internalErr.ErrShutdownServer).Error())
	}
	log.Info().Msg("Server: server has been shut down gracefully")
}
