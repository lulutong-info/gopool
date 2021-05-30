package gopool

import (
	"github.com/stretchr/testify/assert"
	"math"
	"sync"
	"testing"
)

func TestWorker_newWorker(t *testing.T) {
	as := assert.New(t)
	workerQueue := make(chan IWorker, workerQueueCap)
	worker := newWorker(workerQueue)
	defer worker.release()
	worker.start()
	var (
		num int
		wg  sync.WaitGroup
	)
	wg.Add(1)
	worker.do(func() {
		num = math.MaxInt32
		wg.Done()
	})
	wg.Wait()
	as.Equal(num, math.MaxInt32)
}
