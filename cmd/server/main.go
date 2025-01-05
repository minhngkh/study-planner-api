package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"study-planner-api/internal/api"
	handlerImpl "study-planner-api/internal/handler"
	"study-planner-api/internal/utils"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	// The context is used to inform the server it has 5 seconds to finish the request it is
	// currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Error().
			Err(err).
			Msgf("Server forced to shutdown")
	}

	log.Info().Msg("Server exiting")

	done <- true
}

func setupPrettyZeroLog() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}

func main() {
	utils.LoadEnv()

	setupPrettyZeroLog()

	impl := api.NewStrictHandler(handlerImpl.NewHandler(), nil)

	handler := api.NewEchoHandler()
	api.RegisterHandlers(handler, impl)

	port, _ := strconv.Atoi(os.Getenv("PORT"))

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	done := make(chan bool, 1)
	go gracefulShutdown(s, done)

	log.Info().Msgf("Server starting on \x1b[33mhttp://localhost%s\x1b[0m", s.Addr)

	err := s.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal().
			Err(err).
			Msg("Http server error")
	}

	<-done
}
