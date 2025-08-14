package logger

import (
	loggerContract "github.com/barantoraman/microgate/pkg/logger/contract"
	"github.com/barantoraman/microgate/pkg/logger/zaplogger"
)

func GetLogger(level string) loggerContract.Logger {
	return zaplogger.NewLoggerWithLevel(level)
}
