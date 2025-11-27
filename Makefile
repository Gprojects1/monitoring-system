# Makefile для сборки проекта

.PHONY: all build clean install uninstall test

# Переменные
BINARY_DIR=bin
WEBAPP_BINARY=$(BINARY_DIR)/webapp
MONITOR_BINARY=$(BINARY_DIR)/monitor
INSTALL_DIR=/opt/webapp

# Сборка всех бинарников
all: build

# Сборка проекта
build:
	@echo "Сборка веб-приложения..."
	@mkdir -p $(BINARY_DIR)
	@go build -o $(WEBAPP_BINARY) ./cmd/webapp
	@echo "✓ Веб-приложение собрано: $(WEBAPP_BINARY)"
	
	@echo "Сборка монитора..."
	@go build -o $(MONITOR_BINARY) ./cmd/monitor
	@echo "✓ Монитор собран: $(MONITOR_BINARY)"

# Сборка с оптимизацией для продакшена
build-prod:
	@echo "Сборка для продакшена..."
	@mkdir -p $(BINARY_DIR)
	@go build -ldflags="-s -w" -o $(WEBAPP_BINARY) ./cmd/webapp
	@go build -ldflags="-s -w" -o $(MONITOR_BINARY) ./cmd/monitor
	@echo "✓ Сборка завершена"

# Установка зависимостей
deps:
	@echo "Установка зависимостей..."
	@go mod download
	@go mod verify
	@echo "✓ Зависимости установлены"

# Запуск тестов
test:
	@echo "Запуск тестов..."
	@go test -v ./...

# Очистка
clean:
	@echo "Очистка..."
	@rm -rf $(BINARY_DIR)
	@go clean
	@echo "✓ Очистка завершена"

# Форматирование кода
fmt:
	@echo "Форматирование кода..."
	@go fmt ./...
	@echo "✓ Код отформатирован"

# Проверка кода
lint:
	@echo "Проверка кода..."
	@go vet ./...
	@echo "✓ Проверка завершена"

# Запуск веб-приложения локально
run-webapp: build
	@echo "Запуск веб-приложения..."
	@./$(WEBAPP_BINARY)

# Запуск монитора локально
run-monitor: build
	@echo "Запуск монитора..."
	@./$(MONITOR_BINARY)

# Показать справку
help:
	@echo "Доступные команды:"
	@echo "  make build       - Сборка проекта"
	@echo "  make build-prod  - Сборка для продакшена"
	@echo "  make deps        - Установка зависимостей"
	@echo "  make test        - Запуск тестов"
	@echo "  make clean       - Очистка"
	@echo "  make fmt         - Форматирование кода"
	@echo "  make lint        - Проверка кода"
	@echo "  make run-webapp  - Запуск веб-приложения"
	@echo "  make run-monitor - Запуск монитора"
