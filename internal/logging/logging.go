package logging

import (
	"github.com/gofiber/fiber/v2"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Интерфейс, содержащий методы для работы с Логгером
type Logger interface {
	Info(message string, c *fiber.Ctx)
	Warn(message string, c *fiber.Ctx)
	Error(message string, c *fiber.Ctx)
}

// Структура, реализующая Logger
type Log struct {
	logger *zap.Logger
}

// Запускает Zap и возвращает структуру, реализующую Logger
func RunZap() (*Log, error) {
	config := zap.Config{
		Level:    zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding: "console",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:       "MESSAGE",
			LevelKey:         "LEVEL",
			TimeKey:          "TIME",
			EncodeLevel:      zapcore.CapitalLevelEncoder,
			EncodeTime:       zapcore.TimeEncoderOfLayout("2006-01-02, 15:04:05.000"),
			EncodeDuration:   zapcore.StringDurationEncoder,
			ConsoleSeparator: " | ",
		},
		OutputPaths: []string{"stdout"},
	}

	zap, err := config.Build()
	if err != nil {
		return nil, err
	}
	return &Log{logger: zap}, nil
}

// Создает информационный лог
func (l *Log) Info(message string, c *fiber.Ctx) {
	l.logger.Info(
		message,
		zap.String("Method", c.Method()),
		zap.String("Path", c.Path()),
	)
}

// Создает лог-предупреждение
func (l *Log) Warn(message string, c *fiber.Ctx) {
	l.logger.Warn(
		message,
		zap.String("Method", c.Method()),
		zap.String("Path", c.Path()),
	)
}

// Создает лог-ошибку
func (l *Log) Error(message string, c *fiber.Ctx) {
	l.logger.Error(
		message,
		zap.String("Method", c.Method()),
		zap.String("Path", c.Path()),
	)
}
