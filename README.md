# Prometheus Service Monitoring Agent

This Go-based monitoring agent checks the status of systemd services and exposes their metrics for Prometheus.

## Features

- Monitors multiple systemd services.
- Reads the list of services from an environment file (`.env`).
- Exposes metrics on port `8888` for Prometheus.

## Requirements

- Linux-based system with `systemd`.
- Go 1.18+ installed.
- Prometheus for scraping metrics.

## Installation

### 1. Clone the repository

```sh
git clone https://github.com/yourusername/service-monitor.git
cd service-monitor
```

### 2. Create and edit the .env file

The application reads the list of services from an environment file. Create a .env file in the same directory as the program:

SERVICES=xrdp,nginx,postgresql,sshd

#### Explanation of `.env` variables:

| Variable   | Description                                                                                  |
| ---------- | -------------------------------------------------------------------------------------------- |
| `SERVICES` | A comma-separated list of systemd services to monitor. Example: `xrdp,nginx,postgresql,sshd` |

### 3. Install dependencies

```bash
go get github.com/joho/godotenv
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promhttp
```

### 4. Build and run the application

```bash
go run main.go
```

or

```bash
go build -o service-monitor
./service-monitor
```

Prometheus Integration

Add the following job to your Prometheus configuration (prometheus.yml):

```yaml
scrape_configs:
  - job_name: "service-monitor"
    static_configs:
      - targets: ["localhost:8888"]
```

Restart Prometheus for the changes to take effect.

Metrics Example

Once the service is running, visit http://localhost:8888/metrics in your browser or use:

```bash
curl http://localhost:8888/metrics
```

Example output:

```nginx
# HELP service_status_xrdp Current status of the xrdp service (1 = running, 0 = stopped)
# TYPE service_status_xrdp gauge
service_status_xrdp 1
# HELP service_status_nginx Current status of the nginx service (1 = running, 0 = stopped)
# TYPE service_status_nginx gauge
service_status_nginx 1
```

Running as a Systemd Service

To run the monitor as a background service:

Create a systemd service file:

```bash
sudo nano /etc/systemd/system/service-monitor.service
```

Add the following content:

```bash
[Unit]
Description=Prometheus Service Monitoring Agent
After=network.target

[Service]
ExecStart=/path/to/service-monitor
Restart=always
User=nobody
Group=nogroup
EnvironmentFile=/path/to/.env

[Install]
WantedBy=multi-user.target
Reload systemd and start the service:
sudo systemctl daemon-reload
sudo systemctl enable service-monitor
sudo systemctl start service-monitor
```

Check service status:

```bash
sudo systemctl status service-monitor
```

## License

MIT License.
