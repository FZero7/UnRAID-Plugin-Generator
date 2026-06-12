package mapper

import (
	"fmt"
	"strings"

	"github.com/bobbintb/smsups-nut-bridge/internal/prometheus"
)

type NUTVars map[string]string

var analogMap = map[string]string{
	"battery_level":    "battery.charge",
	"input_voltage":    "input.voltage",
	"output_voltage":   "output.voltage",
	"output_frequency": "output.frequency",
	"ups_load":         "ups.load",
	"ups_temperature":  "ups.temperature",
}

func Map(metrics []prometheus.Metric, manufacturer string) NUTVars {
	vars := NUTVars{}
	var statusFlags []string
	var model string

	for _, m := range metrics {
		if host := m.Labels["host"]; host != "" && model == "" {
			model = host
		}

		switch m.Name {
		case "ups_status":
			typ := m.Labels["type"]
			if nutKey, ok := analogMap[typ]; ok {
				vars[nutKey] = formatValue(m.Value)
			}
			if typ == "ups_is_interative" && m.Value == 1 {
				vars["ups.type"] = "line-interactive"
			}

		case "ups_state":
			state := m.Labels["state"]
			active := m.Value == 1
			switch state {
			case "on_grid":
				if active {
					statusFlags = append(statusFlags, "OL")
				} else {
					statusFlags = append(statusFlags, "OB")
				}
			case "battery_fail":
				if active {
					statusFlags = append(statusFlags, "LB")
				}
			case "on_boost":
				if active {
					statusFlags = append(statusFlags, "BOOST")
				}
			case "on_bypass":
				if active {
					statusFlags = append(statusFlags, "BYPASS")
				}
			case "on_test":
				if active {
					statusFlags = append(statusFlags, "TEST")
				}
			case "on_high_power":
				if active {
					statusFlags = append(statusFlags, "OVER")
				}
			}
		}
	}

	if len(statusFlags) > 0 {
		vars["ups.status"] = strings.Join(statusFlags, " ")
	}
	if model != "" {
		vars["ups.model"] = model
	}
	if manufacturer != "" {
		vars["ups.mfr"] = manufacturer
	}
	return vars
}

func formatValue(v float64) string {
	if v == float64(int64(v)) {
		return fmt.Sprintf("%d", int64(v))
	}
	return fmt.Sprintf("%g", v)
}
