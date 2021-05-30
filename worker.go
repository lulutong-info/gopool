package gopool

import (
	"fmt"
	"log"
)

type IWorker interface {
	start()
	do(Task)
	release()
}

type worker struct {
	workerQueue chan IWorker
	taskQueue   chan Task
	stop        chan signal
}

func (w *worker) start() {
	go func() {
		defer func() {
			w.start()
			if e := recover(); e != nil {
				log.Println(fmt.Sprintf("gopool: error:%v", e))
			}
		}()
		defer func() {
			w.workerQueue <- w
		}()
		var t Task
		for {
			select {
			case t = <-w.taskQueue:
				t()
			case <-w.stop:
				w.stop <- stopSignal
				return
			}
		}
	}()
}

func (w *worker) do(task Task) {
	w.taskQueue <- task
}

// worker 出队，释放 worker资源
func (w *worker) release() {
	w.stop <- stopSignal
	<-w.stop
}

func newWorker(workerQueue chan IWorker) IWorker {
	w := &worker{
		workerQueue: workerQueue,
		taskQueue:   make(chan Task),
		stop:        make(chan signal),
	}
	w.workerQueue <- w
	return w
}
