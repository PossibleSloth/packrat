package app_test

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/possiblesloth/packrat/app"
)

func TestServer(t *testing.T) {
	staticDir := t.TempDir()
	baseUrl := "http://127.0.0.1:7777"
	getenv := func(key string) string {
		switch key {
		case "LISTEN_HOST":
			return "127.0.0.1"
		case "LISTEN_PORT":
			return "7777"
		case "STATIC_DIR":
			return staticDir
		case "SERVER_HOST":
			return "127.0.0.1:7777"
		default:
			t.Error("unexpected environment variable")
			return ""
		}
	}
	ctx := context.Background()
	go app.Run(ctx, os.Stdout, getenv)
	err := waitForReady(ctx, 10*time.Second, baseUrl+"/healthz")
	if err != nil {
		t.Errorf("service failed to start: %s", err.Error())
	}

	// Service is finished initializing, functional tests can go below
	resp, err := http.Get(baseUrl + "/api/status")
	if err != nil {
		t.Errorf("error calling status endpoint: %s", err.Error())
	}

	r := struct {
		DownloadQueueLen int `json:"download_queue_len"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Errorf("unable to decode response: %s", err.Error())
	}

	if r.DownloadQueueLen != 0 {
		t.Error("unexpected value from status endpoint")
	}

}
