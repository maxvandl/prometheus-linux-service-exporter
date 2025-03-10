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

// Карта метрик для каждого сервиса
var serviceMetrics = make(map[string]prometheus.Gauge)

// Функция для проверки состояния сервиса
func checkServiceStatus(service string) {
	cmd := exec.Command("systemctl", "is-active", "--quiet", service)
	err := cmd.Run()
	if err == nil {
		serviceMetrics[service].Set(1) // Сервис работает
	} else {
		serviceMetrics[service].Set(0) // Сервис остановлен
	}
}

// Запуск мониторинга в цикле
func startMonitoring(services []string) {
	for {
		for _, service := range services {
			checkServiceStatus(service)
		}
		time.Sleep(10 * time.Second) // Интервал обновления метрик
	}
}

func main() {
	// Загружаем переменные из .env (если файл есть)
	_ = godotenv.Load()

	// Получаем список сервисов из переменной окружения
	servicesStr := os.Getenv("SERVICES")
	if servicesStr == "" {
		log.Fatal("SERVICES environment variable is not set")
	}
	services := strings.Split(servicesStr, ",")

	// Регистрация метрик для каждого сервиса
	for _, service := range services {
		service = strings.TrimSpace(service) // Убираем лишние пробелы
		serviceMetrics[service] = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "service_status_" + service,
			Help: "Current status of the " + service + " service (1 = running, 0 = stopped)",
		})
		prometheus.MustRegister(serviceMetrics[service])
	}

	// Запуск мониторинга в отдельной горутине
	go startMonitoring(services)

	// Экспорт метрик
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Starting Prometheus monitoring agent on :8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
