package main

import (
	"testing"
)

func BenchmarkGenerateTargets(b *testing.B) {
	cfg := &Config{
		Networks: []string{"10.0.0.0/16"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cfg.GenerateTargets()
	}
}
