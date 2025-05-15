package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Zeglius/yafti-go/config"
	"github.com/Zeglius/yafti-go/internal/consts"
	srv "github.com/Zeglius/yafti-go/server"
	webview "github.com/webview/webview_go"
)

//go:embed ../../static/**
var static embed.FS

func main() {
	// Enable WebView mode to prevent auto-shutdown
	consts.IsWebViewMode = true

	// Instantiate server
	server := srv.New()

	// Load static assets
	server.StaticAssets = &static

	// Start server in a goroutine
	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for server to be ready
	waitForServer()

	// Create WebView
	w := webview_go.New(true)
	defer w.Destroy()

	title := config.ConfStatus.Title
	if title == "" {
		title = "Bazzite Portal"
	}

	w.SetTitle(title)
	w.SetSize(1024, 768, webview_go.HintNone)

	// Connect to the local server
	w.Navigate(fmt.Sprintf("http://localhost:%s", consts.PORT))

	// Start the UI loop (blocking call)
	w.Run()

	// When WebView closes, exit the application
	os.Exit(0)
}

// waitForServer polls the server health endpoint until it's ready or times out
func waitForServer() {
	const maxAttempts = 20
	const interval = 100 * time.Millisecond

	url := fmt.Sprintf("http://localhost:%s/health", consts.PORT)

	for i := 0; i < maxAttempts; i++ {
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				log.Printf("Server is ready")
				return
			}
		}
		time.Sleep(interval)
	}

	log.Println("Warning: Server may not be ready yet, continuing anyway")
}
