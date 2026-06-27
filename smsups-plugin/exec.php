<?php
$plugin    = "smsups";
$flash_dir = "/boot/config/plugins/$plugin";
$cfg_file  = "$flash_dir/smsups.cfg";
$ups_cfg   = "$flash_dir/conf/config.yaml";
$bridge_cfg= "$flash_dir/conf/bridge.yaml";

$action  = $_REQUEST['action']  ?? '';
$service = $_REQUEST['service'] ?? '';

$rc = [
    'ups-metrics'  => '/etc/rc.d/rc.ups-metrics',
    'bridge'       => '/etc/rc.d/rc.smsups-bridge',
    'victorialogs' => '/etc/rc.d/rc.victorialogs',
];

function runrc($rc, $cmd) {
    return trim(shell_exec("$rc $cmd 2>&1"));
}

switch ($action) {

    case 'status':
        $script = $rc[$service] ?? null;
        echo $script ? runrc($script, 'status') : "unknown service";
        break;

    case 'start':
    case 'stop':
    case 'restart':
        $script = $rc[$service] ?? null;
        echo $script ? runrc($script, $action) : "unknown service";
        break;

    case 'start_all':
        echo runrc($rc['victorialogs'], 'start') . "\n";
        echo runrc($rc['ups-metrics'],  'start') . "\n";
        echo runrc($rc['bridge'],       'start');
        break;

    case 'stop_all':
        echo runrc($rc['ups-metrics'],  'stop') . "\n";
        echo runrc($rc['bridge'],       'stop') . "\n";
        echo runrc($rc['victorialogs'], 'stop');
        break;

    case 'restart_all':
        echo runrc($rc['victorialogs'], 'restart') . "\n";
        echo runrc($rc['ups-metrics'],  'restart') . "\n";
        echo runrc($rc['bridge'],       'restart');
        break;

    case 'save':
        $cfg_path = ($service === 'bridge') ? $bridge_cfg : $ups_cfg;
        if (!is_dir(dirname($cfg_path))) {
            mkdir(dirname($cfg_path), 0755, true);
        }
        file_put_contents($cfg_path, $_POST['config'] ?? '');
        $script = $rc[$service] ?? null;
        $out = $script ? runrc($script, 'restart') : '';
        echo "Config saved. $out";
        break;

    case 'save_vlogs':
        $retention = preg_replace('/[^a-zA-Z0-9]/', '', $_POST['retention'] ?? '24h');
        $port      = preg_replace('/[^0-9]/',        '', $_POST['port']      ?? '9428');
        if (!$retention) $retention = '24h';
        if (!$port)      $port      = '9428';
        $keys = ['VICTORIALOGS_RETENTION' => $retention, 'VICTORIALOGS_PORT' => $port];
        $lines = file_exists($cfg_file) ? file($cfg_file, FILE_IGNORE_NEW_LINES) : [];
        foreach ($keys as $key => $val) {
            $found = false;
            foreach ($lines as &$line) {
                if (strpos($line, "$key=") === 0) { $line = "$key=$val"; $found = true; }
            }
            if (!$found) $lines[] = "$key=$val";
        }
        if (!is_dir($flash_dir)) mkdir($flash_dir, 0755, true);
        file_put_contents($cfg_file, implode("\n", $lines) . "\n");
        echo "Config saved. " . runrc($rc['victorialogs'], 'restart');
        break;

    case 'set_autostart':
        $enabled = ($_REQUEST['enabled'] ?? 'no') === 'yes' ? 'yes' : 'no';
        $key = ($service === 'bridge') ? 'BRIDGE_ENABLED' : 'UPS_METRICS_ENABLED';

        $lines = file_exists($cfg_file) ? file($cfg_file, FILE_IGNORE_NEW_LINES) : [];
        $found = false;
        foreach ($lines as &$line) {
            if (strpos($line, "$key=") === 0) {
                $line = "$key=$enabled";
                $found = true;
            }
        }
        if (!$found) $lines[] = "$key=$enabled";
        if (!is_dir($flash_dir)) mkdir($flash_dir, 0755, true);
        file_put_contents($cfg_file, implode("\n", $lines) . "\n");
        echo "Auto-start for $service set to $enabled";
        break;

    default:
        echo "invalid action";
}
?>
