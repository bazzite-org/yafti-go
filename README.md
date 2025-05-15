# Yafti-Go

**Yafti-Go** (Yet Another First Time Installation - Go) is a web-based application that helps users install and configure Bazzite OS. It provides a sleek, user-friendly interface for selecting and running installation scripts.


## Features

- **Web-based Interface**: Access the installer through your browser or the integrated WebView
- **Configurable**: Define your installation options in YAML configuration files
- **Visual Selection**: Easily choose which components to install with toggle switches
- **Command Execution**: Runs installation commands with real-time output display
- **Multiple Deployment Options**: Run as a standalone binary, with an external browser, or as a self-contained Flatpak

## Installation

### Requirements

- Go 1.21 or later
- Just command runner (optional, but recommended)

### Quick Start

1. Clone this repository
   ```bash
   git clone https://github.com/Zeglius/yafti-go.git
   cd yafti-go
   ```

2. Install dependencies
   ```bash
   go mod tidy
   go install github.com/a-h/templ/cmd/templ@latest
   ```

3. Run with default config
   ```bash
   export YAFTI_CONF="$(pwd)/yafti.yml" YAFTI_EXEC_WRAPPER="flatpak run org.mozilla.firefox --kiosk --new-instance %u"
   go run main.go
   ```

4. Access the web interface at http://localhost:3169

### Flatpak Version

The Flatpak version includes an integrated WebView, so no external browser is required.

1. Download the Flatpak bundle from the latest release
2. Install it:
   ```bash
   flatpak install --user bazzite-portal.flatpak
   ```
3. Run the app:
   ```bash
   flatpak run com.bazzite.Portal
   ```

### WebView Version

The WebView version provides a standalone window with no external browser dependency:

```bash
# Build the WebView version
go build -o yafti-go-webview ./cmd/webview

# Run the WebView version
YAFTI_CONF="$(pwd)/yafti.yml" ./yafti-go-webview
```

## Configuration

Yafti-Go is configured using YAML files. The configuration file specifies screens (pages) with actions (installable components).

Example configuration:

```yaml
title: Bazzite Portal
screens:
  - title: "Setting up Bazzite"
    actions:
      - id: "decky-loader"
        title: "Decky Loader"
        description: "A plugin loader for the Steam Deck"
        default: false
        script: "echo Installing Decky Loader"
```

By default, Yafti-Go looks for a configuration file at `/usr/share/yafti/conf.yml`, but you can specify a custom path using the `YAFTI_CONF` environment variable.

## Development

This project uses:
- [Echo](https://echo.labstack.com/) - Web framework
- [Templ](https://templ.guide/) - HTML templating
- [HTMX](https://htmx.org/) - Dynamic HTML updates
- [DaisyUI](https://daisyui.com/) - UI components
- [TailwindCSS](https://tailwindcss.com/) - CSS framework
- [WebView](https://github.com/webview/webview) - Embedded WebView for standalone UI

### Common Commands

If you have [Just](https://github.com/casey/just) installed, you can use the following commands:

```bash
# List all available commands
just

# Run with live reload during development
just dev

# Run with a specific config file
just run myconfigfile.yml

# Run with the Bazzite configuration
just bazzite

# Build the standard application
just build

# Build the WebView version
just webview-build

# Run the WebView version
just webview-run

# Build Flatpak package
just flatpak-build

# Install the Flatpak package
just flatpak-install

# Run the Flatpak
just flatpak-run

# Generate templ files
just templ

# Install dependencies
just deps
```

## Project Structure

```
├── config/           # Configuration handling
├── internal/         # Internal packages
├── static/           # Static assets (CSS, JS, images)
├── ui/               # UI components and templates
│   ├── components/   # Reusable UI components
│   └── pages/        # Page templates
├── main.go           # Application entry point
└── yafti.yml         # Example configuration
```

## License

[License details here]

## Contributions

Contributions are welcome! Please feel free to submit a Pull Request.
