package util

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
	"time"
)

var logger *zap.Logger

var once sync.Once

func GetLogger() *zap.Logger {
	once.Do(func() {
		var err error

		cfg := zap.Config{
			Encoding: "json",
			Level: zap.NewAtomicLevelAt(zapcore.DebugLevel),
			OutputPaths: []string{"stdout"},
			EncoderConfig: zapcore.EncoderConfig{
				MessageKey: "message",
				TimeKey: "time",
				EncodeTime: TimeEncoder,
				LevelKey: "level",
				EncodeLevel: zapcore.CapitalLevelEncoder,
			},
		}

		logger, err = cfg.Build()
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to initialize zap logger")
			os.Exit(1)
		}
	})
	return logger
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02T15:04:05.00"))
}
