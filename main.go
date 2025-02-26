package main

import (
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Define a Prometheus gauge metric
var xrdpStatus = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "xrdp_service_status",
	Help: "Current status of the xrdp service (1 = running, 0 = stopped)",
})

// Function to check if xrdp is running
func checkXrdpStatus() {
	cmd := exec.Command("systemctl", "is-active", "--quiet", "xrdp")
	err := cmd.Run()
	if err == nil {
		xrdpStatus.Set(1) // Service is running
	} else {
		xrdpStatus.Set(0) // Service is not running
	}
}

func startMonitoring() {
	for {
		checkXrdpStatus()
		time.Sleep(10 * time.Second) // Adjust the interval as needed
	}
}

func main() {
	// Register metric
	prometheus.MustRegister(xrdpStatus)

	// Start monitoring in a goroutine
	go startMonitoring()

	// Expose metrics endpoint
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Starting Prometheus xrdp monitoring agent on :8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
