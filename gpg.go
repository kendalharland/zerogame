package zerogame

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

const gpgKeyDir = ".gnupg"

// Default GPG key ring files in order of precendence.
// The first file that exists is used for feed verification.
// See https://www.gnupg.org/faq/whats-new-in-2.1.html#keybox
var keyRingPaths = []string{
	"pubring.gpg",
	"pubring.kbx",
}

func pgpDecrypt(message []byte) ([]byte, error) {
	kr, err := readDefaultKeyRing()
	if err != nil {
		return nil, err
	}
	m := crypto.NewPGPMessage(message)
	d, err := kr.Decrypt(m, nil, time.Now().UnixNano())
	if err != nil {
		return nil, err
	}
	return d.Data, nil
}

func showKeyBoxWarning(keyBoxFile, wantKeyRingFile string) {
	fmt.Fprintf(os.Stderr, `
-------
WARNING: A GPG keybox file was found at %s but this tool is
incompatible with the GPG keybox format. To fix this issue,
convert the keybox file to a keyring by running:

    gpg --no-default-keyring --keyring %s --export > %s

Then try running this tool again.
Feed verification will be disabled until this is resolved.
-------

`, keyBoxFile, keyBoxFile, wantKeyRingFile)
}

func readDefaultKeyRing() (*crypto.KeyRing, error) {
	keyRingFile, err := locateKeyRing()
	if err != nil {
		return nil, err
	}
	if strings.HasSuffix(keyRingFile, "kbx") {
		wantKeyringFile := replaceFileExtension(keyRingFile, "gpg")
		showKeyBoxWarning(keyRingFile, wantKeyringFile)
		return nil, errors.New("invalid keyring format")
	}
	return readKeyRing(keyRingFile)
}

func readKeyRing(filename string) (*crypto.KeyRing, error) {
	fmt.Fprintln(os.Stderr, "Loading keyring from "+filename)
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	entities, err := openpgp.ReadKeyRing(fd)
	if err != nil {
		return nil, err
	}
	ring, _ := crypto.NewKeyRing(nil)
	for _, entity := range entities {
		key, err := crypto.NewKeyFromEntity(entity)
		if err != nil {
			return nil, err
		}
		ring.AddKey(key)
	}
	return ring, nil
}

func locateKeyRing() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("could not get user's home directory")
	}
	for _, p := range keyRingPaths {
		filename := filepath.Join(home, gpgKeyDir, p)
		_, err := os.Stat(filename)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			continue
		}
		return filepath.Clean(filename), nil
	}
	return "", errors.New("could not locate keyring file")
}
