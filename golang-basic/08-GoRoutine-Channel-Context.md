# 고루틴/채널/컨텍스트

## 고루틴
* 고루틴은 고 런타임이 제공하는 경량 스레드  
* 스레드 보다 작은 단위(2KB 스택)의 태스크  
* `go` 키워드로 함수를 멀티 스레드로 실행하는 구문으로 매우 심플  

```go
func greeting(to string) {
    fmt.Printf("Hello, %s!\n", to)
}

func main() {
    wg := sync.WaitGroup{}

    // 고루틴 실행
    go gretting("Go")   // 실행이 될수도 있고, 안될 수도 있다. (고런타임의 스케쥴러가 결정)

    wg.Add(1)
    // 익명함수로 고루틴 실행
    go func(to string) {
        fmt.Printf("Hello, %s!\n", to)
        wg.Done()
    }("Go")

    // 메인 고루틴이 종료되면 모든 서브 고루틴도 종료된다.
    wg.Wait()
}
```

### 멀티스레딩 환경에서의 동시성 문제
* 과거의 게임 서버들이 걸핏하면 뻗던 문제가 대부분 동기화 처리 문제  
* 원자성/데드락/라이브락/기아상태
* Golang에서는 기아상태 문제로 select 키워드의 순서를 균등 분포 확률로 처리  
* 심도 있게 다루어야 하는 주제  

동기화의 필요성 (크리티컬 섹션/상호 배제)
```go
func main() {
    var data int
    go func() { data++ }()
    if data == 0 {
        fmt.Printf("data is 0\n")   // 그런데 진짜로 0?
        // fmt.Printf("data is %d\n", data)
    } else {
        fmt.Printf("data is %d\n", data)
    }
    // ...
}

// if data == 0 {} 내부의 data가 0가 아님을 알 수 있는 테스트 코드
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
		if zero == 0 {  // (1): zero의 값이 0이라고 비교하고 들어왔는데
			time.Sleep(time.Millisecond * 10) // sleep 10ms
			if zero != 0 {  // (2): zero가 0이 아니라고? 아니 이게 무슨소리야?
				panic("zero is not 0")
			}
		}
	}
}
```

데드락 유발 예제
```go
type value struct {
    mu sync.Mutex
    val int
}

func main() {
    var wg sync.WaitGroup{}
    printSum := func( x, y *value) {
        defer wg.Done()
        
        // x락
        x.mu.Lock()
        defer x.mu.Unlock()

        // 데드락 유발을 위해 잠시 대기
        time.Sleep(time.Second * 2)

        // y락
        y.mu.Lock()
        defer y.mu.Unlock()

        fmt.Printf("%d+%d=%d\n", x.val, y.val, x.val + y.val)
    }

    var x, y value 
    wg.Add(2)
    go printSum(&x, &y)
    go printSum(&y, &x)
    wg.Wait()
}
```


## 채널
* 채널은 동기화 기능(concurrent-safe)을 제공하는 메세지 큐  
* 고루틴간 흐름을 제어할 수 있도록 한다.  
* `select` 키워드로 채널들을 바인딩하여 기아 상태를 예방  

