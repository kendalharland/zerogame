package main

import (
	"log"

	"github.com/kendalharland/zerogame/cmd/gui/views"
	"github.com/webview/webview"
)

const (
	version    = "0.1"
	title      = "zerogame " + version
	viewWidth  = 600
	viewHeight = 800
)

func main() {
	w := webview.New(true)
	defer w.Destroy()

	w.SetTitle(title)
	w.SetSize(viewWidth, viewHeight, webview.HintNone)
	if err := views.NavigateToIndex(w); err != nil {
		log.Fatal(err)
	}
	w.Run()
}
