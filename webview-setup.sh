# WebView Development Setup for Ubuntu/Debian
sudo apt-get update
sudo apt-get install -y libwebkit2gtk-4.1-dev libgtk-3-dev

# WebView Development Setup for Fedora
# sudo dnf install -y webkit2gtk4.1-devel gtk3-devel

# Build the WebView version
go build -tags webview -o yafti-go-gui

# Run with WebView interface
YAFTI_USE_WEBVIEW=true ./yafti-go-gui
