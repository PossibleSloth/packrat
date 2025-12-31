package app

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

type Config struct {
	ListenHost string
	Port       string
	StaticDir  string
	ServerHost string
}

func NewServer(
	logger *slog.Logger,
	config *Config,
	jobs chan<- Job,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
		logger,
		*config,
		jobs,
	)
	var handler http.Handler = mux
	return handler
}

// Run creates all dependencies and starts the service
// it is called by main()
func Run(ctx context.Context, w io.Writer, getenv func(key string) string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	config := Config{
		ListenHost: getenv("LISTEN_HOST"),
		Port:       getenv("LISTEN_PORT"),
		StaticDir:  getenv("STATIC_DIR"),
		ServerHost: getenv("SERVER_HOST"),
	}

	logger := slog.New(slog.NewJSONHandler(w, nil))
	jobs := make(chan Job, 1000)

	srv := NewServer(
		logger,
		&config,
		jobs,
	)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.ListenHost, config.Port),
		Handler: srv,
	}
	go func() {
		logger.Info("listening", "addr", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	// Create worker threads to download files
	for range 10 {
		go worker(jobs, logger)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()
	return nil
}
