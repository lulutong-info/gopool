package gopool

import "sync"

type IPool interface {
	Do(task Task)
	dispatch()
	getWorkerQueue() chan IWorker
	WaitCount(count int)
	TaskDone()
	WaitAll()
	Release()
	TaskQueueLen() int
}

type pool struct {
	workerQueue chan IWorker
	jobQueue    chan Task
	stop        chan signal
	wg          sync.WaitGroup
}

func (p *pool) dispatch() {
	var t Task
	for {
		select {
		case t = <-p.jobQueue:
			w := <-p.workerQueue
			w.do(t)
		case <-p.stop:
			for i := 0; i < cap(p.workerQueue); i++ {
				w := <-p.workerQueue
				w.release()
			}
			p.stop <- stopSignal
			return
		}
	}
}

func (p *pool) getWorkerQueue() chan IWorker {
	return p.workerQueue
}

func (p *pool) Do(task Task) {
	p.jobQueue <- task
}

func (p *pool) WaitCount(count int) {
	p.wg.Add(count)
}

func (p *pool) TaskDone() {
	p.wg.Done()
}

func (p *pool) WaitAll() {
	p.wg.Wait()
}

func (p *pool) Release() {
	p.stop <- stopSignal
	<-p.stop
}

func (p *pool) TaskQueueLen() int {
	return len(p.jobQueue)
}

func New(workQueueLen, taskQueueLen int) IPool {
	var p IPool = &pool{
		workerQueue: make(chan IWorker, workQueueLen),
		jobQueue:    make(chan Task, taskQueueLen),
		stop:        make(chan signal),
		wg:          sync.WaitGroup{},
	}
	for i := 0; i < cap(p.getWorkerQueue()); i++ {
		w := newWorker(p.getWorkerQueue())
		w.start()
	}
	go p.dispatch()
	return p
}
