package logger

import (
	"log"
	"sync"

	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	once   sync.Once
)

func Init() {
	once.Do(func() {
		var err error
		logger, err = zap.NewProduction()
		if err != nil {
			log.Fatalf("Failed to initialize zap logger: %v", err)
		}
	})
}

func Get() *zap.Logger {
	if logger == nil {
		Init()
	}
	return logger
}
