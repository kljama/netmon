package main

import (
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)
	log.Println("Netmon starting...")

	cfg, err := LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Printf("Loaded config: %v networks, scan interval %v, discovery interval %v, timeout %v", len(cfg.Networks), cfg.ScanInterval, cfg.DiscoveryInterval, cfg.Timeout)

	metrics, err := NewMetricsClient()
	if err != nil {
		log.Fatalf("Failed to initialize metrics client: %v", err)
	}
	defer metrics.Close()

	scanner, err := NewScanner(cfg, metrics)
	if err != nil {
		log.Fatalf("Failed to initialize scanner: %v", err)
	}

	scanner.Run()
}
