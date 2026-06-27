package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Source SourceConfig `yaml:"source"`
	NUT    NUTConfig    `yaml:"nut"`
}

type SourceConfig struct {
	Address  string        `yaml:"address"`
	Port     string        `yaml:"port"`
	Interval time.Duration `yaml:"interval"`
	Timeout  time.Duration `yaml:"timeout"`
}

type NUTConfig struct {
	DevFile      string  `yaml:"dev_file"`
	UPSName      string  `yaml:"ups_name"`
	Description  string  `yaml:"description"`
	Manufacturer string  `yaml:"manufacturer"`
	NominalVA           float64 `yaml:"nominal_va"`            // optional: ups.power.nominal
	NominalWatts        float64 `yaml:"nominal_watts"`         // optional: ups.realpower.nominal (priority over va×pf)
	PowerFactor         float64 `yaml:"power_factor"`          // optional: ups.powerfactor (also used to derive watts from va)
	LowBatteryThreshold float64 `yaml:"low_battery_threshold"` // % below which LB status is set (default 20)
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	if cfg.Source.Interval == 0 {
		cfg.Source.Interval = 15 * time.Second
	}
	if cfg.Source.Timeout == 0 {
		cfg.Source.Timeout = 5 * time.Second
	}
	if cfg.Source.Port == "" {
		cfg.Source.Port = "9110"
	}
	if cfg.NUT.DevFile == "" {
		cfg.NUT.DevFile = "/tmp/smsups.dev"
	}
	if cfg.NUT.LowBatteryThreshold == 0 {
		cfg.NUT.LowBatteryThreshold = 20
	}
	return &cfg, nil
}
