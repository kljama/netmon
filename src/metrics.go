package main

import (
	"fmt"
	"log"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type MetricsClient struct {
	client   influxdb2.Client
	writeAPI api.WriteAPI
	org      string
	bucket   string
}

func NewMetricsClient() (*MetricsClient, error) {
	url := os.Getenv("INFLUX_URL")
	token := os.Getenv("INFLUX_TOKEN")
	org := os.Getenv("INFLUX_ORG")
	bucket := os.Getenv("INFLUX_BUCKET")

	if url == "" || token == "" || org == "" || bucket == "" {
		return nil, fmt.Errorf("missing InfluxDB environment variables (INFLUX_URL, INFLUX_TOKEN, INFLUX_ORG, INFLUX_BUCKET)")
	}

	client := influxdb2.NewClient(url, token)
	writeAPI := client.WriteAPI(org, bucket)

	// Listen to background write errors
	go func() {
		for err := range writeAPI.Errors() {
			log.Printf("Background write error: %v", err)
		}
	}()

	return &MetricsClient{
		client:   client,
		writeAPI: writeAPI,
		org:      org,
		bucket:   bucket,
	}, nil
}

func (m *MetricsClient) RecordPing(target string, rtt time.Duration, up bool) error {
	p := influxdb2.NewPointWithMeasurement("icmp_ping").
		AddTag("target", target).
		AddField("rtt_ms", float64(rtt.Milliseconds())).
		AddField("up", up).
		SetTime(time.Now())

	// WritePoint is non-blocking and adds the point to the buffer
	m.writeAPI.WritePoint(p)
	return nil
}

func (m *MetricsClient) Close() {
	m.writeAPI.Flush()
	m.client.Close()
}
