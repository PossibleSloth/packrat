package app

import (
	"log/slog"
	"net/http"
)

func addRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
	config Config,
	jobs chan<- Job,
) {
	fs := http.FileServer(http.Dir(config.StaticDir))
	mux.Handle("/feeds/", http.StripPrefix("/feeds/", fs))
	mux.Handle("/api/feeds", handleUpdateFeed(logger, config, jobs))
	mux.Handle("/api/status", handleGetStatus(logger, jobs))
	mux.Handle("/healthz", handleHealthzPlease(logger))
}
