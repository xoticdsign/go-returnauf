package logging

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func RunZap() error {
	var err error

	config := zap.Config{
		Level:    zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "Сообщение",
			LevelKey:       "Уровень",
			TimeKey:        "Время",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02, 15:04:05.000"),
			EncodeDuration: zapcore.StringDurationEncoder,
		},
		OutputPaths: []string{func() string {
			logs, err := os.Open("auf-citaty.log")
			if err != nil {
				os.Create("auf-citaty.log")
			}
			defer logs.Close()

			path, _ := filepath.Abs("auf-citaty.log")

			return path
		}()},
	}

	Logger, err = config.Build()
	if err != nil {
		return err
	}
	return nil
}
