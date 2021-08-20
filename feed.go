package zerogame

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

// VerificationMethod controls how an archive is verified before installation.
type VerificationMethod string

const (
	// AutoSelectMethod automatically determines which verification method to use.
	AutoSelectMethod VerificationMethod = "DefaultMethod"

	// DoNotVerifyMethod disables feed archive verification.
	DoNotVerifyMethod VerificationMethod = "NoVerification"

	// GPGDetachedSignatureMethod verifies archives signed with detached GPG signatures.
	//
	// The signer's public key must exist in the key ring.
	GPGDetachedSignatureMethod VerificationMethod = "DetachedSignatureMethod"
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

	// ArchiveuRL is used to GET this feed's archive.
	//
	// Required.
	ArchiveURL string `json:"archive_url"`

	// ArcchiveType is the archive file's extension.
	//
	// This is always zip and exists to extend support to future archive types.
	//
	// Required.
	ArchiveType string `json:"archive_type"`

	// GPGSignatureURL is used to GET this feed's archive's GPG signature.
	GPGSignatureURL string `json:"gpg_signature_url"`
}

// InstallFile describes how to install an archive on several platforms.
type InstallFile struct {
	Platforms []Platform `json:"platforms"`
}

// Platform is a platform-specific installation configuration.
type Platform struct {
	Name string `json:"name"`

	// Installs the archive.
	//
	// Required.
	InstallCommand []string `json:"install"`

	// Uninstalls the archive.
	//
	// Required.
	UninstallCommand []string `json:"uninstall"`

	// Runs the software provided in the archive.
	//
	// Required.
	RunCommand []string `json:"run"`
}

// InstallFeedOptions configures a call to InstallFeed.
type InstallFeedOptions struct {
	// UseCache controls whether existing feeds are refetched from the web.
	UseCache bool

	// VerificationMethod controls how an archive is verified.
	VerificationMethod VerificationMethod
}

// InstallFeed downloads feedURL and installs the corresponding archive on this machine.
func InstallFeed(_ context.Context, feedURL string, opts InstallFeedOptions) error {
	if strings.HasPrefix(feedURL, "file://") {
		feedURL = "file://" + filepath.Clean(feedURL[7:])
	}

	cache, err := newDefaultCache()
	if err != nil {
		return err
	}

	if opts.UseCache && cache.FeedArchiveExists(feedURL) {
		fmt.Fprintf(os.Stderr, "Feed %s is already installed.\n", feedURL)
		return nil
	}

	fmt.Fprintf(os.Stderr, "Loading feed: %s\n", feedURL)
	feed, err := fetchFeed(feedURL)
	if err != nil {
		return err
	}

	archive, err := fetchFeedArchive(feed, opts.VerificationMethod)
	if err != nil {
		return fmt.Errorf("failed to verify feed: %w. aborting", err)
	}

	if err := cache.WriteFeedArchive(feedURL, feed, archive, 0755); err != nil {
		return fmt.Errorf("failed to cache feed archive: %w. aborting", err)
	}

	archivePath, _ := cache.GetFeedArchive(feedURL)
	if err := installArchive(archivePath); err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}
	fmt.Fprintln(os.Stderr, "Installation complete!")
	return nil
}

func installArchive(archivePath string) error {
	files, err := extract(archivePath)
	if err != nil {
		return err
	}
	for _, file := range files {
		if filepath.Base(file) != "install.json" {
			continue
		}
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		var nf InstallFile
		if err := json.Unmarshal(data, &nf); err != nil {
			return err
		}
		for _, p := range nf.Platforms {
			if p.Name == currentPlatform() {
				cmd := exec.Command(p.InstallCommand[0], p.InstallCommand[1:]...)
				cmd.Stdin = os.Stdin
				cmd.Stderr = os.Stderr
				cmd.Stdout = os.Stdout
				return cmd.Run()
			}
		}
		return fmt.Errorf("cannot install the archive on this platform: %q", currentPlatform())
	}

	return errors.New("install.json not found")
}

func fetchFeed(feedURL string) (*Feed, error) {
	feedData, err := getURL(feedURL)
	if err != nil {
		return nil, err
	}
	var feed Feed
	if err := json.Unmarshal(feedData, &feed); err != nil {
		return nil, err
	}
	return &feed, nil
}

func fetchFeedArchive(feed *Feed, method VerificationMethod) ([]byte, error) {
	switch method {
	case AutoSelectMethod:
		method = DoNotVerifyMethod
		if feed.GPGSignatureURL != "" {
			method = GPGDetachedSignatureMethod
		}
		return fetchFeedArchive(feed, method)
	case DoNotVerifyMethod:
		return getURL(feed.ArchiveURL)
	case GPGDetachedSignatureMethod:
		return verifyFeedWithDetachedSignature(feed)
	}
	return nil, fmt.Errorf("unsupported verification method: %v", method)
}

func verifyFeedWithDetachedSignature(feed *Feed) ([]byte, error) {
	data, err := getURL(feed.ArchiveURL)
	if err != nil {
		return nil, fmt.Errorf("feed archive URL is invalid: %w", err)
	}
	signature, err := getURL(feed.GPGSignatureURL)
	if err != nil {
		return nil, fmt.Errorf("feed gpg signature URL is invalid: %w", err)
	}
	if err := verifyDetachedSignature(data, signature); err != nil {
		return nil, err
	}
	return data, nil
}

func verifyDetachedSignature(data, signature []byte) error {
	keyRing, err := readDefaultKeyRing()
	if err != nil {
		return err
	}
	for _, pgpPublicKey := range keyRing.GetKeys() {
		fmt.Fprintln(os.Stderr, "Verifying feed signature...")
		message := crypto.NewPlainMessage(data)
		signature := crypto.NewPGPSignature(signature)
		if err != nil {
			return err
		}
		keyRing, err := crypto.NewKeyRing(pgpPublicKey)
		if err != nil {
			return err
		}
		if keyRing.VerifyDetached(message, signature, crypto.GetUnixTime()) == nil {
			fmt.Fprintln(os.Stderr, "Feed signature verified")
			return nil
		}
	}
	return fmt.Errorf("feed signature could not be verified: %w", err)
}
