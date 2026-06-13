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

echo "=== Building ups-metrics from source ==="
# ups-metrics has no release binaries — compile from latest main branch
UPS_TMP=$(mktemp -d)
wget -qO "$UPS_TMP/source.tar.gz" \
    https://github.com/alexwbaule/ups-metrics/archive/refs/heads/main.tar.gz
mkdir -p "$UPS_TMP/src"
tar -xzf "$UPS_TMP/source.tar.gz" -C "$UPS_TMP/src" --strip-components=1
cd "$UPS_TMP/src"
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
    -ldflags="-s -w" \
    -o "$OLDPWD/$PKG_DIR/usr/local/bin/ups-metrics" \
    ./cmd/ups-metrics
cd "$OLDPWD"
rm -rf "$UPS_TMP"
chmod +x "$PKG_DIR/usr/local/bin/ups-metrics"
echo "ups-metrics binary: OK (latest main)"

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
