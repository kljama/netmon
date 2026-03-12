package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type MetricsClient struct {
	client   influxdb2.Client
	writeAPI api.WriteAPIBlocking
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
	writeAPI := client.WriteAPIBlocking(org, bucket)

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

	err := m.writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		log.Printf("Failed to write point for %s: %v", target, err)
		return err
	}
	return nil
}

func (m *MetricsClient) Close() {
	m.client.Close()
}
