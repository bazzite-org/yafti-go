export YAFTI_EXEC_WRAPPER := "flatpak run org.mozilla.firefox --kiosk --new-instance %u"

default:
    @just --list

# Run dev server with live template reloading
dev:
    go tool templ generate --watch --cmd="env YAFTI_CONF=$PWD/yafti.yml go run ."

# Run with custom config file
run config="yafti.yml":
    env YAFTI_CONF=$PWD/{{config}} go run .

# Build the application
build:
    go build -o yafti-go

# Generate templ files
templ:
    templ generate

# Same as templ but using GOPATH to find the binary
templ-gopath:
    $(go env GOPATH)/bin/templ generate

# Run with bazzite config
bazzite:
    env YAFTI_CONF=$PWD/yafti.yml go run .

# Build and run with bazzite config
bazzite-build:
    go build -o yafti-go && env YAFTI_CONF=$PWD/yafti.yml ./yafti-go

# Clean build artifacts
clean:
    rm -f yafti-go

# Install dependencies
deps:
    go mod tidy
    go install github.com/a-h/templ/cmd/templ@latest

# full clean and rebuild and deploy for bazzite example script
bazzite-full:
    just clean
    just bazzite-build

# Complete rebuild: clean, regenerate templates, build and run with yafti.yml
rebuild-all:
    @echo "=== Cleaning build artifacts ==="
    just clean
    @echo "=== Removing generated templ files ==="
    find ./ui -name "*_templ.go" -delete
    @echo "=== Installing templ tool if needed ==="
    go install github.com/a-h/templ/cmd/templ@latest
    @echo "=== Regenerating templ files ==="
    $(go env GOPATH)/bin/templ generate
    @echo "=== Rebuilding application ==="
    go build -o yafti-go
    @echo "=== Running application with yafti.yml ==="
    go tool templ generate --watch --cmd="env YAFTI_CONF=$PWD/yafti.yml go run ."

# Run with WebView interface
run-webview config="yafti.yml":
    env YAFTI_CONF=$PWD/{{config}} YAFTI_USE_WEBVIEW=true go run .

# Build with WebView support
build-webview:
    go build -o yafti-go-gui -tags webview

# Build both standard and WebView versions
build-all:
    go build -o yafti-go
    go build -o yafti-go-gui -tags webview
