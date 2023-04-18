package main

import (
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"gostudy/pkg/rng"

	"github.com/stretchr/testify/assert"
)

// 우선순위 채널이 (1)우선순위 및 (2)입력순서대로 출력하는지 테스트한다.
func TestPriorityChannelShouldPopOrderedByPriorityAndInputOrder(t *testing.T) {
	assert := assert.New(t)
	repeats := 10000
	pc := NewChannel[int](repeats)

	// push to channel
	for i := 0; i < repeats; i++ {
		// input order, random priority
		pc.Push(i, rng.NextInRange(math.MinInt8, math.MaxInt8))
	}
	assert.Equal(repeats, pc.Count())

	// pop from channel
	items := []element[int]{}
	for i := 0; i < repeats; i++ {
		data, priority, err := pc.PopWithPriority()
		if err != nil {
			panic(err)
		}
		items = append(items, element[int]{data, priority})
	}
	assert.Zero(pc.Count())

	// assert items are ordered by priority and input order
	prev := element[int]{data: -1, priority: math.MinInt8 - 1}
	for _, item := range items {
		// 우선순위 기준으로 오름차순(minHeap)
		assert.GreaterOrEqual(item.priority, prev.priority)
		// 우선순위가 동일한 경우 입력 순서가 유지되는가?
		if item.priority == prev.priority {
			assert.Greater(item.data, prev.data)
		}
		prev = item
	}
	/* OUTPUT
	D:\gitworks\go-study\priority-queue>go test -v -run TestPriorityChannelShouldPopOrderedByPriorityAndInputOrder
	=== RUN   TestPriorityChannelShouldPopOrderedByPriorityAndInputOrder
	--- PASS: TestPriorityChannelShouldPopOrderedByPriorityAndInputOrder (0.02s)
	PASS
	ok      gostudy/priority-queue  0.051s
	*/
}

// 우선순위 채널이 concurrent-safe 한지 테스트한다
func TestPriorityChannelIsConcurrentSafe(t *testing.T) {

	// 우선순위 채널
	pc := NewChannel[time.Time](1000)

	wg := sync.WaitGroup{}
	workers := runtime.NumCPU() * 2 // concurrent-safe 검출을 쉽게하기 위함
	runtime.GOMAXPROCS(workers)

	// timeout trigger
	done := false
	go func(d time.Duration) {
		<-time.After(d) // wait for timeout
		done = true
	}(time.Second * 10)

	// pushed, poped count
	pushed, poped := uint64(0), uint64(0)

	// channel push goroutine
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(x int) {
			defer wg.Done()
			for !done {
				pc.Push(time.Now(), rng.Next[int]())
				atomic.AddUint64(&pushed, 1)
			}
		}(i)
	}

	// channel pop goroutine
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(x int) {
			defer wg.Done()
			for !done {
				pc.Pop()
				atomic.AddUint64(&poped, 1)
			}
		}(i)
	}

	// wait for all goroutine completed
	wg.Wait()

	// print push, pop action count
	t.Logf("push: %d, pop: %d", pushed, poped)

	/* OUTPUT
	D:\gitworks\go-study\priority-queue>go test -v -run TestPriorityChannelIsConcurrentSafe
	=== RUN   TestPriorityChannelIsConcurrentSafe
		priority_queue_test.go:101: push: 18060425, pop: 18060420
	--- PASS: TestPriorityChannelIsConcurrentSafe (10.00s)
	PASS
	ok      gostudy/priority-queue  10.038s
	*/
}
