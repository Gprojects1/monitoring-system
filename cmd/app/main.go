package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Gprojects1/monitoring-system/internal/config"
	"github.com/Gprojects1/monitoring-system/internal/logger"
)

// HealthResponse представляет ответ health check эндпоинта
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
}

// StatusResponse представляет расширенный статус приложения
type StatusResponse struct {
	Status    string    `json:"status"`
	Uptime    string    `json:"uptime"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
}

var (
	startTime time.Time
	cfg       *config.WebAppConfig
	appLogger *logger.Logger
)

func init() {
	startTime = time.Now()
}

// helloHandler обрабатывает главный эндпоинт
func helloHandler(w http.ResponseWriter, r *http.Request) {
	appLogger.Info("Запрос на главную страницу", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
		"remote": r.RemoteAddr,
	})

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Hello World!")
}

// healthHandler обрабатывает health check запросы
func healthHandler(w http.ResponseWriter, r *http.Request) {
	appLogger.Debug("Health check запрос", map[string]interface{}{
		"remote": r.RemoteAddr,
	})

	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Service:   "hello-world-app",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// statusHandler обрабатывает запросы расширенного статуса
func statusHandler(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(startTime)

	response := StatusResponse{
		Status:    "running",
		Uptime:    uptime.String(),
		Version:   "1.0.0",
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// loggingMiddleware логирует все HTTP запросы
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		duration := time.Since(start)

		appLogger.Info("HTTP запрос обработан", map[string]interface{}{
			"method":   r.Method,
			"path":     r.URL.Path,
			"duration": duration.String(),
			"remote":   r.RemoteAddr,
		})
	}
}

// setupRoutes настраивает маршруты HTTP сервера
func setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Регистрация обработчиков с middleware
	mux.HandleFunc("/", loggingMiddleware(helloHandler))
	mux.HandleFunc("/health", loggingMiddleware(healthHandler))
	mux.HandleFunc("/status", loggingMiddleware(statusHandler))

	return mux
}

// gracefulShutdown обрабатывает корректное завершение работы
func gracefulShutdown(server *http.Server) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	appLogger.Info("Получен сигнал завершения работы", nil)

	// Даем серверу время завершить обработку запросов
	appLogger.Info("Завершение работы сервера...", nil)
	server.Close()
}

func main() {
	// Загрузка конфигурации
	var err error
	cfg, err = config.LoadWebAppConfig("configs/webapp.yaml")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Инициализация логгера
	appLogger, err = logger.New(cfg.LogFile, cfg.LogLevel)
	if err != nil {
		log.Fatalf("Ошибка инициализации логгера: %v", err)
	}
	defer appLogger.Close()

	appLogger.Info("Запуск веб-приложения", map[string]interface{}{
		"host": cfg.Host,
		"port": cfg.Port,
	})

	// Настройка HTTP сервера
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      setupRoutes(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запуск горутины для graceful shutdown
	go gracefulShutdown(server)

	// Запуск сервера
	appLogger.Info(fmt.Sprintf("Сервер запущен на %s", addr), nil)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		appLogger.Error("Ошибка при запуске сервера", map[string]interface{}{
			"error": err.Error(),
		})
		os.Exit(1)
	}

	appLogger.Info("Сервер остановлен", nil)
}
