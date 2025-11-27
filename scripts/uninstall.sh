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

# Определяем директорию проекта
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN_DIR="${PROJECT_DIR}/bin"

# Опция полного удаления (включая логи)
REMOVE_LOGS=false
if [[ "$1" == "--full" ]] || [[ "$1" == "-f" ]]; then
    REMOVE_LOGS=true
fi

echo -e "${YELLOW}=== Удаление системы мониторинга ===${NC}"

if [[ $EUID -ne 0 ]]; then
   echo -e "${RED}Ошибка: Скрипт должен быть запущен с правами root${NC}"
   exit 1
fi

echo -e "${YELLOW}[1/7] Остановка сервисов...${NC}"
systemctl stop monitor.timer || true
systemctl stop monitor.service || true
systemctl stop webapp.service || true
# Убить процессы 
pkill -9 webapp 2>/dev/null || true
pkill -9 monitor 2>/dev/null || true
echo -e "${GREEN}✓ Сервисы остановлены${NC}"

echo -e "${YELLOW}[2/7] Отключение автозапуска...${NC}"
systemctl disable monitor.timer || true
systemctl disable webapp.service || true
echo -e "${GREEN}✓ Автозапуск отключен${NC}"

echo -e "${YELLOW}[3/7] Удаление systemd файлов...${NC}"
rm -f /etc/systemd/system/webapp.service
rm -f /etc/systemd/system/monitor.service
rm -f /etc/systemd/system/monitor.timer
systemctl daemon-reload
echo -e "${GREEN}✓ Systemd файлы удалены${NC}"

echo -e "${YELLOW}[4/7] Удаление файлов приложения...${NC}"
rm -rf ${INSTALL_DIR}
echo -e "${GREEN}✓ Файлы приложения удалены${NC}"

echo -e "${YELLOW}[5/7] Удаление локальных бинарников...${NC}"
if [ -d "${BIN_DIR}" ]; then
    rm -rf ${BIN_DIR}
    echo -e "${GREEN}✓ Локальные бинарники удалены (${BIN_DIR})${NC}"
else
    echo -e "${YELLOW}✓ Директория bin/ не найдена${NC}"
fi

echo -e "${YELLOW}[6/7] Удаление пользователя...${NC}"
if id -u ${SERVICE_USER} &> /dev/null; then
    userdel ${SERVICE_USER}
    echo -e "${GREEN}✓ Пользователь удален${NC}"
fi

echo -e "${YELLOW}[7/7] Обработка логов...${NC}"
if [ "$REMOVE_LOGS" = true ]; then
    if [ -d "${LOG_DIR}" ]; then
        rm -rf ${LOG_DIR}
        echo -e "${GREEN}✓ Логи удалены${NC}"
    fi
else
    echo -e "${YELLOW}✓ Логи сохранены в ${LOG_DIR}${NC}"
fi

echo -e "\n${GREEN}=== Удаление завершено ===${NC}"

if [ "$REMOVE_LOGS" = false ]; then
    echo -e "${YELLOW}Примечание: Логи сохранены в ${LOG_DIR}${NC}"
    echo -e "Для полного удаления выполните: ${YELLOW}sudo bash scripts/uninstall.sh --full${NC}"
else
    echo -e "${GREEN}Выполнено полное удаление (включая логи)${NC}"
fi
