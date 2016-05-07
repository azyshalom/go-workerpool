package workerpool

import (

)

type Func func(interface{})

type Worker struct {
    F Func
    I interface{}
}

type WorkerPool struct {
   stop  chan struct{}
   queue chan Worker
}

func NewWorkerPool(threadCount, queueCount int) *WorkerPool {
    ctx := &WorkerPool {
        stop:  make(chan struct{}, 0),
        queue: make(chan Worker, queueCount),
    }
    for i := 0; i < threadCount; i++ {
        go ctx.start(i)
    }
    return ctx
}

func (ctx *WorkerPool) Terminate() {
    close(ctx.stop)
}

func (ctx *WorkerPool) start(nn int) {
    for {
        select {
        case worker := <- ctx.queue:
            worker.F(worker.I)
        case <-ctx.stop:
            return
        }
    }
}

func (ctx *WorkerPool) Push(f Func, i interface{}) {
    select {
    case <-ctx.stop:
        return
    default:
        ctx.queue <- Worker{f, i}
    }
    
}