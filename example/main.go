package main

import (
	"github.com/ullt/gopool"
	"log"
	"runtime"
	"time"
)

func main() {
	const (
		workQueueCap = 100000
		taskQueueCap = 100000
		taskNum      = 100000
		taskDuration = 3 * time.Second
	)
	n := runtime.GOMAXPROCS(runtime.NumCPU())
	log.Printf("NumCPU=%d", n)
	startTime := time.Now().UnixNano() / 1e6
	pool := gopool.New(workQueueCap, taskQueueCap)
	defer pool.Release()
	defer pool.WaitAll()
	pool.WaitCount(taskNum)
	for i := 0; i < taskNum; i++ {
		pool.Do(func() {
			time.Sleep(taskDuration)
			pool.TaskDone()
		})
	}
	pool.WaitAll()
	endTime := time.Now().UnixNano() / 1e6
	log.Printf("%d tasks took %d millisecond", taskNum, endTime-startTime)
}
