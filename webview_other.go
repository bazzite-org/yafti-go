//go:build !webview
// +build !webview

package main

// This file provides stub implementations for non-webview builds
// Because we're using build tags to conditionally include webview support

// RunWebView is a stub for non-webview builds
func RunWebView() {
	// This is intentionally empty on non-webview builds
}
