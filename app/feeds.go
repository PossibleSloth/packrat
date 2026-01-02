package app

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/gopherlibs/feedhub/feedhub"
	"github.com/mmcdole/gofeed"
)

func updateFeedFromURL(url string, staticDir string, host string, jobs chan<- Job) error {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return err
	}

	err = UpdateFeed(feed, staticDir, host, jobs)
	if err != nil {
		return err
	}

	return nil
}

func UpdateFeed(feed *gofeed.Feed, staticDir string, serverHost string, jobs chan<- Job) error {
	// TODO add any other fields that need to be in finished RSS feed
	localfeed := &feedhub.Feed{
		Title:       feed.Title,
		Link:        &feedhub.Link{Href: feed.FeedLink},
		Description: feed.Description,
		Copyright:   feed.Copyright,
		Created:     *feed.PublishedParsed,
		Image: &feedhub.Image{
			Url:   feed.Image.URL, // TODO save feed images locally
			Title: feed.Image.Title,
		},
		Items: make([]*feedhub.Item, len(feed.Items)),
	}

	feedSlug := titleToSlug(feed.Title)
	feedDir := path.Join(staticDir, feedSlug)
	os.MkdirAll(feedDir, os.ModePerm)

	// Add all items to the new feed
	for i, item := range feed.Items {
		if len(item.Enclosures) != 1 {
			return fmt.Errorf("unexpected value for enclosures when updating %s (%d)", feedSlug, len(item.Enclosures))
		}

		// Use hash of GUID as filename. GUID is guaranteed to be unique, but could contain
		// characters we don't want in filename, like slashes.
		hashBytes := sha256.Sum256([]byte(item.GUID))
		localFilename := hex.EncodeToString(hashBytes[:]) + ".mp3"

		// Tell workers to download file
		jobs <- Job{
			Url:      item.Enclosures[0].URL,
			Filename: localFilename,
			LocalDir: feedDir,
		}

		fileUrl := fmt.Sprintf(
			"http://%s/feeds/%s/%s",
			serverHost,
			feedSlug,
			localFilename,
		)

		localfeed.Items[i] = &feedhub.Item{
			Id:          item.GUID,
			Title:       item.Title,
			Description: item.Description,
			Created:     *item.PublishedParsed,
			Link:        &feedhub.Link{Href: item.Link},
			Enclosure: &feedhub.Enclosure{
				Url:    fileUrl,
				Length: item.Enclosures[0].Length,
				Type:   item.Enclosures[0].Type,
			},
		}
	}

	// feedXml, err := localfeed.ToAtom()
	feedXml, err := localfeed.ToRss()
	feedFile := path.Join(feedDir, "feed.xml")
	err = os.WriteFile(feedFile, []byte(feedXml), 0644)
	if err != nil {
		return err
	}

	return nil
}

func titleToSlug(title string) string {
	reg := regexp.MustCompile(`[^\p{L}\p{N}\s-]+`)

	// Replace unwanted characters with an empty string
	cleaned := reg.ReplaceAllString(title, "")

	// Optional: replace spaces with hyphens and convert to lowercase for a cleaner slug
	slug := strings.ToLower(strings.ReplaceAll(cleaned, " ", "-"))

	return slug
}
