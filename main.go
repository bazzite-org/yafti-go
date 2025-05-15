package main

import (
	"embed"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Zeglius/yafti-go/internal/consts"
	srv "github.com/Zeglius/yafti-go/server"
	"golang.org/x/sync/errgroup"

	// Import webview conditionally using build tags
	"github.com/webview/webview"
)

//go:embed static/**
var static embed.FS

func main() {
	// Get the wrapper command from environment variables
	// If YAFTI_EXEC_WRAPPER is set, the server will be started and the wrapper command will be executed
	cmd := os.Getenv("YAFTI_EXEC_WRAPPER")

	// Check if we should use the WebView UI
	useWebview := os.Getenv("YAFTI_USE_WEBVIEW") != "false" // Enable by default

	// Instantiate server
	server := srv.New()

	// Load static assets
	server.StaticAssets = &static

	// Start the server in all cases
	var errg errgroup.Group
	errg.Go(func() error {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Println("Server error:", err)
			return err
		}
		return nil
	})

	// Wait a moment for the server to start
	time.Sleep(500 * time.Millisecond)

	// If webview is enabled, use that and ignore wrapper command
	if useWebview {
		log.Println("Starting WebView UI")
		w := webview.New(true)
		defer w.Destroy()
		w.SetTitle("Yafti")
		w.SetSize(1024, 768, webview.HintNone)
		w.Navigate("http://localhost:" + consts.PORT)
		w.Run()
		return
	}

	// If no webview and no wrapper command is provided, the server is already running...
	if cmd == "" {
		// Wait for server to finish (which it won't unless terminated)
		if err := errg.Wait(); err != nil {
			log.Panic(err)
		}
		return
	}

	// Execute wrapper command
	cmd = strings.ReplaceAll(cmd, "%u", "http://localhost:"+consts.PORT)
	errg.Go(exec.Command("sh", "-c", cmd).Start)

	if err := errg.Wait(); err != nil {
		log.Panic(err)
	}
}
