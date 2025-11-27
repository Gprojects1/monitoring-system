# Система мониторинга веб-приложения

Автоматическая система мониторинга на Go с перезапуском при сбоях.

# Установка

## Установить зависимости
```bash
make deps
```
## Собрать проект
```bash
make build

## Установить
```bash
sudo bash scripts/install.sh
```
# Использование

## Проверить работу
```bash
curl http://localhost:5000/
```
## Посмотреть логи приложения
```bash
sudo tail -f /var/log/webapp/app.log
```
## Посмотреть логи мониторинга
```bash
sudo tail -f /var/log/webapp/monitor.log
```
## Статус сервисов
```bash
sudo systemctl status webapp.service
sudo systemctl status monitor.timer
```
# Конфигурация

Приложение: configs/webapp.yaml
Мониторинг: configs/monitor.yaml

Интервал мониторинга: каждые 30 секунд

# Управление

## Остановить
```bash
sudo systemctl stop webapp.service
sudo systemctl stop monitor.timer
```
## Запустить
```bash
sudo systemctl start webapp.service
sudo systemctl start monitor.timer
```
## Перезапустить
```bash
sudo systemctl restart webapp.service
```
# API Endpoints

GET /          - Hello World
GET /health    - Health check
GET /status    - Расширенный статус

Пример:
```bash
curl http://localhost:5000/health
```
# Удаление

## Обычное удаление (логи остаются)
```bash
sudo bash scripts/uninstall.sh
```
## Полное удаление (включая логи)
```bash
sudo bash scripts/uninstall.sh --full
```
# Структура проекта

monitoring-system/
├── cmd/          - Исходный код приложений
├── internal/     - Внутренние пакеты
├── configs/      - Конфигурационные файлы
├── systemd/      - Systemd сервисы
└── scripts/      - Скрипты установки/удаления

# Команды Make

make build       - Сборка проекта
make build-prod  - Сборка с оптимизацией
make deps        - Установка зависимостей
make clean       - Очистка
make fmt         - Форматирование кода
make lint        - Проверка кода
make run-webapp  - Локальный запуск приложения
make run-monitor - Локальный запуск монитора

# Логи

Логи приложения: /var/log/webapp/app.log
Логи мониторинга: /var/log/webapp/monitor.log

## Посмотреть логи через journald
```bash
sudo journalctl -u webapp -f
sudo journalctl -u monitor.service -f
```
# Диагностика

## Проверить порт
```bash
sudo netstat -tulpn | grep 5000
```
## Проверить статус
```bash
sudo systemctl status webapp.service
```
## Проверить логи
```bash
sudo journalctl -u webapp -n 50
```