채널의 구현 샘플 참고: [/golang-basic/examples/ex8/mutex_test.go](https://github.com/jade-kinx/go-study/blob/main/golang-basic/examples/ex8/mutex_test.go)

채널의 기본적인 사용법
```go
// N senders -> 1 receiver
func ChannelForSendersToReceiver() {
	channel := make(chan string) // 채널
	wgs := sync.WaitGroup{}      // 송신자 wg
	wgr := sync.WaitGroup{}      // 수신자 wg

	// 채널 송신자 N
	for i := 0; i < runtime.NumCPU(); i++ {
		wgs.Add(1)
		go func(ch chan<- string, id int) { // write-only 채널
			defer fmt.Printf("sender[%d] terminated\n", id)
			defer wgs.Done()
			for counter := 0; counter < 10; counter++ {
				// 채널에 메세지 입력
				ch <- fmt.Sprintf("sender[%d]: Hello, Go! (%d)", id, counter)
			}
		}(channel, i)
	}

	// 채널 수신자 1
	wgr.Add(1)
	go func(ch <-chan string) { // read-only 채널
		defer fmt.Println("receiver terminated")
		defer wgr.Done()
		// 채널에서 메세지를 수신 (채널이 닫히면 for loop 종료)
		for msg := range ch {
			fmt.Println(msg)
		}
	}(channel)

	// wait for channel senders
	wgs.Wait()

	// close channel
	close(channel)

	// wait for channel receiver
	wgr.Wait()
}

// 1 sender -> N receivers
func ChannelForSenderToReceivers() {
	channel := make(chan string)
	wgs := sync.WaitGroup{} // 송신자 wg
	wgr := sync.WaitGroup{} // 수신자 wg

	// 채널 송신자 1
	wgs.Add(1)
	go func(ch chan<- string) { // write-only 채널
		defer fmt.Println("sender terminated")
		defer wgs.Done()
		for counter := 0; counter < 10; counter++ {
			// 채널에 메세지 입력
			ch <- fmt.Sprintf("Hello, Go! (%d)", counter)
		}
	}(channel)

	// 채널 수신자 N (어떤 수신자가 메세지를 받을지는 고런타임 스케쥴러가 결정)
	for i := 0; i < runtime.NumCPU(); i++ {
		wgr.Add(1)
		go func(ch <-chan string, id int) { // read-only 채널
			defer fmt.Printf("receiver[%d] terminated\n", id)
			defer wgr.Done()
			// 채널에서 메세지를 수신 (채널이 닫히면 for loop 종료)
			for msg := range ch {
				fmt.Printf("receiver[%d]: %s\n", id, msg)
			}
		}(channel, i)
	}

	// wait for channel senders
	wgs.Wait()

	// close channel
	close(channel)

	// wait for channel receiver
	wgr.Wait()
}

```


## 컨텍스트
* 컨텍스트는 스레드(또는 고루틴) 마다 가지는 고유의 개별 데이터  
* 컨텍스트는 스레드의 맥락! (예: `그제 회의에서 이야기한 메모리 공간 재할당 문제`)  
  * `gin` context = request, response, 기타 변수 등
* 스레드 전환 발생시 컨텍스트 스위칭 비용(overhead) 발생(예: `go언어로 작업하다가 python 언어로 작업해야 할때`)  
* Golang에서 컨텍스트(context)는 작업 명세서와 같은 역할로 작업 가능한 시간, 작업 취소 등 작업의 흐름을 제어하는 데 사용  
* C#에서의 CancellationToken과 같이 고루틴 작업의 취소에 주로 사용  

### context.WithValue
```go
func myFunc(ctx context.Context) {
    if v := ctx.Value("user"); v != nil {
        user, ok := v.(string)
        if !ok {
            fmt.Println("user should be string type")
            return
        }
        fmt.Printf("Hello, %s!\n", user)
        return
    }

    fmt.Println("context value 'user' not found")
}

func main() {
    ctx := context.Background()

    // 컨텍스트에 값 추가
    ctx = context.WithValue(ctx, "user", "john-doe")

    myFunc(ctx)
}
```

### context.WithCancel
```go

// 아주 오래 걸리는 함수가 있다고 가정
func doLongTermWork() string {
    for {
        if r := rand.Intn(100); r == 0 {
            return "ok"
        }
        time.Sleep(time.Second)
    }
}

// 아주 오래 걸리는 함수를 취소할 수 있도록 래핑
func runLongTermWork(ctx context.Context) (string, error) {
    done := make(chan string)

    go func() {
        done <- doLongTermWork()
    }()

    select {
        case result := <- done:
            return result, nil
        case <-ctx.Done():
            return "fail", ctx.Err()
    }
}

func main() {

    // context와 context에 취소 요청을 보낼 수 있는 인터페이스(함수)
    ctx, cancel := context.WithCancel(context.Background())

    // 특정 조건에서 고루틴의 실행을 취소하고 싶다
    go func() {
        <- time.After(time.Second * 10)
        cancel()
    }()

    result, err := runLongTermWork(ctx)
    fmt.Printf("result is %s\n", result)
}
```

