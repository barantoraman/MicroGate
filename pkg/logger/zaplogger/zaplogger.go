package zaplogger

import (
	"os"
	"strings"

	loggerContract "github.com/barantoraman/microgate/pkg/logger/contract"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	logger *zap.Logger
}

// Debug implements contract.Logger.
func (l *zapLogger) Debug(msg string, fields ...zap.Field) {
	panic("unimplemented")
}

// Error implements contract.Logger.
func (l *zapLogger) Error(msg string, fields ...zap.Field) {
	panic("unimplemented")
}

// Fatal implements contract.Logger.
func (l *zapLogger) Fatal(msg string, fields ...zap.Field) {
	panic("unimplemented")
}

// Info implements contract.Logger.
func (l *zapLogger) Info(msg string, fields ...zap.Field) {
	panic("unimplemented")
}

// Warn implements contract.Logger.
func (l *zapLogger) Warn(msg string, fields ...zap.Field) {
	panic("unimplemented")
}

func NewLoggerWithLevel(level string) loggerContract.Logger {
	logLevelMap := map[string]zapcore.Level{
		"debug": zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
		"fatal": zapcore.FatalLevel,
	}

	parsedLevel, ok := logLevelMap[strings.ToLower(level)]
	if !ok {
		parsedLevel = zapcore.InfoLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "timestamp"
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		parsedLevel,
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return &zapLogger{
		logger: logger,
	}
}
