package app_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mmcdole/gofeed"
	"github.com/possiblesloth/packrat/app"
)

func TestFeedUpdate(t *testing.T) {
	path := filepath.Join("testdata", "testfeed.xml")
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

	go func() {
		for j := range jobs {
			t.Log(j.Url)
		}
	}()

	app.UpdateFeed(feed, localDir, "tombstone", jobs)

}
