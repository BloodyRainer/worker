package bottleneck

import (
	"fastworker/util"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type worker struct {
	id int64
}

type task struct {
	handler http.Handler
	wg      *sync.WaitGroup
	rw      http.ResponseWriter
	req     *http.Request
}

var (
	bouncer     chan *worker
	tasks       chan *task
	logger      *zap.Logger
	busyWorkers int64
	stopNotify  chan bool
	stop        bool
)

const (
	defaultWorkers = 1000
	gracePeriodMs  = 5000
)

func initBottleneck(n int) {

	logger = util.GetLogger()
	logger.Debug("initialize bottleneck", zap.Int("workers", n))

	bouncer = make(chan *worker, n)
	tasks = make(chan *task)
	stopNotify = make(chan bool)

	go supplyWorkers(n, bouncer)

	go func() {

	loop:
		for {

			task := <-tasks

			select {
			case w, ok := <-bouncer:
				if !ok {
					logger.Debug("bouncer channel has been closed")
					bouncer = nil
					break loop
				}

				go func(wo *worker) {
					atomic.AddInt64(&busyWorkers, 1)
					wo.do(task)

					go rescheduleWorker(wo)
				}(w)

			case <-time.After(1 * time.Second):
				if stop {
					logger.Debug("closing bouncer after 5 seconds")
					close(bouncer)
				}

			default:
				logger.Warn("no content delivered", zap.Int("status", 204))
				task.rw.WriteHeader(http.StatusNoContent)
				task.wg.Done()
			}
		}
	logger.Debug("don't take tasks no more")
	}()

}

func rescheduleWorker(w *worker) {
	// only rescheduleWorker if bottleneck should continue to do work
	if !stop {
		bouncer <- w
	}

	atomic.AddInt64(&busyWorkers, -1)
}

func (rcv *worker) do(t *task) {
	t.handler.ServeHTTP(t.rw, t.req)
	logger.Debug("finished work", zap.Int64("worker-id", rcv.id))
	t.wg.Done()
}

func supplyWorkers(n int, bouncer chan<- *worker) {
	for i := 0; i < n; i++ {
		bouncer <- &worker{
			id: int64(i),
		}
	}

	<- stopNotify
	logger.Debug("bottleneck received notification to stop work")
	stop = true
	for {
		var tries int
		//logger.Debug("waiting for bottleneck workers to finish their work...")
		if busyWorkers == 0 || tries > 5{
			logger.Debug("all workers have finished their work, closing bouncer", zap.Int64("busy workers", busyWorkers))
			close(bouncer)
			break
		}
		tries++
		time.Sleep(time.Duration(gracePeriodMs/ 5) * time.Millisecond)
	}
	logger.Debug("bottleneck stopped doing work")

}

// Apply the Bottleneck to the Handler.
func Apply(httpHandler http.Handler) http.Handler {
	return ApplyNumWorkers(defaultWorkers, httpHandler)
}

// Apply the Bottleneck to the Handler.
// numWorkers defines the number of workers the load is distributed to.
func ApplyNumWorkers(numWorkers int, httpHandler http.Handler) http.Handler {

	initBottleneck(numWorkers)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		waitTask := new(sync.WaitGroup)
		waitTask.Add(1)

		tasks <- &task{
			handler: httpHandler,
			wg:      waitTask,
			rw:      w,
			req:     r,
		}

		waitTask.Wait()
	})
}

// Logs the number of busy workers with the given time frequency.
func LogBusyWorkers(frequency time.Duration) {

	go func() {
		for range time.Tick(frequency) {
			logger.Info("busy workers", zap.Int64("number", busyWorkers))
		}
	}()

}

// Notify the Bottleneck to stop doing work gracefully.
func NotifyStop() chan<- bool {
	return stopNotify
}
