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
	var targets []string
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

		for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
			// Skip network address and broadcast address for IPv4
			if ip.To4() != nil {
				if ip.Equal(ipnet.IP) {
					continue
				}

				// Calculate broadcast address
				broadcast := make(net.IP, len(ip.To4()))
				for i := range ip.To4() {
					broadcast[i] = ip.To4()[i] | ^ipnet.Mask[i]
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
