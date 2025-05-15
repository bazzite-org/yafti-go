// +build webview

package main

import (
	"log"

	"github.com/Zeglius/yafti-go/internal/consts"
	"github.com/webview/webview"
)

// RunWebView starts a webview window for the application
func RunWebView() {
	log.Println("Starting Linux WebView UI")
	w := webview.New(true)
	defer w.Destroy()
	w.SetTitle(consts.APP_TITLE)
	w.SetSize(1024, 768, webview.HintNone)
	w.Navigate("http://localhost:" + consts.PORT)
	w.Run()
}
