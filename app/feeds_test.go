package app_test

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/mmcdole/gofeed"
	"github.com/possiblesloth/packrat/app"
)

func TestFeedUpdate(t *testing.T) {
	path := filepath.Join("testdata", "feed.xml")
	file, err := os.Open(path)
	if err != nil {
		t.Error("failed to open file")
	}
	defer file.Close()
	fp := gofeed.NewParser()
	feed, err := fp.Parse(file)
	if err != nil {
		t.Error("failed to parse feed file")
	}

	localDir := t.TempDir()

	jobs := make(chan app.Job, 5)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for j := range jobs {
			if strings.Contains(j.Filename, "/") {
				t.Errorf("filename contained slash character: %s", j.Filename)
			}
		}
		wg.Done()
	}()

	app.UpdateFeed(feed, localDir, "tombstone", jobs)
	close(jobs)
	wg.Wait()
}
