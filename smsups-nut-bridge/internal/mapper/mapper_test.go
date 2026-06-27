package mapper_test

import (
	"testing"

	"smsups-nut-bridge/internal/config"
	"smsups-nut-bridge/internal/mapper"
	"smsups-nut-bridge/internal/prometheus"
)

func baseMetrics() []prometheus.Metric {
	return []prometheus.Metric{
		{Name: "ups_status", Labels: map[string]string{"host": "TestUPS", "type": "battery_level", "unit": "%"}, Value: 85},
		{Name: "ups_status", Labels: map[string]string{"host": "TestUPS", "type": "input_voltage", "unit": "VAC"}, Value: 230.0},
		{Name: "ups_status", Labels: map[string]string{"host": "TestUPS", "type": "output_voltage", "unit": "VAC"}, Value: 220.5},
		{Name: "ups_status", Labels: map[string]string{"host": "TestUPS", "type": "output_frequency", "unit": "Hz"}, Value: 60.0},
		{Name: "ups_status", Labels: map[string]string{"host": "TestUPS", "type": "ups_load", "unit": "%"}, Value: 15.0},
		{Name: "ups_status", Labels: map[string]string{"host": "TestUPS", "type": "ups_temperature", "unit": "°C"}, Value: 40.0},
		{Name: "ups_status", Labels: map[string]string{"host": "TestUPS", "type": "ups_is_interative", "unit": ""}, Value: 1},
		{Name: "ups_state", Labels: map[string]string{"host": "TestUPS", "state": "on_grid"}, Value: 1},
		{Name: "ups_state", Labels: map[string]string{"host": "TestUPS", "state": "battery_fail"}, Value: 0},
		{Name: "ups_state", Labels: map[string]string{"host": "TestUPS", "state": "on_boost"}, Value: 0},
		{Name: "ups_state", Labels: map[string]string{"host": "TestUPS", "state": "on_bypass"}, Value: 0},
		{Name: "ups_state", Labels: map[string]string{"host": "TestUPS", "state": "on_test"}, Value: 0},
		{Name: "ups_state", Labels: map[string]string{"host": "TestUPS", "state": "on_high_power"}, Value: 0},
	}
}

func TestMap_AnalogMetrics(t *testing.T) {
	vars := mapper.Map(baseMetrics(), config.NUTConfig{Manufacturer: "SMS Brasil"})

	cases := map[string]string{
		"battery.charge":   "85",
		"input.voltage":    "230",
		"output.voltage":   "220.5",
		"output.frequency": "60",
		"ups.load":         "15",
		"ups.temperature":  "40",
		"ups.type":         "line-interactive",
	}
	for k, want := range cases {
		if vars[k] != want {
			t.Errorf("%s: expected %q, got %q", k, want, vars[k])
		}
	}
}

func TestMap_StatusOnGrid(t *testing.T) {
	vars := mapper.Map(baseMetrics(), config.NUTConfig{Manufacturer: "SMS Brasil"})
	if vars["ups.status"] != "OL" {
		t.Errorf("ups.status: expected OL, got %q", vars["ups.status"])
	}
}

func TestMap_StatusOnBattery(t *testing.T) {
	metrics := baseMetrics()
	for i := range metrics {
		if metrics[i].Labels["state"] == "on_grid" {
			metrics[i].Value = 0
		}
	}
	vars := mapper.Map(metrics, config.NUTConfig{Manufacturer: "SMS Brasil"})
	if vars["ups.status"] != "OB" {
		t.Errorf("ups.status: expected OB, got %q", vars["ups.status"])
	}
}

func TestMap_StatusReplaceBattery(t *testing.T) {
	metrics := baseMetrics()
	for i := range metrics {
		if metrics[i].Labels["state"] == "battery_fail" {
			metrics[i].Value = 1
		}
	}
	vars := mapper.Map(metrics, config.NUTConfig{Manufacturer: "SMS Brasil"})
	status := vars["ups.status"]
	if !containsStr(status, "RB") {
		t.Errorf("ups.status: expected RB, got %q", status)
	}
}

func TestMap_StatusLowBatteryThreshold(t *testing.T) {
	metrics := baseMetrics()
	for i := range metrics {
		if metrics[i].Labels["type"] == "battery_level" {
			metrics[i].Value = 15
		}
		if metrics[i].Labels["state"] == "on_grid" {
			metrics[i].Value = 0
		}
	}
	vars := mapper.Map(metrics, config.NUTConfig{Manufacturer: "SMS Brasil", LowBatteryThreshold: 20})
	if vars["ups.status"] != "OB LB" {
		t.Errorf("ups.status: expected 'OB LB', got %q", vars["ups.status"])
	}
}

func TestMap_StatusNotLowBatteryAboveThreshold(t *testing.T) {
	metrics := baseMetrics() // battery_level = 85
	vars := mapper.Map(metrics, config.NUTConfig{Manufacturer: "SMS Brasil", LowBatteryThreshold: 20})
	if containsStr(vars["ups.status"], "LB") {
		t.Errorf("ups.status: should not contain LB at 85%%, got %q", vars["ups.status"])
	}
}

func TestMap_StatusFlags(t *testing.T) {
	metrics := baseMetrics()
	for i := range metrics {
		switch metrics[i].Labels["state"] {
		case "on_boost":
			metrics[i].Value = 1
		case "on_bypass":
			metrics[i].Value = 1
		}
	}
	vars := mapper.Map(metrics, config.NUTConfig{Manufacturer: "SMS Brasil"})
	status := vars["ups.status"]
	if !contains(status, "BOOST") {
		t.Errorf("expected BOOST in status, got %q", status)
	}
	if !contains(status, "BYPASS") {
		t.Errorf("expected BYPASS in status, got %q", status)
	}
}

func TestMap_ModelFromHost(t *testing.T) {
	vars := mapper.Map(baseMetrics(), config.NUTConfig{Manufacturer: "SMS Brasil"})
	if vars["ups.model"] != "TestUPS" {
		t.Errorf("ups.model: expected TestUPS, got %q", vars["ups.model"])
	}
}

func TestMap_Manufacturer(t *testing.T) {
	vars := mapper.Map(baseMetrics(), config.NUTConfig{Manufacturer: "SMS Brasil"})
	if vars["ups.mfr"] != "SMS Brasil" {
		t.Errorf("ups.mfr: expected 'SMS Brasil', got %q", vars["ups.mfr"])
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
