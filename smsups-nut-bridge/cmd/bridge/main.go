package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"smsups-nut-bridge/internal/config"
	"smsups-nut-bridge/internal/mapper"
	"smsups-nut-bridge/internal/nut"
	"smsups-nut-bridge/internal/prometheus"
)

func main() {
	cfgPath := flag.String("config", "conf/bridge.yaml", "path to bridge.yaml")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("loading config: %v", err)
	}

	log.Printf("smsups-nut-bridge starting: source=%s:%s interval=%s dev=%s",
		cfg.Source.Address, cfg.Source.Port, cfg.Source.Interval, cfg.NUT.DevFile)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	ticker := time.NewTicker(cfg.Source.Interval)
	defer ticker.Stop()

	if err := run(ctx, cfg); err != nil {
		log.Printf("initial run error (will retry): %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("shutting down")
			os.Exit(0)
		case <-ticker.C:
			if err := run(ctx, cfg); err != nil {
				log.Printf("run error (will retry on next tick): %v", err)
			}
		}
	}
}

func run(ctx context.Context, cfg *config.Config) error {
	metrics, err := prometheus.Fetch(ctx, cfg.Source.Address, cfg.Source.Port, cfg.Source.Timeout)
	if err != nil {
		return err
	}
	vars := mapper.Map(metrics, cfg.NUT)
	if cfg.NUT.UPSName != "" {
		vars["ups.id"] = cfg.NUT.UPSName
	}
	if cfg.NUT.Description != "" {
		vars["ups.desc"] = cfg.NUT.Description
	}
	if err := nut.Write(cfg.NUT.DevFile, vars); err != nil {
		return err
	}
	log.Printf("updated %s: status=%s charge=%s%%", cfg.NUT.DevFile, vars["ups.status"], vars["battery.charge"])
	return nil
}
