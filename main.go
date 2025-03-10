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

// Map для хранения метрик сервисов
var serviceMetrics = make(map[string]prometheus.Gauge)

// Функция для загрузки переменных окружения из .env
func loadEnv() []string {
	err := godotenv.Load(".env") // Загружаем файл .env


	services := os.Getenv("SERVICES") // Получаем список сервисов


	return strings.Split(services, ",") // Разбиваем строку в массив
}

// Функция для проверки статуса сервиса
func checkServiceStatus(service string) {
	cmd := exec.Command("pgrep", "-x", service)
	err := cmd.Run()
	if err == nil {
		log.Printf("%s is running\n", service)
		serviceMetrics[service].Set(1)
	} else {
		log.Printf("%s is stopped\n", service)
		serviceMetrics[service].Set(0)
	}
}

// Мониторинг всех сервисов
func startMonitoring(services []string) {
	for {
		for _, service := range services {
			checkServiceStatus(service)
		}
		time.Sleep(10 * time.Second)
	}
}

func main() {
	// Загружаем список сервисов из .env
	services := loadEnv()

	// Регистрируем метрики для каждого сервиса
	for _, service := range services {
		metric := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "service_status_" + service,
			Help: "Current status of " + service + " (1 = running, 0 = stopped)",
		})
		prometheus.MustRegister(metric)
		serviceMetrics[service] = metric
	}

	// Запускаем мониторинг в отдельной горутине
	go startMonitoring(services)

	// Запускаем HTTP-сервер для Prometheus
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Starting Prometheus service monitoring agent on :8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
