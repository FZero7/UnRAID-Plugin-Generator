package prometheus_test

import (
	"testing"

	"github.com/bobbintb/smsups-nut-bridge/internal/prometheus"
)

const sampleMetrics = `# HELP ups_status The status of the UPS
# TYPE ups_status gauge
ups_status{host="SMSGamer2000",type="battery_level",unit="%"} 100
ups_status{host="SMSGamer2000",type="input_voltage",unit="VAC"} 233.4
ups_status{host="SMSGamer2000",type="output_voltage",unit="VAC"} 216.9
ups_status{host="SMSGamer2000",type="output_frequency",unit="Hz"} 59.8
ups_status{host="SMSGamer2000",type="ups_load",unit="%"} 12.4
ups_status{host="SMSGamer2000",type="ups_temperature",unit="°C"} 53.8
ups_status{host="SMSGamer2000",type="ups_is_interative",unit=""} 1
# HELP ups_state The states of the UPS
# TYPE ups_state gauge
ups_state{host="SMSGamer2000",state="on_grid"} 1
ups_state{host="SMSGamer2000",state="battery_fail"} 0
ups_state{host="SMSGamer2000",state="on_boost"} 0
ups_state{host="SMSGamer2000",state="on_bypass"} 0
ups_state{host="SMSGamer2000",state="on_test"} 0
ups_state{host="SMSGamer2000",state="on_high_power"} 0
ups_state{host="SMSGamer2000",state="battery_is_full"} 1
`

func TestParse_UPSStatus(t *testing.T) {
	metrics, err := prometheus.Parse(sampleMetrics)
	if err != nil {
		t.Fatal(err)
	}

	found := findMetric(metrics, "ups_status", "type", "battery_level")
	if found == nil {
		t.Fatal("battery_level not found")
	}
	if found.Value != 100 {
		t.Errorf("battery_level: expected 100, got %f", found.Value)
	}

	found = findMetric(metrics, "ups_status", "type", "input_voltage")
	if found == nil {
		t.Fatal("input_voltage not found")
	}
	if found.Value != 233.4 {
		t.Errorf("input_voltage: expected 233.4, got %f", found.Value)
	}
}

func TestParse_UPSState(t *testing.T) {
	metrics, err := prometheus.Parse(sampleMetrics)
	if err != nil {
		t.Fatal(err)
	}

	found := findMetric(metrics, "ups_state", "state", "on_grid")
	if found == nil {
		t.Fatal("on_grid not found")
	}
	if found.Value != 1 {
		t.Errorf("on_grid: expected 1, got %f", found.Value)
	}

	found = findMetric(metrics, "ups_state", "state", "battery_fail")
	if found == nil {
		t.Fatal("battery_fail not found")
	}
	if found.Value != 0 {
		t.Errorf("battery_fail: expected 0, got %f", found.Value)
	}
}

func TestParse_IgnoresComments(t *testing.T) {
	metrics, err := prometheus.Parse("# HELP foo bar\n# TYPE foo gauge\nfoo{a=\"b\"} 1\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(metrics) != 1 {
		t.Errorf("expected 1 metric, got %d", len(metrics))
	}
}

func findMetric(metrics []prometheus.Metric, name, labelKey, labelVal string) *prometheus.Metric {
	for i := range metrics {
		if metrics[i].Name == name && metrics[i].Labels[labelKey] == labelVal {
			return &metrics[i]
		}
	}
	return nil
}
