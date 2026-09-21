// Команда shortener — сервис сокращения ссылок.
package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"urlshortener/internal/config"
	"urlshortener/internal/httpapi"
	"urlshortener/internal/storage"
)

// TODO: прочитать конфигурацию, собрать зависимости, поднять HTTP-сервер.
func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var level slog.Level

	switch cfg.LOG_LEVEL {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "error":
		level = slog.LevelError
	case "warn":
		level = slog.LevelWarn
	}
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	logger := slog.New(logHandler)
	store := storage.New()
	handler := httpapi.New(store, cfg, logger)
	srv := &http.Server{Addr: cfg.HTTPAddr, Handler: handler}
	err = srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server stopped", "err", err)
		os.Exit(1)
	}
}
