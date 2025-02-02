package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type AppLogger interface {
	Debug(msg string, fields ...any)
	Error(msg string, fields ...any)
	Fatal(msg string, fields ...any)
	Info(msg string, fields ...any)
}

type Logger struct {
	logger *zap.Logger
}

var _ AppLogger = (*Logger)(nil)

func Initialize(level string, outPath ...string) (*Logger, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	if outPath != nil {
		cfg.OutputPaths = append(outPath, "stderr")
	}

	zl, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{zl}, nil
}

func (l *Logger) Debug(msg string, fields ...any) {
	l.logger.Debug(msg, toZapFields(fields...)...)
}
func (l *Logger) Error(msg string, fields ...any) {
	l.logger.Error(msg, toZapFields(fields...)...)
}
func (l *Logger) Fatal(msg string, fields ...any) {
	l.logger.Fatal(msg, toZapFields(fields...)...)
}
func (l *Logger) Info(msg string, fields ...any) {
	l.logger.Info(msg, toZapFields(fields...)...)
}

func toZapField(field any) zapcore.Field {
	if v, ok := field.(zapcore.Field); ok {
		return v
	}
	return zapcore.Field{}
}

func toZapFields(fields ...any) []zapcore.Field {
	if len(fields) > 0 {
		zapFields := make([]zapcore.Field, len(fields))

		for i := range fields {
			zapFields[i] = toZapField(fields[i])
		}

		return zapFields
	}
	return []zapcore.Field{}
}
