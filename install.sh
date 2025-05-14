#!/bin/bash
set -e

# Configuration
INSTALL_DIR="${INSTALL_DIR:-/usr/local}"
CONFIG_DIR="${CONFIG_DIR:-/etc}"
DATA_DIR="${DATA_DIR:-/usr/share/yafti-go}"

# Create directories
mkdir -p "${INSTALL_DIR}/bin"
mkdir -p "${CONFIG_DIR}/default"
mkdir -p "${DATA_DIR}"

# Install binary
cp yafti-go "${INSTALL_DIR}/bin/"
chmod 755 "${INSTALL_DIR}/bin/yafti-go"

# Install config files
cp yafti.yml "${DATA_DIR}/"

# Install static files and templates
cp -r static "${DATA_DIR}/"
cp -r templates "${DATA_DIR}/" 2>/dev/null || true

# Create default config file
cat > "${CONFIG_DIR}/default/yafti-go" << EOF
YAFTI_CONF=${DATA_DIR}/yafti.yml
YAFTI_EXEC_WRAPPER=flatpak run org.mozilla.firefox --kiosk --new-instance %u
EOF

# Create systemd service file
mkdir -p "${DATA_DIR}/systemd"
cat > "${DATA_DIR}/systemd/yafti-go.service" << EOF
[Unit]
Description=Yet Another First Time Installer (Go version)
Wants=network-online.target
After=network-online.target

[Service]
Type=simple
EnvironmentFile=/etc/default/yafti-go
ExecStart=${INSTALL_DIR}/bin/yafti-go
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOF

echo "Installation completed successfully."
