package zerogame

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type FeedRecord struct {
	FeedURL               string `json:"feed_url"`
	FeedArchivePath       string `json:"feed_archive_path"`
	FeedSignedArchivePath string `json:"feed_signed_archive_path"`
}

type FeedRepository struct {
	db     *FeedRecordDB
	tmpdir string
}

func NewFeedRepository(tmpdir string, db *FeedRecordDB) *FeedRepository {
	return &FeedRepository{
		tmpdir: tmpdir,
		db:     db,
	}
}

func (r *FeedRepository) ImportFeed(ctx context.Context, url string, client *http.Client) error {
	feed, err := getFeed(client, url)
	if err != nil {
		return err
	}

	root, err := ioutil.TempDir(r.tmpdir, "")
	if err != nil {
		return err
	}

	archive, err := getURL(client, feed.ArchiveURL)
	if err != nil {
		return err
	}
	ef, err := ioutil.TempFile(root, "archive")
	if err != nil {
		return err
	}
	defer ef.Close()
	if _, err := ef.Write(archive); err != nil {
		return err
	}

	rec := &FeedRecord{
		FeedURL:         url,
		FeedArchivePath: ef.Name(),
	}

	if feed.IsArchiveSigned {
		df, err := ioutil.TempFile(root, "archive.decrypted")
		if err != nil {
			return err
		}
		defer df.Close()
		if err := decrypt(ef.Name(), df.Name()); err != nil {
			return err
		}
		rec.FeedSignedArchivePath = rec.FeedArchivePath
		rec.FeedArchivePath = df.Name()
	}

	return r.db.PutRecord(url, rec)
}

func (r *FeedRepository) RemoveFeed(url string) error {
	return r.db.RemoveRecord(url)
}

func (r *FeedRepository) GetFeed(url string) (*FeedRecord, error) {
	return r.db.GetRecord(url)
}

func (r *FeedRepository) Feeds() ([]*FeedRecord, error) {
	return r.db.Records()
}

type FeedRecordDB struct {
	root string
}

func NewFeedRecordDB(root string) *FeedRecordDB {
	return &FeedRecordDB{root: root}
}

func (db *FeedRecordDB) GetRecord(url string) (*FeedRecord, error) {
	basename := uuid.NewMD5(uuid.NameSpaceURL, []byte(url)).String()
	filename := filepath.Join(db.root, basename)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	rec := &FeedRecord{}
	if err := json.Unmarshal(content, rec); err != nil {
		return nil, err
	}
	return rec, nil
}

func (db *FeedRecordDB) PutRecord(url string, rec *FeedRecord) error {
	ensureDir(db.root)
	content, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	basename := uuid.NewMD5(uuid.NameSpaceURL, []byte(url)).String()
	filename := filepath.Join(db.root, basename)
	return ioutil.WriteFile(filename, content, 0755)
}

func (db *FeedRecordDB) RemoveRecord(url string) error {
	basename := uuid.NewMD5(uuid.NameSpaceURL, []byte(url)).String()
	filename := filepath.Join(db.root, basename)
	return os.Remove(filename)
}

func (db *FeedRecordDB) Records() ([]*FeedRecord, error) {
	return nil, nil
}

func ensureDir(p string) {
	os.MkdirAll(p, os.FileMode(0755))
}
