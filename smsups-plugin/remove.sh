#!/bin/bash
PLUGIN="smsups"

echo "Removing $PLUGIN..."

/etc/rc.d/rc.victorialogs stop 2>/dev/null
/etc/rc.d/rc.ups-metrics stop 2>/dev/null
/etc/rc.d/rc.smsups-bridge stop 2>/dev/null

rm -f /usr/local/bin/ups-metrics
rm -f /usr/local/bin/smsups-nut-bridge
rm -rf /usr/local/emhttp/plugins/$PLUGIN
rm -f /etc/rc.d/rc.victorialogs
rm -f /etc/rc.d/rc.ups-metrics
rm -f /etc/rc.d/rc.smsups-bridge

echo "Note: Config files at /boot/config/plugins/$PLUGIN were NOT removed."
echo "$PLUGIN removed."
