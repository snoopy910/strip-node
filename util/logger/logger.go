package logger

import (
	"log"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	sugar     *zap.SugaredLogger
	once      sync.Once
	initError error
)

// Init initializes the logger. It only executes once.
func Init() error {
	once.Do(func() {
		config := zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		logger, initError := config.Build(
			zap.AddStacktrace(zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= zap.ErrorLevel
			})),
		)
		if initError == nil {
			sugar = logger.Sugar()
		}
	})
	return initError
}

// Sugar returns the SugaredLogger instance
func Sugar() *zap.SugaredLogger {
	if sugar == nil {
		if err := Init(); err != nil {
			// Fallback to a no-op logger rather than return nil
			logger, err := zap.NewProduction()
			if err != nil {
				log.Fatal(err)
			}
			sugar = logger.Sugar()
		}
	}
	return sugar
}

// Sync flushes any buffered log entries
func Sync() error {
	if sugar != nil {
		return sugar.Sync()
	}
	return nil
}
