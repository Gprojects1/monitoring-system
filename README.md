# monitoring-system
# Система мониторинга веб-приложения

Автоматическая система мониторинга на Go с перезапуском при сбоях.

## Установка

# Установить зависимости
make deps

# Собрать проект
make build

# Установить
sudo bash scripts/install.sh

## Использование

# Проверить работу
curl http://localhost:5000/

# Посмотреть логи приложения
sudo tail -f /var/log/webapp/app.log

# Посмотреть логи мониторинга
sudo tail -f /var/log/webapp/monitor.log

# Статус сервисов
sudo systemctl status webapp.service
sudo systemctl status monitor.timer

## Конфигурация

Приложение: configs/webapp.yaml
Мониторинг: configs/monitor.yaml

Интервал мониторинга: каждые 30 секунд

## Управление

# Остановить
sudo systemctl stop webapp.service
sudo systemctl stop monitor.timer

# Запустить
sudo systemctl start webapp.service
sudo systemctl start monitor.timer

# Перезапустить
sudo systemctl restart webapp.service

## API Endpoints

GET /          - Hello World
GET /health    - Health check
GET /status    - Расширенный статус

Пример:
curl http://localhost:5000/health

## Удаление

# Обычное удаление (логи остаются)
sudo bash scripts/uninstall.sh

# Полное удаление (включая логи)
sudo bash scripts/uninstall.sh --full

## Структура проекта

monitoring-system/
├── cmd/          - Исходный код приложений
├── internal/     - Внутренние пакеты
├── configs/      - Конфигурационные файлы
├── systemd/      - Systemd сервисы
└── scripts/      - Скрипты установки/удаления

## Команды Make

make build       - Сборка проекта
make build-prod  - Сборка с оптимизацией
make deps        - Установка зависимостей
make clean       - Очистка
make fmt         - Форматирование кода
make lint        - Проверка кода
make run-webapp  - Локальный запуск приложения
make run-monitor - Локальный запуск монитора

## Логи

Логи приложения: /var/log/webapp/app.log
Логи мониторинга: /var/log/webapp/monitor.log

# Посмотреть логи через journald
sudo journalctl -u webapp -f
sudo journalctl -u monitor.service -f

## Troubleshooting

# Проверить порт
sudo netstat -tulpn | grep 5000

# Проверить статус
sudo systemctl status webapp.service

# Проверить логи
sudo journalctl -u webapp -n 50

# Восстановить права
sudo chown -R webapp:webapp /opt/webapp
sudo chown -R webapp:webapp /var/log/webapp

## Лицензия

MIT
