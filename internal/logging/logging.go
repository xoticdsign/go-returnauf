package logging

import (
	"github.com/gofiber/fiber/v2"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Info(message string, c *fiber.Ctx)
	Warn(message string, c *fiber.Ctx)
	Error(message string, c *fiber.Ctx)
}

type Service struct {
	logger *zap.Logger
}

func RunZap() (*Service, error) {
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
	return &Service{logger: zap}, nil
}

func (s *Service) Info(message string, c *fiber.Ctx) {
	s.logger.Info(
		message,
		zap.String("Method", c.Method()),
		zap.String("Path", c.Path()),
	)
}

func (s *Service) Warn(message string, c *fiber.Ctx) {
	s.logger.Warn(
		message,
		zap.String("Method", c.Method()),
		zap.String("Path", c.Path()),
	)
}

func (s *Service) Error(message string, c *fiber.Ctx) {
	s.logger.Error(
		message,
		zap.String("Method", c.Method()),
		zap.String("Path", c.Path()),
	)
}
