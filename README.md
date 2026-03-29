# Netmon

Netmon is a performant, containerized network monitoring tool that periodically scans IP subnets, monitors host reachability and latency via ICMP, and visualizes the results.

## Features

- Fast asynchronous ICMP monitoring for active hosts.
- Periodic background discovery sweeps for new hosts.
- Time-series metric storage using InfluxDB v2.
- Fully automated Grafana dashboard provisioning.
- Secure deployment via Docker Compose with least-privilege containers.

## Getting Started

1. Clone the repository.
2. Formally set up your credentials by copying `.env.example` to `.env` and filling in secure passwords.
3. Configure your network topology by copying `src/config.example.yaml` to `src/config.yaml` and adding your target IP ranges.
4. Start the stack:
   ```bash
   docker compose up -d --build
   ```

## Services

- Grafana: `http://localhost:3000` (Login using the credentials defined in your `.env` file)
- InfluxDB: `http://localhost:8086`

## Architecture

The stack consists of three core containers:
1. `netmon-service`: A Go daemon that executes concurrent ICMP pings. It operates securely by dropping all Linux capabilities except `CAP_NET_RAW`.
2. `influxdb`: A time-series database optimized for storing ping latency (`rtt_ms`) and reachability state (`up`).
3. `grafana`: A visualization frontend automatically provisioned with a highly scalable master dashboard showing network KPIs and historical host state.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
