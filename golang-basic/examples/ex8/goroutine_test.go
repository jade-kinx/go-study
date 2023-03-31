package main

import (
	"testing"
	"time"
)

func TestGoRoutineConcurrentNotSafe(t *testing.T) {

	// timeout 10 seconds
	done := false
	go func() {
		<-time.After(time.Second * 10)
		done = true
	}()

	for !done {
		zero := 0
		go func() { zero++ }()
		if zero == 0 {
			time.Sleep(time.Millisecond * 10) // sleep 10ms
			if zero != 0 {
				panic("zero is not 0")
			}
		}
	}
}
