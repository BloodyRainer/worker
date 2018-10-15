package util

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"sync"
)

var logger *zap.Logger

var once sync.Once

func GetLogger() *zap.Logger {
	once.Do(func() {
		var err error
		logger, err = zap.NewProduction()
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to initialize zap logger")
		}
	})
	return logger
}
