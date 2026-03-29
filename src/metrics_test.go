package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func BenchmarkRecordPing(b *testing.B) {
	// Create a dummy HTTP server to simulate InfluxDB
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent) // InfluxDB typically returns 204 No Content for successful writes
	}))
	defer ts.Close()

	os.Setenv("INFLUX_URL", ts.URL)
	os.Setenv("INFLUX_TOKEN", "dummy-token")
	os.Setenv("INFLUX_ORG", "dummy-org")
	os.Setenv("INFLUX_BUCKET", "dummy-bucket")
	defer func() {
		os.Unsetenv("INFLUX_URL")
		os.Unsetenv("INFLUX_TOKEN")
		os.Unsetenv("INFLUX_ORG")
		os.Unsetenv("INFLUX_BUCKET")
	}()

	client, err := NewMetricsClient()
	if err != nil {
		b.Fatalf("Failed to create metrics client: %v", err)
	}
	defer client.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := client.RecordPing("127.0.0.1", time.Millisecond*10, true)
		if err != nil {
			b.Fatalf("RecordPing failed: %v", err)
		}
	}
}
