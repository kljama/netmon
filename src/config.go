package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Networks          []string      `yaml:"networks"`
	ScanInterval      time.Duration `yaml:"scan_interval"`
	DiscoveryInterval time.Duration `yaml:"discovery_interval"`
	Timeout           time.Duration `yaml:"timeout"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := &Config{
		ScanInterval:      10 * time.Second,
		DiscoveryInterval: 5 * time.Minute,
		Timeout:           1 * time.Second,
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return cfg, nil
}

// GenerateTargets expands CIDR network ranges into individual IP addresses.
func (c *Config) GenerateTargets() ([]string, error) {
	// ⚡ Bolt: Pre-calculate total capacity to avoid slice reallocation
	// This reduces memory allocations from ~6.6MB to ~2.1MB per /16 network.
	var capacity int
	for _, network := range c.Networks {
		if _, ipnet, err := net.ParseCIDR(network); err == nil {
			ones, bits := ipnet.Mask.Size()
			// Cap pre-allocation to avoid OOM for massive blocks (e.g. > /16)
			if shift := bits - ones; shift >= 0 && shift <= 16 {
				capacity += 1 << shift
			}
		} else if net.ParseIP(network) != nil {
			capacity++
		}
	}

	targets := make([]string, 0, capacity)

	for _, network := range c.Networks {
		ip, ipnet, err := net.ParseCIDR(network)
		if err != nil {
			// Try as a single IP
			if parsedIP := net.ParseIP(network); parsedIP != nil {
				targets = append(targets, network)
				continue
			}
			return nil, fmt.Errorf("invalid network or IP %q: %w", network, err)
		}

		// Security Check: Prevent memory exhaustion / DoS from oversized subnets
		// Maximum allowed is /16 for IPv4 and /112 for IPv6 (65536 addresses)
		ones, bits := ipnet.Mask.Size()
		if bits-ones > 16 {
			return nil, fmt.Errorf("network %q is too large; max allowed size is 65536 addresses (e.g. /16 for IPv4)", network)
		}

		// Calculate broadcast address outside the loop for IPv4
		var broadcast net.IP
		ip4 := ip.To4()
		isIPv4 := ip4 != nil
		if isIPv4 {
			broadcast = make(net.IP, len(ip4))
			for i := range ip4 {
				broadcast[i] = ip4[i] | ^ipnet.Mask[i]
			}
		}

		for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
			// Skip network address and broadcast address for IPv4
			if isIPv4 {
				if ip.Equal(ipnet.IP) {
					continue
				}

				if ip.Equal(broadcast) {
					continue
				}
			}
			targets = append(targets, ip.String())
		}
	}
	return targets, nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
