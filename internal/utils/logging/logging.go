package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func RunZap() error {
	var err error

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

	Logger, err = config.Build()
	if err != nil {
		return err
	}
	return nil
}
