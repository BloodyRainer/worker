package bottleneck

import (
	"fastworker/util"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

type Worker struct {
	id int64
}

type task struct {
	handler http.Handler
	wg      *sync.WaitGroup
	rw      http.ResponseWriter
	req     *http.Request
}

var (
	bouncer chan *Worker
	tasks chan *task
	logger *zap.Logger
)

// Initialize the Bottleneck
func Init(n int) {

	logger = util.GetLogger()
	bouncer = make(chan *Worker, n)
	tasks = make(chan *task)

	go supplyWorkers(n, bouncer)

	go func() {
		for {

			task := <-tasks

			select {
			case worker := <-bouncer:

				go func() {
					worker.do(task)
					go rescheduleWorker(worker)
				}()

			default:
				logger.Warn("no content delivered", zap.Int("status", 204))
				task.rw.WriteHeader(http.StatusNoContent)
				task.wg.Done()
			}
		}
	}()

}

// Apply the Bottleneck to the Handler
func Apply(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		waitGrp := new(sync.WaitGroup)
		waitGrp.Add(1)

		t := &task{
			handler: next,
			wg:      waitGrp,
			rw:      w,
			req:     r,
		}

		tasks <- t
		waitGrp.Wait()
	})
}

func (rcv *Worker) do(t *task) {
	t.handler.ServeHTTP(t.rw, t.req)
	logger.Info("finished work", zap.Int64("worker-id", rcv.id))
	t.wg.Done()
}

func supplyWorkers(n int, bouncer chan<- *Worker) {
	for i := 0; i < n; i++ {
		bouncer <- &Worker{
			id: int64(i),
		}
	}
}

func rescheduleWorker(w *Worker) {
	bouncer <- w
}
