package main

import (
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/prometheus-community/pro-bing"
)

type Scanner struct {
	targets     []string
	activeIPs   sync.Map
	activeCount atomic.Int64
	cfg         *Config
	metrics     *MetricsClient
}

func NewScanner(cfg *Config, metrics *MetricsClient) (*Scanner, error) {
	targets, err := cfg.GenerateTargets()
	if err != nil {
		return nil, err
	}

	return &Scanner{
		targets: targets,
		cfg:     cfg,
		metrics: metrics,
	}, nil
}

func (s *Scanner) Run() {
	log.Printf("Starting scanner for %d targets...", len(s.targets))
	ticketChan := make(chan struct{}, runtime.NumCPU()*10) // Concurrency limit

	discoveryTicker := time.NewTicker(s.cfg.DiscoveryInterval)
	defer discoveryTicker.Stop()

	scanTicker := time.NewTicker(s.cfg.ScanInterval)
	defer scanTicker.Stop()

	// Initial discovery scan
	s.runDiscovery(ticketChan)

	for {
		select {
		case <-discoveryTicker.C:
			s.runDiscovery(ticketChan)
		case <-scanTicker.C:
			s.runMonitoring(ticketChan)
		}
	}
}

func (s *Scanner) runDiscovery(ticketChan chan struct{}) {
	log.Printf("Running discovery sweep on all %d targets...", len(s.targets))
	var wg sync.WaitGroup
	for _, target := range s.targets {
		// ⚡ Bolt: Fast-path early return.
		// Since activeIPs is append-only in this application (hosts are never removed),
		// we can safely check if the host is already known before acquiring a concurrency
		// ticket and spawning a goroutine. This avoids thousands of unnecessary goroutines
		// and lock acquisitions on subsequent sweeps.
		if _, exists := s.activeIPs.Load(target); exists {
			continue
		}

		wg.Add(1)
		targetCopy := target
		ticketChan <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-ticketChan }()

			// Only discover if we don't already know they are active
			// (Double-check in case another goroutine just discovered it)
			if _, exists := s.activeIPs.Load(targetCopy); !exists {
				up, _, _ := s.pingInternal(targetCopy)
				if up {
					log.Printf("Discovered new active host: %s", targetCopy)
					_, loaded := s.activeIPs.LoadOrStore(targetCopy, true)
					if !loaded {
						s.activeCount.Add(1)
					}
				}
			}
		}()
	}
	wg.Wait()

	log.Printf("Discovery sweep complete. Tracking %d active hosts.", s.activeCount.Load())
}

func (s *Scanner) runMonitoring(ticketChan chan struct{}) {
	var wg sync.WaitGroup
	s.activeIPs.Range(func(key, value interface{}) bool {
		target := key.(string)
		wg.Add(1)
		ticketChan <- struct{}{}
		go func(t string) {
			defer wg.Done()
			defer func() { <-ticketChan }()

			up, rtt, err := s.pingInternal(t)

			// Wait before persisting to avoid overwhelming if we needed to, but handled by concurrency
			if err == nil {
				metricErr := s.metrics.RecordPing(t, rtt, up)
				if metricErr != nil {
					log.Printf("Failed to record metrics for %s: %v", t, metricErr)
				}
			}
			// NOTE: We do NOT remove down hosts from s.activeIPs as per user request.
			// They will remain tracked to report 'up=false' to Grafana.
		}(target)
		return true
	})
	wg.Wait()
}

func (s *Scanner) pingInternal(target string) (bool, time.Duration, error) {
	pinger, err := probing.NewPinger(target)
	if err != nil {
		log.Printf("Failed to create pinger for %s: %v", target, err)
		return false, 0, err
	}

	// We run as non-root with CAP_NET_RAW which allows us to open raw sockets.
	// Privileged=true tells the library to use raw sockets rather than UDP pings.
	pinger.SetPrivileged(true)
	pinger.Count = 1
	pinger.Timeout = s.cfg.Timeout

	err = pinger.Run()
	if err != nil {
		log.Printf("Failed to ping %s: %v", target, err)
		return false, 0, err
	}

	stats := pinger.Statistics()
	if stats.PacketsRecv > 0 {
		return true, stats.AvgRtt, nil
	}
	return false, 0, nil
}
