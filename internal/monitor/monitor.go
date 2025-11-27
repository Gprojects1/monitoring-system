package monitor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"time"

	"github.com/Gprojects1/monitoring-system/internal/config"
	"github.com/Gprojects1/monitoring-system/internal/logger"
)

// Monitor представляет систему мониторинга приложения
type Monitor struct {
	config *config.MonitorConfig
	logger *logger.Logger
	client *http.Client
}

// HealthResponse структура ответа health check
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
}

// New создает новый экземпляр монитора
func New(cfg *config.MonitorConfig, log *logger.Logger) *Monitor {
	return &Monitor{
		config: cfg,
		logger: log,
		client: &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
		},
	}
}

// CheckHealth проверяет доступность приложения
func (m *Monitor) CheckHealth() bool {
	url := fmt.Sprintf("%s%s", m.config.AppURL, m.config.HealthEndpoint)

	m.logger.Debug("Проверка доступности", map[string]interface{}{
		"url": url,
	})

	resp, err := m.client.Get(url)
	if err != nil {
		m.logger.Error("✗ Ошибка при подключении к приложению", map[string]interface{}{
			"error": err.Error(),
		})
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		m.logger.Warning("✗ Приложение вернуло некорректный статус", map[string]interface{}{
			"status_code": resp.StatusCode,
		})
		return false
	}

	// Проверяем тело ответа
	var healthResp HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&healthResp); err != nil {
		m.logger.Warning("✗ Некорректный формат ответа health check", map[string]interface{}{
			"error": err.Error(),
		})
		return false
	}

	if healthResp.Status != "healthy" {
		m.logger.Warning("✗ Приложение сообщает о проблемах", map[string]interface{}{
			"status": healthResp.Status,
		})
		return false
	}

	m.logger.Info("✓ Приложение доступно и работает корректно", nil)
	return true
}

// RestartApplication перезапускает приложение через systemd
func (m *Monitor) RestartApplication() bool {
	m.logger.Warning("Попытка перезапуска службы", map[string]interface{}{
		"service": m.config.ServiceName,
	})

	if !m.config.UseSystemd {
		m.logger.Error("✗ Systemd отключен в конфигурации", nil)
		return false
	}

	// Выполняем команду перезапуска
	cmd := exec.Command("sudo", "systemctl", "restart", m.config.ServiceName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		m.logger.Error("✗ Ошибка при выполнении systemctl", map[string]interface{}{
			"error":  err.Error(),
			"output": string(output),
		})
		return false
	}

	m.logger.Info("✓ Служба успешно перезапущена", map[string]interface{}{
		"service": m.config.ServiceName,
	})
	return true
}

// WaitForStartup ожидает успешного запуска приложения
func (m *Monitor) WaitForStartup() bool {
	m.logger.Info("Ожидание запуска приложения...", map[string]interface{}{
		"max_retries": m.config.MaxRetries,
		"retry_delay": m.config.RetryDelay,
	})

	for attempt := 1; attempt <= m.config.MaxRetries; attempt++ {
		// Задержка перед проверкой
		time.Sleep(time.Duration(m.config.RetryDelay) * time.Second)

		m.logger.Info("Попытка проверки доступности", map[string]interface{}{
			"attempt": attempt,
			"total":   m.config.MaxRetries,
		})

		if m.CheckHealth() {
			m.logger.Info("✓ Приложение успешно запущено", map[string]interface{}{
				"attempt": attempt,
			})
			return true
		}

		if attempt < m.config.MaxRetries {
			m.logger.Info("Приложение еще не доступно, ожидание...", map[string]interface{}{
				"attempt": attempt,
				"total":   m.config.MaxRetries,
			})
		}
	}

	m.logger.Error("✗ Приложение не запустилось после перезапуска", map[string]interface{}{
		"attempts": m.config.MaxRetries,
	})
	return false
}

// Run выполняет основной цикл мониторинга
func (m *Monitor) Run() int {
	// Проверяем доступность
	if m.CheckHealth() {
		m.logger.Info("Мониторинг завершен: приложение работает нормально", nil)
		return 0
	}

	// Приложение недоступно - пытаемся перезапустить
	m.logger.Warning("Приложение недоступно. Требуется перезапуск.", nil)

	if !m.RestartApplication() {
		m.logger.Error("Не удалось перезапустить приложение", nil)
		return 1
	}

	// Ожидаем запуска и проверяем
	if m.WaitForStartup() {
		m.logger.Info("Приложение успешно восстановлено", nil)
		return 0
	}

	m.logger.Error("Не удалось восстановить приложение", nil)
	return 1
}
