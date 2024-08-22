package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang-island/internal/api"
	"golang-island/internal/api/middleware/logger"
	"golang-island/internal/config"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.MustLoad()

	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	log.Info(
		"starting image-converter",
		slog.String("version", "123"),
	)
	log.Debug("debug messages are enabled")

	log.Info("starting server", slog.String("address", cfg.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(logger.New(log))

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	controller := api.NewController(log)

	router.Post("/algorithm/island", func(w http.ResponseWriter, r *http.Request) {
		controller.CalculateIsland(w, r)
	})

	router.Get("/algorithm/island/{id}", func(w http.ResponseWriter, r *http.Request) {
		controller.GetIslandResult(w, r)
	})

	router.Get("/algorithm/island/ready/{id}", func(w http.ResponseWriter, r *http.Request) {
		controller.IsReady(w, r)
	})

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", err)

		return
	}

	log.Info("server stopped")
}
