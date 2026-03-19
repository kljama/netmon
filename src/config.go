package main

import (
	"fmt"
	"net"
	"net/netip"
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
	// ⚡ Bolt: Fast path capacity calculation to pre-allocate targets
	// This reduces allocations and memory copying by predicting size
	var totalCap int
	for _, network := range c.Networks {
		if prefix, err := netip.ParsePrefix(network); err == nil {
			if prefix.Addr().Is4() {
				count := 1 << (32 - prefix.Bits())
				if count > 2 {
					count -= 2 // Skip network and broadcast
				} else if count == 2 {
					count = 2 // /31 network
				} else {
					count = 1 // /32 network
				}
				totalCap += count
			} else {
				// Wild guess for IPv6 capacity (we don't want to over-allocate massively)
				totalCap += 256
			}
		} else if parsedAddr, err := netip.ParseAddr(network); err == nil && parsedAddr.IsValid() {
			totalCap += 1
		} else {
			// Fallback if netip fails (e.g. extremely old format or legacy code might use net.ParseCIDR instead)
			totalCap += 1
		}
	}

	targets := make([]string, 0, totalCap)

	for _, network := range c.Networks {
		prefix, err := netip.ParsePrefix(network)
		if err != nil {
			// Try as a single IP
			if parsedAddr, err := netip.ParseAddr(network); err == nil && parsedAddr.IsValid() {
				targets = append(targets, parsedAddr.String())
				continue
			}

			// Legacy fallback for net.ParseIP compatibility
			if parsedIP := net.ParseIP(network); parsedIP != nil {
				targets = append(targets, network)
				continue
			}
			return nil, fmt.Errorf("invalid network or IP %q: %w", network, err)
		}

		if prefix.Addr().Is4() {
			// ⚡ Bolt: Fast path for IPv4 using netip
			count := 1 << (32 - prefix.Bits())
			if count > 2 {
				count -= 2
			} else if count == 2 {
				count = 2
			} else {
				count = 1
			}

			addr := prefix.Masked().Addr()

			// For standard subnets (> /31), skip the network address
			if count > 1 && prefix.Bits() < 31 {
				addr = addr.Next()
			}

			for i := 0; i < count; i++ {
				targets = append(targets, addr.String())
				addr = addr.Next()
			}
		} else {
			// Legacy loop for IPv6 or unsupported configurations using net
			ip, ipnet, err := net.ParseCIDR(network)
			if err != nil {
				return nil, fmt.Errorf("invalid network %q: %w", network, err)
			}

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
