package app

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path"
)

type Job struct {
	Url      string
	Filename string
	LocalDir string
}

func worker(jobs <-chan Job, logger *slog.Logger) {
	for j := range jobs {
		logger.Info("downloading file", "url", j.Url, "filename", j.Filename, "localdir", j.LocalDir)
		filepath := path.Join(j.LocalDir, j.Filename)

		// Check whether the file already exists
		// If it does we don't need to download it again
		if _, err := os.Stat(filepath); !errors.Is(err, os.ErrNotExist) {
			continue
		}

		// Create local file to write to
		out, err := os.Create(filepath)
		if err != nil {
			logger.Error("error creating file for download", "error", err.Error())
			continue
		}
		defer out.Close()

		resp, err := http.Get(j.Url)
		if err != nil {
			logger.Error("error making request to file URL", "error", err.Error())
			continue
		}
		defer resp.Body.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			logger.Error("error writing to file", "error", err.Error())
			continue
		}
	}
}
