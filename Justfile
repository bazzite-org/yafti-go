default:
    @just --list

# Run dev server with live template reloading
dev:
    go tool templ generate --watch --cmd="env YAFTI_CONF=$PWD/example_config.yml go run ."

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
