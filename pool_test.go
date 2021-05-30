package gopool

import (
	"github.com/stretchr/testify/assert"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

const (
	workerQueueCap = 100
	taskQueueCap   = 100
)

func TestMain(m *testing.M) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	m.Run()
}

func TestNew(t *testing.T) {
	as := assert.New(t)
	pool := New(workerQueueCap, taskQueueCap)
	defer pool.Release()
	as.NotNil(pool)
}

func TestPool_Release(t *testing.T) {
	as := assert.New(t)
	pool := New(workerQueueCap, taskQueueCap)
	as.NotNil(pool)
	count := workerQueueCap
	pool.WaitCount(count)
	for i := 0; i < count; i++ {
		pool.Do(func() {
			demoTask()
			pool.TaskDone()
		})
	}
	pool.WaitAll()
	time.Sleep(100 * time.Millisecond)
	pool.Release()
}

func TestPool_TaskQueueLen(t *testing.T) {
	as := assert.New(t)
	pool := New(workerQueueCap, taskQueueCap)
	defer pool.Release()
	count := 1
	for i := 0; i < count; i++ {
		pool.Do(func() {
			demoTask()
		})
	}
	as.Equal(pool.TaskQueueLen(), count)
}

func TestPanic(t *testing.T) {
	pool := New(workerQueueCap, taskQueueCap)
	pool.WaitCount(taskQueueCap)
	defer pool.Release()
	for i := 0; i < taskQueueCap; i++ {
		pool.Do(func() {
			defer pool.TaskDone()
			panic("create a panic")
		})
	}
	pool.WaitAll()
}

func BenchmarkPool(b *testing.B) {
	b.ReportAllocs()
	pool := New(workerQueueCap, taskQueueCap)
	defer pool.Release()
	defer pool.WaitAll()
	for n := 0; n < b.N; n++ {
		for i := 0; i < 1000000; i++ {
			pool.WaitCount(1)
			pool.Do(func() {
				demoTask()
				pool.TaskDone()
			})
		}
	}
}

var num int32 = 1

func demoTask() {
	atomic.AddInt32(&num, 1)
}
