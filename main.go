package main

import (
	"fastworker/handlers"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strconv"
)

const port = 8080

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to initialize zap logger")
	}
	defer logger.Sync()

	http.Handle("/", handlers.RootHandler{})

	logger.Info("starting server", zap.Int("port", port))
	err = http.ListenAndServe(":" + strconv.Itoa(8080), nil)
	if err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
	}


}
