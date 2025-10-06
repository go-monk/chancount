package main

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/go-monk/chancount"
)

const n = 10

func main() {
	ch := generateInts(n)

	var count atomic.Uint64
	ch = chancount.Elems(&count, ch)

	// Print the count of received elements every second.
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			fmt.Printf("chan elements: %d\n", count.Load())
		}
	}()

	processInts(ch)

}

func processInts(ch <-chan int) {
	for range ch {
	}
}

func generateInts(n int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := range n {
			time.Sleep(time.Second)
			ch <- i
		}
	}()
	return ch
}
