package main

import (
	"container/heap"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// waitForPriorityQueue는 동시성을 보장하지 않는다.
func TestWaitForPriorityQueueIsConcurrentUnsafe(t *testing.T) {

	// 우선순위 큐 생성 및 초기화
	pq := &waitForPriorityQueue{}
	heap.Init(pq)

	wg := sync.WaitGroup{}
	workers := runtime.NumCPU() * 2
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
				heap.Push(pq, &waitFor{data: x, readyAt: time.Now()})
				atomic.AddUint64(&pushed, 1) // concurrent safe
			}
		}(i)
	}

	// channel pop goroutine
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(x int) {
			defer wg.Done()
			for !done {
				if pq.Len() > 0 {
					heap.Pop(pq)
				}
				atomic.AddUint64(&poped, 1) // concurrent safe
			}
		}(i)
	}

	// wait for all goroutine completed
	wg.Wait()

	// print push, pop action count
	t.Logf("push: %d, pop: %d", pushed, poped)

	/* OUTPUT
	D:\gitworks\go-study\priority-queue>go test -v -run ^TestWaitForPriorityQueueIsConcurrentUnsafe$
	=== RUN   TestWaitForPriorityQueueIsConcurrentUnsafe
	panic: runtime error: invalid memory address or nil pointer dereference
	[signal 0xc0000005 code=0x1 addr=0x28 pc=0xe29263]

	goroutine 88 [running]:
	gostudy/priority-queue.(*waitForPriorityQueue).Pop(0xe8efa0?)
			D:/gitworks/go-study/priority-queue/waitForPriorityQueue.go:70 +0x23
	container/heap.Pop({0xe8efa0, 0xc000008108})
			C:/Program Files/Go/src/container/heap/heap.go:63 +0x6b
	gostudy/priority-queue.TestWaitForPriorityQueueIsConcurrentUnsafe.func3(0x0?)
			D:/gitworks/go-study/priority-queue/waitForPriorityQueue_test.go:50 +0x65
	created by gostudy/priority-queue.TestWaitForPriorityQueueIsConcurrentUnsafe
			D:/gitworks/go-study/priority-queue/waitForPriorityQueue_test.go:47 +0x277
	exit status 2
	FAIL    gostudy/priority-queue  0.043s
	*/
}

// Push()와 Pop() 사이에 크리티컬 섹션을 설정하여 동시성을 보장
func TestWaitForPriorityQueueIsConcurrentSafe(t *testing.T) {

	// lock
	lock := sync.Mutex{}

	// 우선순위 큐
	pq := &waitForPriorityQueue{}
	heap.Init(pq)

	wg := sync.WaitGroup{}
	workers := runtime.NumCPU() * 2
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
			for !done {
				func() {
					////////////////////
					// CRITICAL SECTION
					lock.Lock()
					defer lock.Unlock()
					heap.Push(pq, &waitFor{data: x, readyAt: time.Now()})
					pushed++
					////////////////////
					// waitForPriorityQueue.Push() 내부에서 크리티컬 섹션을 구현하는 편이 좋음
					// 참고: https://github.com/caffix/queue/blob/master/queue.go
				}()
			}
			wg.Done()
		}(i)
	}

	// channel pop goroutine
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(x int) {
			for !done {
				func() {
					///////////////////////
					// CRITICAL SECTION
					lock.Lock()
					defer lock.Unlock()
					if pq.Len() > 0 {
						heap.Pop(pq)
						poped++
					}
					///////////////////////
					// waitForPriorityQueue.Pop() 내부에서 크리티컬 섹션을 구현하는 편이 좋음
					// 참고: https://github.com/caffix/queue/blob/master/queue.go
				}()
			}
			wg.Done()
		}(i)
	}

	// wait for all goroutine completed
	wg.Wait()

	// print push, pop action count
	t.Logf("push: %d, pop: %d", pushed, poped)

	/* OUTPUT
	D:\gitworks\go-study\priority-queue>go test -v -run ^TestWaitForPriorityQueueIsConcurrentSafe$
	=== RUN   TestWaitForPriorityQueueIsConcurrentSafe
		waitForPriorityQueue_test.go:132: push: 48849518, pop: 48845173
	--- PASS: TestWaitForPriorityQueueIsConcurrentSafe (10.00s)
	PASS
	ok      gostudy/priority-queue  10.033s
	*/
}
