package zerogame

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

const markerName = "marker.zg"

type cache string

func newDefaultCache() (cache, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get the user's home directory: %w", err)
	}
	return cache(filepath.Join(home, ".zerogame")), nil
}

func (c cache) FeedArchiveExists(feedURL string) bool {
	_, err := c.GetFeedArchive(feedURL)
	return err == nil
}

func (c cache) GetFeedArchive(feedURL string) (string, error) {
	marker := filepath.Join(c.feedDir(feedURL), markerName)
	if _, err := os.Stat(marker); os.IsNotExist(err) {
		return "", errors.New("not found")
	}
	basename, err := ioutil.ReadFile(marker)
	if err != nil {
		return "", errors.New("failed to read marker")
	}
	filename := filepath.Join(c.feedDir(feedURL), string(basename))
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return "", errors.New("not found")
	}
	return filename, nil
}

func (c cache) WriteFeedArchive(feedURL string, feed *Feed, archive []byte, mode os.FileMode) error {
	ensureDir(c.feedDir(feedURL))
	basename := fmt.Sprintf("%s-%s.%s", feed.Name, feed.Version, feed.ArchiveType)
	filename := filepath.Join(c.feedDir(feedURL), basename)
	fmt.Fprintf(os.Stderr, "Writing feed archive to %s...\n", filename)
	if err := ioutil.WriteFile(filename, archive, mode); err != nil {
		return err
	}
	marker := filepath.Join(c.feedDir(feedURL), markerName)
	if err := ioutil.WriteFile(marker, []byte(basename), 0755); err != nil {
		return err
	}
	return nil
}

func (c cache) feedDir(feedURL string) string {
	return filepath.Join(string(c), uniqueFeedID(feedURL))
}

func ensureDir(p string) {
	os.MkdirAll(p, os.FileMode(0755))
}

func uniqueFeedID(feedURL string) string {
	id := uuid.NewMD5(uuid.NameSpaceURL, []byte(feedURL))
	return id.String()
}
