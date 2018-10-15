package work

import (
	"bytes"
	"fastworker/util"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

type Worker struct {
	id int
}

type Task struct {
	request        *http.Request
	responseWriter http.ResponseWriter
	wait           *sync.WaitGroup
}

var bouncer chan *Worker
var tasks chan *Task
var once sync.Once
var logger *zap.Logger

func SubmitTask(w http.ResponseWriter, r *http.Request) *sync.WaitGroup{

	t := &Task{
		request:        r,
		responseWriter: w,
		wait:           &sync.WaitGroup{},
	}

	t.wait.Add(1)
	tasks <- t

	return t.wait
}

func (rcv *Worker) do(t *Task) error {
	buf := new(bytes.Buffer)

	d := 100 * time.Millisecond
	time.Sleep(d)
	s := fmt.Sprintf("worker-nr: %v finished some work of %v\n", rcv.id, d)
	_, err := buf.WriteString(s)
	if err != nil {
		return fmt.Errorf("error writing string to buffer: %v", err)
	}

	mw := io.MultiWriter(t.responseWriter, os.Stdout)
	io.WriteString(mw, buf.String())

	return nil
}

func InitWorkers(n int) {

	logger = util.GetLogger()
	bouncer = make(chan *Worker)
	tasks = make(chan *Task)

	go func() {
		for i := 0; i < n; i++ {
			bouncer <- &Worker{
				id: i,
			}
		}
	}()

	for {

		task := <-tasks

		select {
		case worker := <-bouncer:

			go func() {

				err := worker.do(task)
				if err != nil {
					logger.Error("error doing work", zap.Error(err))
				}
				task.wait.Done()

				go rescheduleWorker(worker)

			}()

		default:
			logger.Warn("no content delivered", zap.Int("status", 204))
			task.responseWriter.WriteHeader(http.StatusNoContent)
			task.wait.Done()
		}
	}
}

func rescheduleWorker(w *Worker) {
	bouncer <- w
}
