package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	Sugar *zap.SugaredLogger
}

func SetupLogger() (*Logger, error) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.DisableStacktrace = true

	logger, err := config.Build()

	if err != nil {
		return nil, err
	}

	return &Logger{logger.Sugar()}, nil
}

func (l *Logger) Sync() {
	l.Sugar.Sync()
}
