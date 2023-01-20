package workerpool

import (
        "github.com/azyshalom/go-semaphore"
        "sync"
)

type Handler func(interface{})

type WorkerPool struct {
        stop chan struct{}
        pool chan interface{}
        sem  *semaphore.Semaphore
        wg   sync.WaitGroup
}

type workerItem struct {
        handler   Handler
        parameter interface{}
}

func New(maxSize int) *WorkerPool {
        wp := &WorkerPool{
                stop: make(chan struct{}),
                pool: make(chan interface{}, maxSize),
                sem:  semaphore.New(maxSize, maxSize),
        }
        go wp.start()
        return wp
}

func (wp *WorkerPool) Stop() {
        close(wp.stop)
        wp.wg.Wait()
}

func (wp *WorkerPool) Push(handler Handler, val interface{}) {
        wp.sem.Wait()
        wp.pool <- workerItem{handler: handler, parameter: val}
}

func (wp *WorkerPool) start() {
        for {
                select {
                case i := <-wp.pool:
                        go wp.work(i)
                case <-wp.stop:
                        return
                }
        }
}

func (wp *WorkerPool) work(i interface{}) {
        wp.wg.Add(1)
        defer wp.wg.Done()
        w, ok := i.(workerItem)
        if ok {
                w.handler(w.parameter)
        }
        wp.sem.Post()
}
