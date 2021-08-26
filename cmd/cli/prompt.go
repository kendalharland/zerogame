package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
)

type prompt struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func (p *prompt) ReadString(prompt string) (string, error) {
	return p.readStringUntil(prompt, nil)
}

func (p *prompt) ReadNonEmptyString(prompt string) (string, error) {
	return p.readStringUntil(prompt, func(value string) error {
		if value == "" {
			return fmt.Errorf("please enter a non-empty string")

		}
		return nil
	})
}

func (p *prompt) ReadOneOf(prompt string, options ...string) (string, error) {
	return p.readStringUntil(prompt, func(value string) error {
		for i := range options {
			if value == options[i] {
				return nil
			}
		}
		return errors.New("please enter a valid archive type (one of: zip)")
	})
}

func (p *prompt) ReadURL(prompt string) (string, error) {
	return p.readStringUntil(prompt, func(value string) error {
		u, err := url.Parse(value)
		if err != nil {
			return errors.New("please enter a valid URL")
		}
		switch u.Scheme {
		case "file", "http", "https":
			return nil
		default:
			return errors.New("only file, http, and https schemes are allowed")
		}
	})
}

func (p *prompt) readStringUntil(prompt string, test func(string) error) (string, error) {
	r := bufio.NewReader(p.stdin)
	for {
		fmt.Fprint(p.stdout, prompt)
		value, err := r.ReadString('\n')
		if err != nil {
			return "", err
		}
		value = strings.TrimSpace(value)
		if test == nil {
			return value, nil
		} else if err := test(value); err == nil {
			return value, nil
		} else {
			fmt.Fprintln(p.stderr, err)
		}
	}
}
