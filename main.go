package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// serviceMetrics stores Prometheus gauges for each monitored service
var serviceMetrics = make(map[string]prometheus.Gauge)

// loadServicesFromEnv loads the list of services from the .env file
func loadServicesFromEnv() []string {
    // Try to load .env file, but don't fail if it doesn't exist
    _ = godotenv.Load(".env")
    
    // Get the list of services from the environment variable
    services := os.Getenv("SERVICES")
    if services == "" {
        log.Println("SERVICES environment variable is not set")
        return nil
    }
    
    return strings.Split(services, ",")
}

// checkServiceStatus checks if a service is running and updates its metric
func checkServiceStatus(service string) {
	cmd := exec.Command("pgrep", "-x", service)
	err := cmd.Run()
	
	if err == nil {
		log.Printf("%s is running", service)
		serviceMetrics[service].Set(1)
	} else {
		log.Printf("%s is stopped", service)
		serviceMetrics[service].Set(0)
	}
}

// monitorServices periodically checks the status of all services
func monitorServices(services []string) {
	for {
		for _, service := range services {
			checkServiceStatus(service)
		}
		time.Sleep(10 * time.Second)
	}
}

func main() {
	// Load the list of services from .env
	services := loadServicesFromEnv()
	if services == nil || len(services) == 0 {
		log.Fatal("No services to monitor")
	}
	
	// Register metrics for each service
	for _, service := range services {
		metric := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "service_status_" + service,
			Help: "Current status of " + service + " (1 = running, 0 = stopped)",
		})
		prometheus.MustRegister(metric)
		serviceMetrics[service] = metric
	}

	// Start monitoring in a separate goroutine
	go monitorServices(services)

	// Start HTTP server for Prometheus metrics
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Starting Prometheus service monitoring agent on :8888, monitoring services: %s", strings.Join(services, ", "))
	log.Fatal(http.ListenAndServe(":8888", nil))
}