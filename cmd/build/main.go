package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var configs = []BuildConfig{
	{GOOS: "linux", GOARCH: "amd64", ExecutableName: "zerogame"},
	{GOOS: "windows", GOARCH: "amd64", ExecutableName: "zerogame", ExecutableSuffix: "exe"},
	{GOOS: "windows", GOARCH: "386", ExecutableName: "zerogame", ExecutableSuffix: "exe"},
}

type BuildConfig struct {
	GOOS             string
	GOARCH           string
	ExecutableName   string
	ExecutableSuffix string
}

func main() {
	version := flag.String("v", "", "Release version string")
	dst := flag.String("o", "bin", "Build output directory")
	src := flag.String("i", filepath.Join("cmd", "zg"), "Source directory")
	gobin := flag.String("go", "/usr/bin/go", "Optional path to the go binary")

	flag.Parse()

	if err := execute(configs, *src, *dst, *version, *gobin); err != nil {
		log.Fatal(err)
	}
	os.Exit(1)
}

func execute(configs []BuildConfig, src, dst, version, gobin string) error {
	if version == "" {
		return errors.New("missing -v")
	}
	if dst == "" {
		return errors.New("missing -o")
	}
	if src == "" {
		return errors.New("missing -i")
	}

	dst, _ = filepath.Abs(filepath.Join(strings.Split(dst, string(os.PathSeparator))...))
	src, _ = filepath.Abs(filepath.Join(strings.Split(src, string(os.PathSeparator))...))

	var didFail bool
	for _, c := range configs {
		if err := build(c, version, src, dst, gobin); err != nil {
			didFail = true
			log.Println(err)
		}
	}

	if didFail {
		return errors.New("some builds failed")
	}
	return nil
}

func build(conf BuildConfig, version, src, dst, gobin string) error {
	output := filepath.Join(dst, conf.ExecutableName+version+"."+conf.GOOS+"-"+conf.GOARCH)
	if conf.ExecutableSuffix != "" {
		output += "." + conf.ExecutableSuffix
	}

	os.MkdirAll(dst, os.FileMode(0666))
	os.Remove(output)

	cmd := &exec.Cmd{
		Path:   "/usr/bin/env",
		Args:   []string{"/usr/bin/env", "GOOS=" + conf.GOOS, "GOARCH=" + conf.GOARCH, gobin, "build", "-o", output, src},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env:    os.Environ(),
	}

	//cmd.Env = append(cmd.Env, "GOOS", conf.GOOS)
	//cmd.Env = append(cmd.Env, "GOARCH", conf.GOARCH)

	fmt.Fprintf(os.Stderr, "Running %v\n", cmd.Args)
	return cmd.Run()
}
