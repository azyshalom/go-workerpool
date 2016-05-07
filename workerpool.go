package workerpool

import (
    "sync"
)

type Func func(interface{})
type Semaphore chan struct{}

type Worker struct {
    F Func
    I interface{}
}

type WorkerPool struct {
   stop    chan struct{}
   queue   chan Worker
   empty   Semaphore
   wg      sync.WaitGroup
}

func NewWorkerPool(maxCount int) *WorkerPool {
    ctx := &WorkerPool {
        stop:    make(chan struct{}, 0),
        queue:   make(chan Worker, maxCount),
        empty:   make(Semaphore, maxCount),
    }
    for i := 0; i < maxCount; i++ {
        ctx.empty <- struct{}{}
    }
    go ctx.start()
    return ctx
}

func (ctx *WorkerPool) Terminate() {
    close(ctx.stop)
    ctx.wg.Wait()
}

func (ctx *WorkerPool) Push(f Func, i interface{}) {
    select {
    case <-ctx.stop:
        return
    case <-ctx.empty:
        ctx.queue <- Worker{f, i}
    }
}

func (ctx *WorkerPool) start() {
    for {
        select {
        case worker := <- ctx.queue:
            go ctx.work(worker)
        case <-ctx.stop:
            return
        }
    }
}

func (ctx *WorkerPool) work(worker Worker) {
    ctx.wg.Add(1)
    defer ctx.wg.Done()
    worker.F(worker.I)
    ctx.empty <- struct{}{}
}

