package app

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func handleUpdateFeed(logger *slog.Logger, config Config, jobs chan<- Job) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			body := struct {
				URL string
			}{}
			err := json.NewDecoder(r.Body).Decode(&body)
			if err != nil {
				logger.Error("unable to decode request", "error", err.Error())
			}

			// TODO make sure URL is valid
			logger.Info("got request to update URL", "URL", body.URL)
			updateFeedFromURL(body.URL, config.StaticDir, config.ServerHost, jobs)
			w.WriteHeader(http.StatusAccepted)
		},
	)
}

func handleHealthzPlease(logger *slog.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			logger.Info("got request for healthcheck")
			w.Write([]byte("OK"))
		},
	)
}
