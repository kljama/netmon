package main

import (
	"testing"
)

func BenchmarkGenerateTargets_IPv4_24(b *testing.B) {
	cfg := &Config{
		Networks: []string{"192.168.1.0/24"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cfg.GenerateTargets()
	}
}

func BenchmarkGenerateTargets_IPv4_16(b *testing.B) {
	cfg := &Config{
		Networks: []string{"10.0.0.0/16"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cfg.GenerateTargets()
	}
}

func BenchmarkGenerateTargets_IPv6_120(b *testing.B) {
	cfg := &Config{
		Networks: []string{"2001:db8::/120"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cfg.GenerateTargets()
	}
}
