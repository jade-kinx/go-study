package main

import (
	"fmt"
	"sync"
	"time"
)

// tick for busy waiting
// const tick = time.Millisecond * 10
var (
	tick = time.Millisecond * 0
)

// 채널
type Channel[T any] struct {
	q      []T        // 채널 데이터 버퍼
	l      sync.Mutex // 채널 락
	cap    int        // 채널 버퍼 크기
	closed bool       // 채널 종료
}

// 채널 생성
func NewChannel[T any](cap int) *Channel[T] {
	return &Channel[T]{cap: cap}
}

// 채널 종료
func (c *Channel[T]) Close() {
	c.closed = true
}

func (c *Channel[T]) Closed() bool {
	return c.closed
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

	// 채널 종료 검사
	if c.closed {
		return false
	}

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
func (c *Channel[T]) Push(o T) error {
	for {
		// channel closed?
		if c.closed {
			return fmt.Errorf("push to closed channel")
		}

		// push data
		if ok := c.TryPush(o); ok {
			return nil
		}

		// wait for channel ready (busy-waiting)
		time.Sleep(tick)
	}
}

// 채널에서 데이터를 팝(데이터 있을때 까지 대기)
// o <- ch
func (c *Channel[T]) Pop() (o T, err error) {
	for {
		// channel closed?
		if c.closed {
			return o, fmt.Errorf("channel closed")
		}

		// data found on channel?
		if o, ok := c.TryPop(); ok {
			return o, nil
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
