package util

import (
	"sync"
	"time"
)

type TaskExe func()

type WorkerPool struct {
	taskQueue chan TaskExe
	wg        sync.WaitGroup
}

func NewWorkerPool(size int) *WorkerPool {
	pool := &WorkerPool{
		taskQueue: make(chan TaskExe),
	}

	pool.wg.Add(size)
	for i := 0; i < size; i++ {
		go pool.worker()
	}

	return pool
}

func (p *WorkerPool) worker() {
	defer p.wg.Done()
	for task := range p.taskQueue {
		task()
	}
}

func (p *WorkerPool) Schedule(task TaskExe, delay time.Duration) {
	time.AfterFunc(delay, func() {
		p.taskQueue <- task
	})
}

func (p *WorkerPool) Stop() {
	close(p.taskQueue)
	p.wg.Wait()
}

//func main() {
//	pool := NewWorkerPool(5)
//
//	for i := 0; i < 10; i++ {
//		i := i
//		pool.Schedule(func() {
//			fmt.Printf("Executing task %d\n", i)
//		}, time.Duration(i)*time.Second)
//	}
//
//	time.Sleep(15 * time.Second)
//
//	pool.Stop()
//}
