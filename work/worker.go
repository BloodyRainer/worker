package work

import (
	"bytes"
	"fmt"
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
}

var bouncer chan *Worker
var tasks chan *Task
var once sync.Once

func AddTask(w http.ResponseWriter, r *http.Request) {
	fmt.Println("add task")

	t := &Task{
		request:        r,
		responseWriter: w,
	}

	tasks <- t
}

func (rcv *Worker) Do(w http.ResponseWriter, r *http.Request) error {
	buf := new(bytes.Buffer)

	time.Sleep(1000 * time.Millisecond)
	s := fmt.Sprintf("worker-nr: %v finished some work \n", rcv.id)
	_, err := buf.WriteString(s)
	if err != nil {
		return fmt.Errorf("error writing string to buffer: %v", err)
	}

	mw := io.MultiWriter(os.Stdout, w)

	_, err = io.Copy(mw, buf)
	if err != nil {
		return fmt.Errorf("error writing to multiwriter: %v", err)
	}

	return nil
}

func InitWorkers(n int) {

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

		fmt.Println("before receiving tasks")
		t := <-tasks
		fmt.Println("after receiving tasks")

		select {
		case worker := <-bouncer:

			fmt.Println("got worker")
			go func() {

				err := worker.Do(t.responseWriter, t.request)
				if err != nil {
					// TODO: error to responsewriter
					fmt.Fprintf(os.Stderr, "error doing work: %v\n", err)
				}

				go rescheduleWorker(worker)

			}()

		default:
			fmt.Fprintf(t.responseWriter, "no content\n")
			fmt.Fprintf(os.Stdout, "no content\n")

		}
	}
}

func rescheduleWorker(w *Worker) {
	bouncer <- w
}
