package main

import (
	"RestApiFP/1thRestApiPr/internal/config"
	"RestApiFP/1thRestApiPr/internal/http-server/handlers/redirect"
	"RestApiFP/1thRestApiPr/internal/http-server/handlers/url/save"
	mwLogger "RestApiFP/1thRestApiPr/internal/http-server/middleware/logger"
	"RestApiFP/1thRestApiPr/internal/lib/logger/sl"
	"RestApiFP/1thRestApiPr/internal/storage/sqlite"
	"net/http"

	"fmt"
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev = "dev"
	envProd = "prod"
)

func setupLogger(env string) *slog.Logger{

	// err := godotenv.Load()

	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg) // necessary delete

	log := setupLogger(cfg.Env)

	log.Info("starting url-shotener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")
	
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()
	// middleware
	router.Use(middleware.RequestID)
	// two aplicable Loggers
	// to Log a requests is a good practice
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer) // If panic is happen
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(log, storage))
	router.Get("/{alias}", redirect.New(log, storage))
	// TODO write a delete method


	log.Info("staring server", slog.String("address", cfg.Server.Address))
	
	srv := &http.Server{
		Addr: cfg.Server.Address,
		Handler: router,
		ReadTimeout: cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout: cfg.Server.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
	
	// _ = storage

	// middleware

	// TODO: run server

}