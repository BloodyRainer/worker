package work

import (
	"fastworker/util"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

type Worker struct {
	id int
}

type Task struct {
	request        *http.Request
	responseWriter http.ResponseWriter
	wait           *sync.WaitGroup
	do             func(http.ResponseWriter, *http.Request) error
}

var bouncer chan *Worker
var tasks chan *Task
var once sync.Once
var logger *zap.Logger

func SubmitTask(w http.ResponseWriter, r *http.Request, do func(http.ResponseWriter, *http.Request) error) *sync.WaitGroup {

	t := &Task{
		request:        r,
		responseWriter: w,
		wait:           &sync.WaitGroup{},
		do:             do,
	}

	t.wait.Add(1)
	tasks <- t

	return t.wait
}

func (rcv *Worker) work(t *Task) error {
	err := t.do(t.responseWriter, t.request)
	if err != nil {
		return err
	}
	logger.Info("finished work", zap.Int("worker-id", rcv.id))
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
				err := worker.work(task)
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
