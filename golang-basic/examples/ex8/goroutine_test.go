package main

import (
	"testing"
	"time"
)

func TestGoRoutineConcurrentNotSafe(t *testing.T) {
	noop := func() bool { time.Sleep(0); return true }

	// timeout for 10 seconds
	done := false
	go func() {
		<-time.After(time.Second * 10)
		done = true
	}()

	for !done {
		data := 0
		go func() { data++ }()
		if data == 0 && noop() && data != 0 {
			panic("wtf! data is zero or not?")
		}
	}
}
