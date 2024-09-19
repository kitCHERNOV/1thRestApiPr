package main

import (
	"RestApiFP/1thRestApiPr/internal/config"
	"fmt"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev = "dev"
	envProd = "prod"
)

func setupLogger(env string) *slog.Logger{
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

	// TODO: init config: cleanenv

	// TODO: init logger: slog

	// TODO: init storage: sqlite

	// TODO: init router: chi, "chi render"

	// TODO: run server

}