package main

import (
	"reflect"
	"testing"
)

func TestGenerateTargets(t *testing.T) {
	tests := []struct {
		name     string
		networks []string
		expected []string
		wantErr  bool
	}{
		{
			name:     "Single IP",
			networks: []string{"192.168.1.1"},
			expected: []string{"192.168.1.1"},
			wantErr:  false,
		},
		{
			name:     "Slash 32",
			networks: []string{"10.0.0.5/32"},
			expected: []string{"10.0.0.5"},
			wantErr:  false,
		},
		{
			name:     "Slash 31",
			networks: []string{"192.168.1.0/31"},
			expected: []string{"192.168.1.0", "192.168.1.1"},
			wantErr:  false,
		},
		{
			name:     "Slash 30",
			networks: []string{"192.168.1.0/30"}, // 192.168.1.0 (net), 192.168.1.1, 192.168.1.2, 192.168.1.3 (bcast)
			expected: []string{"192.168.1.1", "192.168.1.2"},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Networks: tt.networks,
			}
			targets, err := cfg.GenerateTargets()
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateTargets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(targets, tt.expected) {
				t.Errorf("GenerateTargets() = %v, expected %v", targets, tt.expected)
			}
		})
	}
}
