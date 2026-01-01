package app_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/possiblesloth/packrat/app"
)

func TestServer(t *testing.T) {
	getenv := func(key string) string {
		switch key {
		case "LISTEN_HOST":
			return "localhost"
		case "LISTEN_PORT":
			return "7777"
		case "STATIC_DIR":
			return "/tmp"
		case "SERVER_HOST":
			return "localhost:7777"
		default:
			t.Error("unexpected environment variable")
			return ""
		}
	}
	ctx := context.Background()
	go app.Run(ctx, os.Stdout, getenv)
	err := waitForReady(ctx, 10*time.Second, "http://localhost:7777/healthz")
	if err != nil {
		t.Errorf("service failed to start: %s", err.Error())
	}
}
