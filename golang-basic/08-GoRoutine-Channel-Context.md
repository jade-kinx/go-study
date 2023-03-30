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
* 원자성(atomic operation)
* 데드락/라이브락/기아상태
* 심도 있게 다루어야 하는 주제  

동기화의 필요성
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

[/golang-basic/examples/ex8/mutex_test.go](https://github.com/jade-kinx/go-study/blob/main/golang-basic/examples/ex8/mutex_test.go) 참고

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

