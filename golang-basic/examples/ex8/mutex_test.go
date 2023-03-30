package main

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// tick for busy waiting
// const tick = time.Millisecond * 10
const tick = time.Millisecond * 0

// 채널
type Channel[T any] struct {
	q   []T        // 채널 데이터 버퍼
	l   sync.Mutex // 채널 락
	cap int        // 채널 버퍼 크기
}

// 채널 생성
func NewChannel[T any](cap int) *Channel[T] {
	return &Channel[T]{cap: cap}
}

// 채널 락
func (c *Channel[T]) lock() {
	c.l.Lock()
}

// 채널 언락
func (c *Channel[T]) unlock() {
	c.l.Unlock()
}

// 채널에서 데이터를 팝
// o <- ch
func (c *Channel[T]) TryPop() (o T, ok bool) {
	c.lock()
	defer c.unlock()

	if c.Count() > 0 {
		o = c.q[0]
		c.q = c.q[1:]
		return o, true
	}

	return o, false
}

// 채널에 데이터를 푸시
// ch <- o
func (c *Channel[T]) TryPush(o T) bool {
	c.lock()
	defer c.unlock()

	// cap 초과 검사
	if c.Count() >= c.cap {
		return false
	}

	c.q = append(c.q, o)
	return true
}

// 채널의 현재 버퍼 카운트
func (c *Channel[T]) Count() int {
	return len(c.q)
}

// 이중 잠금 문제를 보여주기 위함
func (c *Channel[T]) IsEmpty() bool {
	c.lock()
	defer c.unlock()
	return len(c.q) == 0
}

// 채널에 데이터를 푸시
// ch <- o
func (c *Channel[T]) Push(o T) {
	for {
		if ok := c.TryPush(o); ok {
			return
		}
		// wait for channel ready (busy-waiting)
		time.Sleep(tick)
	}
}

// 채널에서 데이터를 팝(데이터 있을때 까지 대기)
// o <- ch
func (c *Channel[T]) Pop() (o T) {
	for {
		// data found on channel?
		if o, ok := c.TryPop(); ok {
			return o
		}
		// wait for data (busy-waiting)
		time.Sleep(tick)
	}
}

// 데드락이 발생하는 팝
func (c *Channel[T]) TryPopForDeadLock() (o T, ok bool) {
	c.lock()
	defer c.unlock()

	if !c.IsEmpty() { // deadlock here! c.IsEmpty() lock again
		o = c.q[0]
		c.q = c.q[1:]
		return o, true
	}

	return o, false
}

// 동기화 문제가 발생하는 팝
// 이중 잠금 문제는 수정했지만, (1)내부와 (3)에서의 버퍼의 상태(갯수)가 다를 수 있다!
func (c *Channel[T]) TryPopForConcurrentNotSafe() (o T, ok bool) {
	if !c.IsEmpty() { // (1): 버퍼의 크기 > 0
		c.lock()   // (2): 다른 고루틴에 제어권 넘어갈 수 있음
		o = c.q[0] // (3): 버퍼의 크기 == 0 일수도 있음(패닉)
		c.q = c.q[1:]
		c.unlock() // !!(4): c.unlock()을 해주기 전에 panic이 발생하고 콜스택의 상단에서 복구되면,
		// 해당 뮤텍스 자원은 여전히 락이 걸려있는 상태가 되어, 이후 다시 락 진입시 데드락 발생!
		// 따라서, sync.mutex.Unlock() 메소드는 가급적이면 defer 지연 함수와 함께 사용하는 것을 권장
		return o, true
	}
	return o, false
}

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

func TestMutexLockTwice(t *testing.T) {
	l := sync.Mutex{}

	t.Logf("test started!")
	l.Lock()
	l.Lock() // *deadlock*
	t.Logf("test completed!")

	// 다른 언어와 달리 동일한 고루틴에서도 lock이 2번 걸리면 데드락이 발생하기 때문에
	// mutex lock을 사용시 매우 조심해서 사용해야 한다.
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
