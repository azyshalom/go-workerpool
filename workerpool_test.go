package workerpool

import (
    "testing"
    "fmt"
    "time"
    "math/rand"
    "runtime"
)

const (
    THREAD_COUNT = 10
)

func TestWorkerPool(t *testing.T) {
   runtime.GOMAXPROCS(runtime.NumCPU())
   wp := NewWorkerPool(THREAD_COUNT)
   for i := 0; i < 500; i++ {
        go wp.Push(test, i)
   }
   time.Sleep(30 * time.Second)
   wp.Terminate()
}

func test(i interface{}) {
    n := i.(int)
    fmt.Printf("*** %03d\n", n)
    time.Sleep(time.Duration(rand.Int31n(1000)) * time.Millisecond)
}