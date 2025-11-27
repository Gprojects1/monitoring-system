package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// Logger обертка над logrus для удобного логирования
type Logger struct {
	log  *logrus.Logger
	file *os.File
}

// New создает новый экземпляр логгера
func New(logFile string, level string) (*Logger, error) {
	// Создаем директорию для логов
	logDir := filepath.Dir(logFile)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("ошибка создания директории для логов: %w", err)
	}

	// Открываем файл для логов
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия файла логов: %w", err)
	}

	// Создаем logrus логгер
	log := logrus.New()
	log.SetOutput(file)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Устанавливаем уровень логирования
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	log.SetLevel(logLevel)

	return &Logger{
		log:  log,
		file: file,
	}, nil
}

// Info логирует информационное сообщение
func (l *Logger) Info(msg string, fields map[string]interface{}) {
	if fields != nil {
		l.log.WithFields(fields).Info(msg)
	} else {
		l.log.Info(msg)
	}
}

// Debug логирует отладочное сообщение
func (l *Logger) Debug(msg string, fields map[string]interface{}) {
	if fields != nil {
		l.log.WithFields(fields).Debug(msg)
	} else {
		l.log.Debug(msg)
	}
}

// Warning логирует предупреждение
func (l *Logger) Warning(msg string, fields map[string]interface{}) {
	if fields != nil {
		l.log.WithFields(fields).Warn(msg)
	} else {
		l.log.Warn(msg)
	}
}

// Error логирует ошибку
func (l *Logger) Error(msg string, fields map[string]interface{}) {
	if fields != nil {
		l.log.WithFields(fields).Error(msg)
	} else {
		l.log.Error(msg)
	}
}

// Close закрывает файл логов
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}
