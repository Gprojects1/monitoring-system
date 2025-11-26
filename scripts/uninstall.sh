#!/bin/bash
# Скрипт удаления системы мониторинга

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

INSTALL_DIR="/opt/webapp"
LOG_DIR="/var/log/webapp"
SERVICE_USER="webapp"

echo -e "${YELLOW}=== Удаление системы мониторинга ===${NC}"

if [[ $EUID -ne 0 ]]; then
   echo -e "${RED}Ошибка: Скрипт должен быть запущен с правами root${NC}"
   exit 1
fi

echo -e "${YELLOW}[1/5] Остановка сервисов...${NC}"
systemctl stop monitor.timer || true
systemctl stop monitor.service || true
systemctl stop webapp.service || true
echo -e "${GREEN}✓ Сервисы остановлены${NC}"

echo -e "${YELLOW}[2/5] Отключение автозапуска...${NC}"
systemctl disable monitor.timer || true
systemctl disable webapp.service || true
echo -e "${GREEN}✓ Автозапуск отключен${NC}"

echo -e "${YELLOW}[3/5] Удаление systemd файлов...${NC}"
rm -f /etc/systemd/system/webapp.service
rm -f /etc/systemd/system/monitor.service
rm -f /etc/systemd/system/monitor.timer
systemctl daemon-reload
echo -e "${GREEN}✓ Systemd файлы удалены${NC}"

echo -e "${YELLOW}[4/5] Удаление файлов приложения...${NC}"
rm -rf ${INSTALL_DIR}
echo -e "${GREEN}✓ Файлы приложения удалены${NC}"

echo -e "${YELLOW}[5/5] Удаление пользователя...${NC}"
if id -u ${SERVICE_USER} &> /dev/null; then
    userdel ${SERVICE_USER}
    echo -e "${GREEN}✓ Пользователь удален${NC}"
fi

echo -e "\n${GREEN}=== Удаление завершено ===${NC}"
echo -e "${YELLOW}Примечание: Логи сохранены в ${LOG_DIR}${NC}"
echo -e "Для полного удаления выполните: ${YELLOW}rm -rf ${LOG_DIR}${NC}"
