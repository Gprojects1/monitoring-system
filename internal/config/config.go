package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// WebAppConfig конфигурация веб-приложения
type WebAppConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	LogFile  string `yaml:"log_file"`
	LogLevel string `yaml:"log_level"`
}

// MonitorConfig конфигурация системы мониторинга
type MonitorConfig struct {
	AppURL         string `yaml:"app_url"`
	HealthEndpoint string `yaml:"health_endpoint"`
	Timeout        int    `yaml:"timeout"`
	MaxRetries     int    `yaml:"max_retries"`
	RetryDelay     int    `yaml:"retry_delay"`
	ServiceName    string `yaml:"service_name"`
	LogFile        string `yaml:"log_file"`
	LogLevel       string `yaml:"log_level"`
	UseSystemd     bool   `yaml:"use_systemd"`
}

// LoadWebAppConfig загружает конфигурацию веб-приложения из YAML файла
func LoadWebAppConfig(path string) (*WebAppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла конфигурации: %w", err)
	}

	var cfg WebAppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("ошибка парсинга YAML: %w", err)
	}

	// Валидация конфигурации
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("ошибка валидации конфигурации: %w", err)
	}

	return &cfg, nil
}

// Validate проверяет корректность конфигурации веб-приложения
func (c *WebAppConfig) Validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("некорректный порт: %d", c.Port)
	}

	if c.Host == "" {
		return fmt.Errorf("хост не может быть пустым")
	}

	if c.LogFile == "" {
		return fmt.Errorf("путь к логу не может быть пустым")
	}

	return nil
}

// LoadMonitorConfig загружает конфигурацию мониторинга из YAML файла
func LoadMonitorConfig(path string) (*MonitorConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла конфигурации: %w", err)
	}

	var cfg MonitorConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("ошибка парсинга YAML: %w", err)
	}

	// Валидация конфигурации
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("ошибка валидации конфигурации: %w", err)
	}

	return &cfg, nil
}

// Validate проверяет корректность конфигурации мониторинга
func (c *MonitorConfig) Validate() error {
	if c.AppURL == "" {
		return fmt.Errorf("URL приложения не может быть пустым")
	}

	if c.Timeout < 1 {
		return fmt.Errorf("таймаут должен быть больше 0")
	}

	if c.MaxRetries < 1 {
		return fmt.Errorf("количество попыток должно быть больше 0")
	}

	if c.ServiceName == "" {
		return fmt.Errorf("имя службы не может быть пустым")
	}

	return nil
}
