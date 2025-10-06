package chancount

import "sync/atomic"

// Elems keeps count of elements coming from the input channel and forwards them
// to the returned channel.
func Elems[T any](count *atomic.Uint64, input <-chan T) <-chan T {
	out := make(chan T, cap(input))
	go func() {
		for e := range input {
			out <- e
			count.Add(1)
		}
		close(out)
	}()
	return out
}
