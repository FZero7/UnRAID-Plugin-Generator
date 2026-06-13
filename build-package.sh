#!/bin/bash
set -e

PLUGIN="smsups"
PKG_DIR="smsups-package"

echo "=== Building smsups-nut-bridge ==="
cd smsups-nut-bridge
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
    -ldflags="-w -s" \
    -o "../$PKG_DIR/usr/local/bin/smsups-nut-bridge" \
    ./cmd/bridge
cd ..
echo "Bridge binary: OK"

echo "=== Downloading ups-metrics from GitHub releases ==="
LATEST_TAG=$(curl -s https://api.github.com/repos/alexwbaule/ups-metrics/releases/latest \
    | grep '"tag_name":' | head -n1 | awk -F '"' '{print $4}')

if [ -z "$LATEST_TAG" ]; then
    echo "ERROR: Could not fetch latest ups-metrics release tag"
    exit 1
fi
echo "Latest ups-metrics release: $LATEST_TAG"

ASSET_URL=$(curl -s https://api.github.com/repos/alexwbaule/ups-metrics/releases/latest \
    | grep '"browser_download_url"' \
    | grep -i 'linux' | grep -i 'amd64' \
    | awk -F '"' '{print $4}' | head -1)

if [ -z "$ASSET_URL" ]; then
    echo "ERROR: No linux/amd64 asset found in latest release"
    echo "Available assets:"
    curl -s https://api.github.com/repos/alexwbaule/ups-metrics/releases/latest \
        | grep '"browser_download_url"' | awk -F '"' '{print $4}'
    exit 1
fi

echo "Downloading: $ASSET_URL"
wget -qO "$PKG_DIR/usr/local/bin/ups-metrics" "$ASSET_URL"
chmod +x "$PKG_DIR/usr/local/bin/ups-metrics"
echo "ups-metrics binary: OK ($LATEST_TAG)"

echo "=== Copying plugin files ==="
cp smsups-plugin/smsups.page          "$PKG_DIR/usr/local/emhttp/plugins/$PLUGIN/"
cp smsups-plugin/exec.php             "$PKG_DIR/usr/local/emhttp/plugins/$PLUGIN/php/"
cp smsups-plugin/default_config.yaml  "$PKG_DIR/usr/local/emhttp/plugins/$PLUGIN/"
cp smsups-plugin/default_bridge.yaml  "$PKG_DIR/usr/local/emhttp/plugins/$PLUGIN/"
cp smsups-plugin/rc.ups-metrics       "$PKG_DIR/etc/rc.d/rc.ups-metrics"
cp smsups-plugin/rc.smsups-bridge     "$PKG_DIR/etc/rc.d/rc.smsups-bridge"

chmod +x "$PKG_DIR/usr/local/bin/smsups-nut-bridge"
chmod +x "$PKG_DIR/etc/rc.d/"*

echo "=== Package contents ==="
find "$PKG_DIR" -type f | sort

echo ""
echo "=== Next step: run upg.py to generate .plg ==="
echo "python upg.py smsups.toml -p ./$PKG_DIR -o smsups.plg"
