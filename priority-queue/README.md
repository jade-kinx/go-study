# 채널과 우선순위 큐(Priority Queue)의 구현

## 우선선위 큐 개요

> 우선순위 큐(Priority Queue)는 평범한 `큐`나 `스택`과 비슷한 축약 자료형이다. 그러나 각 원소들은 우선순위를 갖고 있다. 우선순위 큐에서, 높은 우선 순위를 가진 원소는 낮은 우선순위를 가진 원소보다 먼저 처리된다. 만약 두 원소가 같은 우선순위를 가진다면 그들은 큐에서 그들의 순서에 의해 처리된다.  
[위키백과](https://ko.wikipedia.org/wiki/%EC%9A%B0%EC%84%A0%EC%88%9C%EC%9C%84_%ED%81%90)

| 컨테이너 | 설명 |
| :---: | --- |
| 스택 | 후입선출(FILO) |
| 큐 | 선입선출(FIFO) |
| 우선순위 큐 | 우선순위 순으로 처리하되, 우선순위가 같으면 선입선출 |

* 우선순위 큐는 큐 자료 구조에서 필요에 의해 우선순위 개념을 추가한 것으로, 필요하다면 우선순위 스택을 만들 수도 있다.  


### 우선순위 큐의 활용
* 작업 스케쥴링(CPU, Thread, Task, co-routine 등)
* 네트워크 QoS(Quality of Service), OOB 등
* 트랜잭션(DB, BlockChain 등) 처리 순서 조정 등
* API 요청(http request)에 대해 우선순위 부여하여 처리 등
* 큐의 구조를 갖되 우선 처리를 요하는 모든 작업(은행/마트 등 대기시간)


### 구현 방식에 따른 시간복잡도(BigO 표기법)

| 구현 방법 | enque | deque |
| :---: | --- | --- |
| 배열(unsorted) | O(1) | O(N) |
| 배열(sorted) | O(N) | O(1) |
| 연결 리스트(unsorted) | O(1) | O(N) |
| 연결 리스트(sorted) | O(N) | O(1) |
| 힙(heap) | O(logN) | O(logN) |

일반적으로 우선순위 큐는 [힙(자료구조)](https://ko.wikipedia.org/wiki/%ED%9E%99_(%EC%9E%90%EB%A3%8C_%EA%B5%AC%EC%A1%B0))를 통해 구현되지만, 반드시 그럴 필요는 없고, 상황에 따라 적당한 방법으로 구현해도 무방하다.  

### 기본 인터페이스

| 함수 | 설명 |
| :---: | --- |
| enque(push) | 우선순위 큐에 원소를 추가한다 |
| deque(pop) | 우선순위 큐에서 원소를 하나 꺼낸다(제거) |
| peek(top) | 우선순위 큐의 맨 앞에 있는 원소(우선순위가 가장 높은)를 반환(제거X) |
| size(len) | 우선순위 큐에 보관되어 있는 원소의 수를 반환한다 |
| empty | (optional) 우선순위 큐가 비어있는지 확인한다 |

## Go 채널을 이용한 우선순위 큐의 구현

```go
// 2개의 채널로 우선순위를 구분하여 높은 우선순위의 채널을 먼저 읽는다.
func SelectWithPriority(highChan <-chan int, lowChan <-chan int) int {
    select {
        case val := <-highChan:
            return val
        default:
            select {
                case val := <-highChan:
                    return val
                case val := <-lowChan:
                    return val
            }
    }
}
```

* High-Priority, Normal-Priority, Low-Priority 3개의 우선순위를 갖는 경우?
* 우선순위 항목이 가변적인 경우? (대기시간, 파일크기 등)
* Go 채널만으로는 우선순위 큐를 구현하기 적합하지 않다.(채널은 기본적으로 순서를 보장)  
* 일반적인 큐에 우선순위 처리 로직을 더하고, `concurrent-safe` 처리를 해주면 `우선순위 채널`을 만들 수 있지 않을까?  

## Go의 heap 컨테이너를 이용한 우선순위 큐

go.dev 예제: [Example(PriorityQueue)](https://pkg.go.dev/container/heap)
```go
// This example demonstrates a priority queue built using the heap interface.
package main

import (
	"container/heap"
	"fmt"
)

// An Item is something we manage in a priority queue.
type Item struct {
	value    string // The value of the item; arbitrary.
	priority int    // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, value string, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

// This example creates a PriorityQueue with some items, adds and manipulates an item,
// and then removes the items in priority order.
func main() {
	// Some items and their priorities.
	items := map[string]int{
		"banana": 3, "apple": 2, "pear": 4,
	}

	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := make(PriorityQueue, len(items))
	i := 0
	for value, priority := range items {
		pq[i] = &Item{
			value:    value,
			priority: priority,
			index:    i,
		}
		i++
	}
	heap.Init(&pq)

	// Insert a new item and then modify its priority.
	item := &Item{
		value:    "orange",
		priority: 1,
	}
	heap.Push(&pq, item)
	pq.update(item, item.value, 5)

	// Take the items out; they arrive in decreasing priority order.
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		fmt.Printf("%.2d:%s ", item.priority, item.value)
	}

    /* OUTPUT
    05:orange 04:pear 03:banana 02:apple
    */
}
```

## 쿠버네티스 client-go의 PriorityQueue 예
깃허브 링크: [https://github.com/kubernetes/client-go/blob/master/util/workqueue/delaying_queue.go](https://github.com/kubernetes/client-go/blob/master/util/workqueue/delaying_queue.go)

* [waitForPriorityQueue.go](./waitForPriorityQueue.go)  
* [waitForPriorityQueue_test.go](./waitForPriorityQueue_test.go)  

```go
// waitforPriorityQueue_test.go

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
					defer lock.Unlock() // unlock은 반드시 defer와 함께!
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
					defer lock.Unlock() // unlock은 반드시 defer와 함께
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

```


## 우선순위 채널의 구현

* [priority_channel.go](./priority_channel.go)
```go
////////////////////////
// priority_channel.go
////////////////////////

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

// 새 원소의 priority가 위치할 index를 찾는다.
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
```

* [priority_channel_test.go](./priority_channel_test.go)
```go
/////////////////////////////
// priority_channel_test.go
/////////////////////////////

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
```

### 추가적으로 해볼만 한 작업들
* `heap`을 이용한 구현(성능 개선)
* 이미 포함된 원소의 우선순위 변경
* 원소의 삭제
* 양방향(우선순위) 출력
* `queue`, `stack`, `map`, `list` 등의 자료 구조를 `concurrent-safe`하게 만들기  




