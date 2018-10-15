package main

import (
	"fastworker/handlers"
	"fastworker/util"
	"fastworker/work"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

const port = 8080
const workers = 4

func main() {
	logger := util.GetLogger()
	defer logger.Sync()

	go work.InitWorkers(workers)

	http.Handle("/", handlers.RootHandler{})

	logger.Info("starting server", zap.Int("port", port))
	err := http.ListenAndServe(":" + strconv.Itoa(8080), nil)
	if err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
	}
}
