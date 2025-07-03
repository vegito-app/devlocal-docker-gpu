package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog/log"
)

func ListenAndServe(ctx context.Context, addr string, handler http.Handler) error {
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// Handle graceful shutdown
	waitShutdown := make(chan struct{})
	shutDownCtx, cancelWaitShutdown := context.WithCancel(context.Background())
	defer func() {
		cancelWaitShutdown()
		<-waitShutdown
	}()
	var serverError error
	go func() {
		defer close(waitShutdown)

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		select {
		case <-ctx.Done():
			log.Info().Msg("Context done, shutting down HTTP server")
			shutdown(ctx, server)

		case <-c:
			log.Info().Msg("Received OS interruption signal, shutting down HTTP server")
			shutdown(ctx, server)

		case <-shutDownCtx.Done():
			log.Error().Err(serverError).Msg("HTTP server has been stopped internally")

		}
	}()

	serverError = server.ListenAndServe()
	return serverError
}

func shutdown(ctx context.Context, server *http.Server) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error().Msg("HTTP server shutdown, error: " + err.Error())
		return
	}
	log.Info().Msg("HTTP server has gracefully shutdown")
}
