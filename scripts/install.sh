#!/bin/bash
# Скрипт автоматической установки системы мониторинга (Go версия)

set -e

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Конфигурация
INSTALL_DIR="/opt/webapp"
BIN_DIR="${INSTALL_DIR}/bin"
CONFIG_DIR="${INSTALL_DIR}/configs"
LOG_DIR="/var/log/webapp"
SERVICE_USER="webapp"
SERVICE_GROUP="webapp"

echo -e "${GREEN}=== Установка системы мониторинга (Go) ===${NC}"

# Проверка прав root
if [[ $EUID -ne 0 ]]; then
   echo -e "${RED}Ошибка: Скрипт должен быть запущен с правами root${NC}"
   exit 1
fi

# Проверка Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}Ошибка: Go не установлен${NC}"
    echo -e "${YELLOW}Установите Go: https://golang.org/doc/install${NC}"
    exit 1
fi

echo -e "${YELLOW}[1/8] Сборка проекта...${NC}"
make clean
make build-prod
echo -e "${GREEN}✓ Проект собран${NC}"

echo -e "${YELLOW}[2/8] Создание пользователя и группы...${NC}"
if ! id -u ${SERVICE_USER} &> /dev/null; then
    useradd --system --no-create-home --shell /bin/false ${SERVICE_USER}
    echo -e "${GREEN}✓ Пользователь ${SERVICE_USER} создан${NC}"
else
    echo -e "${YELLOW}✓ Пользователь ${SERVICE_USER} уже существует${NC}"
fi

echo -e "${YELLOW}[3/8] Создание директорий...${NC}"
mkdir -p ${BIN_DIR}
mkdir -p ${CONFIG_DIR}
mkdir -p ${LOG_DIR}
echo -e "${GREEN}✓ Директории созданы${NC}"

echo -e "${YELLOW}[4/8] Копирование файлов...${NC}"
cp bin/webapp ${BIN_DIR}/
cp bin/monitor ${BIN_DIR}/
cp configs/*.yaml ${CONFIG_DIR}/
chmod +x ${BIN_DIR}/webapp
chmod +x ${BIN_DIR}/monitor
echo -e "${GREEN}✓ Файлы скопированы${NC}"

echo -e "${YELLOW}[5/8] Настройка прав доступа...${NC}"
chown -R ${SERVICE_USER}:${SERVICE_GROUP} ${INSTALL_DIR}
chown -R ${SERVICE_USER}:${SERVICE_GROUP} ${LOG_DIR}
chmod -R 755 ${INSTALL_DIR}
echo -e "${GREEN}✓ Права доступа настроены${NC}"

echo -e "${YELLOW}[6/8] Установка systemd сервисов...${NC}"
cp systemd/webapp.service /etc/systemd/system/
cp systemd/monitor.service /etc/systemd/system/
cp systemd/monitor.timer /etc/systemd/system/
systemctl daemon-reload
echo -e "${GREEN}✓ Systemd сервисы установлены${NC}"

echo -e "${YELLOW}[7/8] Запуск сервисов...${NC}"
systemctl enable webapp.service
systemctl start webapp.service
sleep 2

systemctl enable monitor.timer
systemctl start monitor.timer
echo -e "${GREEN}✓ Сервисы запущены${NC}"

echo -e "${YELLOW}[8/8] Проверка работоспособности...${NC}"
sleep 2
if curl -f http://localhost:5000/health &> /dev/null; then
    echo -e "${GREEN}✓ Приложение работает корректно${NC}"
else
    echo -e "${RED}✗ Приложение не отвечает${NC}"
fi

# Проверка статуса
echo -e "\n${GREEN}=== Статус сервисов ===${NC}"
systemctl status webapp.service --no-pager || true
echo ""
systemctl status monitor.timer --no-pager || true

echo -e "\n${GREEN}=== Установка завершена успешно! ===${NC}"
echo -e "Веб-приложение доступно по адресу: ${YELLOW}http://localhost:5000${NC}"
echo -e "Логи приложения: ${YELLOW}${LOG_DIR}/app.log${NC}"
echo -e "Логи мониторинга: ${YELLOW}${LOG_DIR}/monitor.log${NC}"
echo -e "\nПолезные команды:"
echo -e "  ${YELLOW}systemctl status webapp.service${NC}   - статус приложения"
echo -e "  ${YELLOW}systemctl status monitor.timer${NC}    - статус мониторинга"
echo -e "  ${YELLOW}journalctl -u webapp -f${NC}           - логи приложения"
echo -e "  ${YELLOW}tail -f ${LOG_DIR}/monitor.log${NC}    - логи мониторинга"
