// Package chancount provides utilities for counting and reporting elements
// passing through Go channels.
//
// Example usage:
//
//	ch := generateInts(n)
//	ch = chancount.Elems(ch)
//	processInts(ch)
package chancount

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type chanCounter[T any] struct {
	name string
	freq time.Duration
	ch   <-chan T
	w    io.Writer
	sync.RWMutex
	count      int
	start      time.Time
	timeFormat string
}

type option[T any] func(*chanCounter[T])

// WithName sets the name for the channel. The name is visible in reports. It's
// useful when you count elements of more than one channel.
func WithName[T any](name string) option[T] {
	return func(c *chanCounter[T]) {
		c.name = name
	}
}

// WithFrequency sets the reporting frequency.
func WithFrequency[T any](freq time.Duration) option[T] {
	return func(c *chanCounter[T]) {
		c.freq = freq
	}
}

// WithOutput sets where reports should go.
func WithOutput[T any](w io.Writer) option[T] {
	return func(c *chanCounter[T]) {
		c.w = w
	}
}

// WithTimeFormat sets the time format layout used in reports.
func WithTimeFormat[T any](layout string) option[T] {
	return func(c *chanCounter[T]) {
		c.timeFormat = layout
	}
}

// Elems keeps counting and reporting the number of elements received on ch
// until ch is closed. The default report frequency is [time.Second], the
// default report time layout is [time.TimeOnly] and the default report output
// is [os.Stdout].
func Elems[T any](ch <-chan T, opts ...option[T]) <-chan T {
	c := newChanCounter(ch, opts...)
	return c.Count()
}

func newChanCounter[T any](ch <-chan T, opts ...option[T]) *chanCounter[T] {
	c := &chanCounter[T]{
		ch:         ch,
		freq:       time.Second,
		w:          os.Stdout,
		start:      time.Now(),
		timeFormat: time.TimeOnly,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *chanCounter[T]) Count() <-chan T {
	out := make(chan T, cap(c.ch))
	go func() {
		defer close(out)

		t := time.NewTicker(c.freq)
		defer t.Stop()

		start := time.Now()

		for {
			select {
			case <-t.C:
				c.report(start)
			case e, ok := <-c.ch:
				if !ok { // channel closed
					return
				}
				out <- e
				c.Lock()
				c.count++
				c.Unlock()
			}
		}
	}()
	return out
}

func (c *chanCounter[T]) report(start time.Time) {
	c.RLock()
	count := c.count
	c.RUnlock()

	var msg string
	if c.name != "" {
		msg += fmt.Sprintf("%s ", c.name)
	}

	elements := "elements"
	if count == 1 {
		elements = "element"
	}

	msg += fmt.Sprintf("channel received %d %s since %s",
		count,
		elements,
		start.Format(c.timeFormat))

	fmt.Fprintln(c.w, msg)
}
