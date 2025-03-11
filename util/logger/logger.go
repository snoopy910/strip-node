package logger

import (
	"log"
	"sync"

	"go.uber.org/zap"
)

var (
	sugar     *zap.SugaredLogger
	once      sync.Once
	initError error
)

// Init initializes the logger. It only executes once.
func Init() error {
	once.Do(func() {
		logger, initError := zap.NewProduction()
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
