package main

import (
	"time"

	"github.com/go-monk/chancount"
)

const n = 10

func main() {
	ch := generateInts(n)
	ch = chancount.Elems(ch)
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
