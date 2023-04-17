package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// 동기화 문제가 발생하면 패닉이 발생
// 동기화 문제를 확인하려면 Channel.lock(), Channel.unlock() 메소드 내의 코드를 주석처리
func TestChannelShouldBeConcurrentSafe(t *testing.T) {
	ch := NewChannel[int](10000) // 채널 생성
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
				ch.Push(x)                   // ch <- x (push to channel)
				atomic.AddUint64(&pushed, 1) // concurrent safe
			}
			wg.Done()
		}(i)
	}

	// channel pop goroutine
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(x int) {
			for !done {
				ch.Pop()                    // <- ch (pop from channel)
				atomic.AddUint64(&poped, 1) // concurrent safe
			}
			wg.Done()
		}(i)
	}

	// wait for all goroutine completed
	wg.Wait()

	// print push, pop action count
	t.Logf("push: %d, pop: %d", pushed, poped)
}

func TestChannelUsageForSelectLike(t *testing.T) {

	ch1 := NewChannel[int](100)
	ch2 := NewChannel[int](100)

	sender := func(ch *Channel[int], to int) {
		defer ch.Close() // close channel
		for i := 0; i < to; i++ {
			ch.Push(i)
		}
	}

	go sender(ch1, 100)
	go sender(ch2, 100)

	for {
		// select {
		// case v, ok := <- ch1:
		if v, ok := ch1.TryPop(); ok {
			fmt.Println(v)
		}

		// case v, ok := <- ch2:
		if v, ok := ch2.TryPop(); ok {
			fmt.Println(v)
		}
		// }

		// stop the loop
		if ch1.Closed() && ch1.Count() == 0 && ch2.Closed() && ch2.Count() == 0 {
			break
		}
	}
}

func TestMutexLockTwice(t *testing.T) {
	// timeout
	go func() {
		<-time.After(time.Second * 10)
		panic("timeout")
	}()

	l := sync.Mutex{}
	t.Logf("test started!")
	l.Lock()
	l.Lock() // *deadlock*
	t.Logf("test completed!")

	// 다른 언어와 달리 동일한 고루틴에서도 lock이 2번 걸리면 데드락이 발생하기 때문에
	// mutex lock을 사용시 매우 조심해서 사용해야 한다.
}

func TestRWMutexLock(t *testing.T) {
	// timeout
	go func() {
		<-time.After(time.Second * 10)
		panic("timeout")
	}()

	l := sync.RWMutex{}
	t.Logf("test started!")
	l.Lock()
	l.RLock() // ReadLock, WriteLock 간에도 deadlock 발생!
	l.RUnlock()
	l.Unlock()
	t.Logf("test completed!")
}

// 고 채널에 대한 벤치마크
func BenchmarkGoChannel(b *testing.B) {
	ch := make(chan int, 100)
	// ch := NewChannel[int](100)
	wg := sync.WaitGroup{}
	workers := runtime.NumCPU()
	repeats := b.N

	// channel push/pop goroutines
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(x int) {
			for i := 0; i < repeats; i++ {
				ch <- i
				// ch.Push(i)

				<-ch
				// ch.Pop()
			}
			wg.Done()
		}(i)
	}

	// wait for all goroutine completed
	wg.Wait()
}

// 커스텀 채널에 대한 벤치마크
func BenchmarkMyChannel(b *testing.B) {
	// ch := make(chan int, 100)
	ch := NewChannel[int](100)
	wg := sync.WaitGroup{}
	workers := runtime.NumCPU()
	repeats := b.N

	// channel push/pop goroutines
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(x int) {
			for i := 0; i < repeats; i++ {
				// ch <- i
				ch.Push(i)

				// <- ch
				ch.Pop()
			}
			wg.Done()
		}(i)
	}

	// wait for all goroutine completed
	wg.Wait()
}

// 고채널 벤치마크 결과
/*
Running tool: C:\Program Files\Go\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkGoChannel$ github.com/jade-kinx/go-study/golang-basic/examples/ex8 -v

goos: windows
goarch: amd64
pkg: github.com/jade-kinx/go-study/golang-basic/examples/ex8
cpu: 12th Gen Intel(R) Core(TM) i7-12700
BenchmarkGoChannel
BenchmarkGoChannel-20
  302956	      4767 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/jade-kinx/go-study/golang-basic/examples/ex8	1.519s
*/

// 채널 벤치마크 결과
/*
Running tool: C:\Program Files\Go\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkMyChannel$ github.com/jade-kinx/go-study/golang-basic/examples/ex8 -v

goos: windows
goarch: amd64
pkg: github.com/jade-kinx/go-study/golang-basic/examples/ex8
cpu: 12th Gen Intel(R) Core(TM) i7-12700
BenchmarkMyChannel
BenchmarkMyChannel-20
  662766	      1785 ns/op	     319 B/op	       2 allocs/op
PASS
ok  	github.com/jade-kinx/go-study/golang-basic/examples/ex8	2.230s
*/
