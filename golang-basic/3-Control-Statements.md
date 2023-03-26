# 흐름 제어

## if문 (조건문/분기문)

```go
// 조건문 기본 if
if x > max {
    x = max
}

// if-else if-else
if y < min {
    y = min
} else if y > max {
    y = max
} else {
    return y
}

// if 초기문; 조건문 {
if x, err := f(); err != nil {
    fmt.Println(err)
}
// 단, x, err은 if문 밖에서 사용할 수 없음
// x 또는 err이 if문 위에서 선언되어 if문에서 재할당된 경우는 사용할 수 있음
```

## switch문
```go

// 기본형
switch tag {
    default: s3()
    case 0, 1, 2, 3: s1()
    case 4, 5, 6, 7: s2()
}

// 1
switch x := getSomeValue(); {
case x > 0:
    fmt.Println("greater than 0")
case x < 0:
    fmt.Println("less than 0")
default:
    fmt.Println("is zero")
}

// 2. 1은 2와 동일
x := getSomeValue()
switch {
    case x > 0: fmt.Println("greater than 0")
    case x < 0: fmt.Println("less than 0")
    default: fmt.Println("is 0")
}

// 3. 타입 스위치
switch i := x.(type) {
    case nil: fmt.Println("is nil")
    case int: fmt.Println("is int")
    case float64: fmt.Println("is float64")
    case func(int) float64: fmt.Println("is funci(int) float64")
    case bool, string: fmt.Println("type is bool or string")
    default: fmt.Println("unknown type")
}

// 4. fallthough
switch x := 1 {
    case 0: fallthrough
    case 1: fmt.Println("x is 0 or 1")
    default: fmt.Println("x is unknown")
}
```

## for문

```go
// 원형
// for 초기문; 조건문; 후처리 {
//     코드 블록
// }

// 1. 기본형
sum := 0
for i := 0; i < 100; i++ {
    sum += i
}
fmt.Printf("sum(100) is %d\n", sum)

// 2. 초기문 생략
repeats := 100
for ; repeats > 0; repeats-- {
    fmt.Printf("repeats=%d", repeats)
}

// 3. 후처리 생략
for i := 0; i < 100; {
    i += 1 + (i%2)
}

// 4. 조건문 only
repeats := 100
for repeats > 0 {
    repeats--
}

// 5. while문 (w/continue,break)
i := 0
for {
    if i < 10 {
        continue
    }

    if i > 100 {
        break
    }

    fmt.Println(i)
    i++
}

// 6. for range문 (range: 문자열, 배열, 슬라이스, 맵, 채널 등)
values := []string{"hello", "world"}
for i, value := range values {
    fmt.Println(i, value)
}
// output
// 0 hello
// 1 world
```

## select문
채널의 입출력을 선택하여 수행한다.  
POSIX select() 시스템콜과 유사해 보인다. 이름이 select인 이유?  

```go 

var ch1, ch2, ch3 chan int

// ch1, ch2, ch3 모두에 데이터가 들어왔을 때 어떤 채널이 실행될지는 알 수 없다.
// 아마도 조금이라도 먼저 들어온 놈이 실행될듯?
select {
    case c1 := <- ch1: fmt.Println("received ", c1, " from ch1")
    case c2 := <- ch2: fmt.Println("received ", c2, " from ch2")
    case c3 := <- ch3: fmt.Println("received ", c3, " from ch3")
    default: fmt.Println("no message received")
}

```

