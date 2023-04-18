package main

import (
	"fmt"
	"sync"
	"time"
)

// 채널 데이터 래퍼
type element[T any] struct {
	data     T   // 데이터
	priority int // 우선순위
}

// 우선순위 채널
type PriorityChannel[T any] struct {
	q      []element[T]  // 채널 아이템 버퍼
	l      sync.RWMutex  // 채널 락
	cap    int           // 채널 버퍼 크기
	closed bool          // 채널 종료
	tick   time.Duration // tick for busy-waiting
}

// 채널 생성
func NewChannel[T any](cap int) *PriorityChannel[T] {
	return &PriorityChannel[T]{cap: cap}
}

func (c *PriorityChannel[T]) SetTick(tick time.Duration) {
	c.tick = tick
}

// 채널 종료
func (c *PriorityChannel[T]) Close() {
	c.closed = true
}

// 쓰기 락
func (c *PriorityChannel[T]) wLock() {
	c.l.Lock()
}

// 쓰기 언락
func (c *PriorityChannel[T]) wUnlock() {
	c.l.Unlock()
}

// 읽기 락
func (c *PriorityChannel[T]) rLock() {
	c.l.RLock()
}

// 읽기 언락
func (c *PriorityChannel[T]) rUnlock() {
	c.l.RUnlock()
}

//////////////////////////////////////////////////////////////////////////////////
// 읽기락/읽기언락은 여기서는 굳이 필요없지만, read/write lock에 대해서 설명하기 위함
//////////////////////////////////////////////////////////////////////////////////

// 채널의 현재 원소 갯수
func (c *PriorityChannel[T]) Count() int {
	c.rLock()
	defer c.rUnlock()
	return len(c.q)
}

// 채널에서 데이터를 하나 꺼낸다
func (c *PriorityChannel[T]) TryPop() (item element[T], ok bool) {
	c.wLock()
	defer c.wUnlock()

	if len(c.q) > 0 {
		item = c.q[0]
		c.q = c.q[1:]
		return item, true
	}

	return item, false
}

// 채널에 데이터를 추가
func (c *PriorityChannel[T]) TryPush(data T, priority int) bool {
	c.wLock()
	defer c.wUnlock()

	// 채널이 닫혔으면 실패 처리
	if c.closed {
		return false
	}

	// cap 초과 검사
	if len(c.q) >= c.cap {
		return false
	}

	// 새로운 아이템
	item := element[T]{data: data, priority: priority}

	/////////////////////////////////////////
	// 새 아이템을 우선순위 위치에 추가
	// heap을 이용하는 편이 성능상 우수하다
	at := c.find(priority)
	if len(c.q) == at {
		c.q = append(c.q, item)
		return true
	}
	c.q = append(c.q[:at+1], c.q[at:]...)
	c.q[at] = item
	return true
}

// priority가 위치할 index를 찾는다.
func (c *PriorityChannel[T]) find(priority int) int {
	for i, item := range c.q {
		if item.priority > priority {
			return i
		}
	}

	return len(c.q)
}

// 채널에 데이터를 푸시(입력 완료까지 대기)
// ch <- data
func (c *PriorityChannel[T]) Push(data T, priority int) error {
	for {
		// channel closed?
		if c.closed {
			return fmt.Errorf("push to closed channel")
		}

		// push data
		if ok := c.TryPush(data, priority); ok {
			return nil
		}

		// wait for channel ready (busy-waiting)
		time.Sleep(c.tick)
	}
}

// 채널에서 데이터를 팝(데이터 있을때 까지 대기)
// data <- ch
func (c *PriorityChannel[T]) Pop() (data T, err error) {
	data, _, err = c.PopWithPriority()
	if err != nil {
		return data, err
	}
	return data, nil
}

// 채널에서 데이터를 우선순위와 함께 팝(데이터 있을때 까지 대기)
func (c *PriorityChannel[T]) PopWithPriority() (data T, priority int, err error) {
	for {
		// channel closed?
		if c.closed {
			return data, priority, fmt.Errorf("priority channel closed")
		}

		// data available on p-channel?
		if item, ok := c.TryPop(); ok {
			return item.data, item.priority, nil
		}

		// wait for data (busy-waiting)
		time.Sleep(c.tick)
	}
}
