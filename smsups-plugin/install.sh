#!/bin/bash
PLUGIN="smsups"
FLASH_DIR="/boot/config/plugins/$PLUGIN"
EMHTTP="/usr/local/emhttp/plugins/$PLUGIN"

echo "Configuring $PLUGIN..."

mkdir -p "$FLASH_DIR/conf"

if [ ! -f "$FLASH_DIR/conf/config.yaml" ]; then
    cp "$EMHTTP/default_config.yaml" "$FLASH_DIR/conf/config.yaml"
    echo "Default ups-metrics config installed"
fi

if [ ! -f "$FLASH_DIR/conf/bridge.yaml" ]; then
    cp "$EMHTTP/default_bridge.yaml" "$FLASH_DIR/conf/bridge.yaml"
    echo "Default bridge config installed"
fi

if [ ! -f "$FLASH_DIR/smsups.cfg" ]; then
    cat > "$FLASH_DIR/smsups.cfg" <<EOF
UPS_METRICS_ENABLED=yes
BRIDGE_ENABLED=yes
EOF
    echo "Default smsups.cfg created"
fi

echo ""
echo "NOTE: Edit $FLASH_DIR/conf/config.yaml before starting."
echo "      Set device.address and device.login.password."
echo ""
echo "NUT dummy-ups configuration (add to /etc/nut/ups.conf):"
echo "  [smsups]"
echo "    driver = dummy-ups"
echo "    port = /tmp/smsups.dev"
echo "    desc = \"SMS Brasil UPS via Bridge\""
echo ""

/etc/rc.d/rc.ups-metrics start
/etc/rc.d/rc.smsups-bridge start

echo "$PLUGIN installed successfully!"
