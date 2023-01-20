package workerpool

import (
        "fmt"
        "testing"
        "time"
)

func TestWorkerPool(t *testing.T) {

        wp := New(10)
        for i := 0; i < 100; i++ {
                wp.Push(work, fmt.Sprintf("work %d", i))
        }

        time.Sleep(time.Second * 20)
        wp.Stop()
}

func work(i interface{}) {
        fmt.Printf("%s\n", i.(string))
        time.Sleep(time.Second * 1)
}
