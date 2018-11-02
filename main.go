package main

import (
	"fastworker/bottleneck"
	"fastworker/handlers"
	"fastworker/util"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	port                = 8080
	workers             = 100
	secGracefulShutdown = 25
)

var (
	logger *zap.Logger
	done chan bool
)

func main() {
	logger = util.GetLogger()
	defer logger.Sync()

	sigs := make(chan os.Signal)
	done := make(chan bool)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	bottleneck.LogBusyWorkers(5 * time.Second)

	go func() {
		<-sigs
		bottleneck.NotifyStop() <- true
		done <- true
	}()

	go startWebserver()

	<- done
	logger.Info("waiting for graceful shutdown", zap.Int("seconds", secGracefulShutdown))
	time.Sleep(time.Duration(secGracefulShutdown) * time.Second)

}

func startWebserver() {
	http.Handle("/", bottleneck.ApplyNumWorkers(workers, handlers.RootHandler{}))

	fmt.Printf("starting server on port %v\n", port)
	err := http.ListenAndServe(":" + strconv.Itoa(8080), nil)
	if err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
	}
}
