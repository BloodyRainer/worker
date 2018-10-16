package bottleneck

import (
	"fastworker/util"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

type Worker struct {
	id int
}

type task struct {
	handler http.Handler
	wg      *sync.WaitGroup
	rw      http.ResponseWriter
	req     *http.Request
}

var bouncer chan *Worker
var tasks chan *task
var logger *zap.Logger

func Init(n int) {

	go func() {
		logger = util.GetLogger()
		bouncer = make(chan *Worker)
		tasks = make(chan *task)

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
	logger.Info("finished work", zap.Int("worker-id", rcv.id))
	t.wg.Done()
}

func rescheduleWorker(w *Worker) {
	bouncer <- w
}
