package views

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	_ "embed"

	"github.com/kendalharland/zerogame"
	"github.com/webview/webview"
)

//go:embed index.html
var indexHtml []byte

func NavigateToIndex(w webview.WebView) error {
	p := &indexPage{
		LogFunc: func(m string, args ...interface{}) {
			s := fmt.Sprintf(m, args...)
			w.Eval(fmt.Sprintf("log(%q);", s))
		},
	}
	w.Bind("zg_install_feed", p.installFeed)
	w.Navigate(fmt.Sprintf(`data:text/html,<!doctype html><html><body>%s</body></html>`, string(indexHtml)))
	return nil
}

type indexPage struct {
	LogFunc func(string, ...interface{})
}

func (p *indexPage) installFeed(feedURL string, noCache bool) error {
	home, _ := os.UserHomeDir()
	workspace := filepath.Join(home, ".config", "zerogame")
	db := zerogame.NewFeedRecordDB(filepath.Join(workspace, "db"))
	repo := zerogame.NewFeedRepository(workspace, db)
	return zerogame.InstallFeed(context.Background(), feedURL, repo)
}
