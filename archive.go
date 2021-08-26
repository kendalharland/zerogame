package zerogame

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type manifest struct {
	Platforms []platform `json:"platforms"`
}

type platform struct {
	Name             string   `json:"name"`
	InstallCommand   []string `json:"install"`
	UninstallCommand []string `json:"uninstall"`
	RunCommand       []string `json:"run"`
}

func InstallArchive(path string) error {
	workspace, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(workspace)

	log.Printf("[zerogame] extracting archive to %s", workspace)
	manifest, err := extract(path, workspace)
	if err != nil {
		return err
	}

	var args []string
	for _, p := range manifest.Platforms {
		if p.Name == currentPlatform() {
			args = p.InstallCommand
			break
		}
	}
	if args == nil {
		return fmt.Errorf("cannot be installed on platform: %s", currentPlatform())
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = workspace
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return err
	}
	if cmd.ProcessState.ExitCode() != 0 {
		return fmt.Errorf("process exited with code: %d", cmd.ProcessState.ExitCode())
	}

	fmt.Fprintln(os.Stderr, "Installation complete!")
	return nil
}

func decrypt(src, dst string) error {
	bytes, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	data, err := pgpDecrypt(bytes)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(dst, data, 0755); err != nil {
		return err
	}
	return nil
}

func extract(src, dst string) (*manifest, error) {
	files, err := unzip(src, dst)
	if err != nil {
		return nil, fmt.Errorf("unzip: %w", err)
	}

	var m manifest
	for _, file := range files {
		if filepath.Base(file) != "install.json" {
			continue
		}
		bytes, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(bytes, &m); err != nil {
			return nil, err
		}
		break
	}

	return &m, nil
}
