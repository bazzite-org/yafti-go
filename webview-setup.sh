#!/bin/bash

# This script helps set up the WebView version of the yafti-go app on Linux

# Detect the Linux distribution
if [ -f /etc/os-release ]; then
    . /etc/os-release
    DISTRO=$ID
else
    echo "Unable to determine Linux distribution"
    exit 1
fi

# Install dependencies based on the distribution
case $DISTRO in
    ubuntu|debian|pop|elementary|linuxmint|neon)
        echo "Installing WebView dependencies for Debian-based distribution..."
        sudo apt-get update
        sudo apt-get install -y libwebkit2gtk-4.1-dev libgtk-3-dev
        ;;
    fedora|rhel|centos|rocky|almalinux)
        echo "Installing WebView dependencies for Red Hat-based distribution..."
        sudo dnf install -y webkit2gtk4.1-devel gtk3-devel
        ;;
    arch|manjaro|endeavouros)
        echo "Installing WebView dependencies for Arch-based distribution..."
        sudo pacman -Sy webkit2gtk gtk3
        ;;
    opensuse*)
        echo "Installing WebView dependencies for openSUSE..."
        sudo zypper install webkit2gtk3-devel gtk3-devel
        ;;
    *)
        echo "Unsupported distribution: $DISTRO"
        echo "Please manually install WebKit2GTK and GTK3 development packages"
        ;;
esac

# Build the WebView version
echo "Building WebView version..."
go build -tags webview -o yafti-go-gui

# Check if build was successful
if [ $? -eq 0 ]; then
    echo "Build successful! WebView binary created as 'yafti-go-gui'"
    echo "You can run it with:"
    echo "YAFTI_USE_WEBVIEW=true ./yafti-go-gui"
    
    # Option to run immediately
    read -p "Would you like to run it now? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        YAFTI_USE_WEBVIEW=true ./yafti-go-gui
    fi
else
    echo "Build failed. Please check the error messages above."
fi
