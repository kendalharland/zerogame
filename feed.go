package zerogame

import (
	"context"
	"net/http"
	"time"
)

// Feed describes an archive to install.
type Feed struct {
	// Name is the display name of this Feed.
	//
	// Required.
	Name string `json:"name"`

	// Version is this feed's version string.
	//
	// This is used to distinguish between feeds with the same name.
	//
	// Required.
	Version string `json:"version"`

	// ArchiveURL is the remote location of this feed's archive.
	//
	// Required.
	ArchiveURL string `json:"archive_url"`

	// IsArchiveSigned indicates whether the archive has been signed with a PGP key.
	IsArchiveSigned bool `json:"is_archive_signed"`
}

func InstallFeed(ctx context.Context, url string, repo *FeedRepository) error {
	client := &http.Client{Timeout: 5 * time.Second}
	if err := repo.ImportFeed(ctx, url, client); err != nil {
		return err
	}
	record, _ := repo.GetFeed(url)
	return InstallArchive(record.FeedArchivePath)
}

func UninstallFeed(_ context.Context, feedURL string) error {
	return nil
}

func UpdateFeed(_ context.Context, feedURL string) error {
	return nil
}

func RunFeed(_ context.Context, feedURL string) error {
	return nil
}

func RepairFeed(_ context.Context, feedURL string) error {
	return nil
}
