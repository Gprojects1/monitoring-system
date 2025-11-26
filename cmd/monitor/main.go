package main

import (
	"log"
	"os"

	"github.com/Gprojects1/monitoring-system/internal/config"
	"github.com/Gprojects1/monitoring-system/internal/logger"
	"github.com/Gprojects1/monitoring-system/internal/monitor"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadMonitorConfig("configs/monitor.yaml")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Инициализация логгера
	monLogger, err := logger.New(cfg.LogFile, cfg.LogLevel)
	if err != nil {
		log.Fatalf("Ошибка инициализации логгера: %v", err)
	}
	defer monLogger.Close()

	// Создание монитора
	mon := monitor.New(cfg, monLogger)

	// Запуск проверки
	monLogger.Info("========================================", nil)
	monLogger.Info("Начало проверки мониторинга", nil)
	monLogger.Info("========================================", nil)

	exitCode := mon.Run()

	if exitCode == 0 {
		monLogger.Info("Мониторинг завершен успешно", nil)
	} else {
		monLogger.Error("Мониторинг завершен с ошибками", map[string]interface{}{
			"exit_code": exitCode,
		})
	}

	os.Exit(exitCode)
}
