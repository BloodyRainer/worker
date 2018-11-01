package main

import (
	"fastworker/bottleneck"
	"fastworker/handlers"
	"fastworker/util"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

const (
	port = 8080
	workers = 400
)

func main() {
	logger := util.GetLogger()
	defer logger.Sync()

	bottleneck.Init(workers)

	http.Handle("/", bottleneck.Apply(handlers.RootHandler{}))

	fmt.Printf("starting server on port %v\n", port)
	err := http.ListenAndServe(":" + strconv.Itoa(8080), nil)
	if err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
	}
}
