package main

import "fmt"

type MyLocker struct {
	locked bool
}

func (l *MyLocker) Lock() {
	l.locked = true
}

func (l *MyLocker) Unlock() {
	l.locked = false
}

type MyType struct {
	locker MyLocker
}

func (mt MyType) Test() {
	mt.locker.Lock()
	fmt.Println("test")
	mt.locker.Unlock()
}

// This example creates a PriorityQueue with some items, adds and manipulates an item,
// and then removes the items in priority order.
func main() {
	fmt.Println("Hello, Go!")
	defer fmt.Println("Bye, Go!")

	t := MyType{}
	t.Test()
}
