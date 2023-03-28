# 함수

## 형태
```go
// func: 함수 키워드
// Sum: 함수 이름
// (a int, b int): 함수 인자
// int: 리턴타입
func Sum(a int, b int) int {
    return a + b
}
```
* 함수의 `visibility`를 제한하는 `modifier keyword(public, private 등)`가 없음  
* 함수 이름의 첫글자가 대문자이면 패키지 외부에 대해 `public`, 그렇지 않으면 `private`
  
## Pass-By-Value or Pass-By-Reference?
* Golang에서는 기본적으로 모든 타입에 대해 `Pass-By-Value`  
* `Pass-By-Value`는 메모리 복사가 발생, `Pass-By-Reference` 형태로 사용하기 위해서는 `Pointer`를 전달

```go
func IncrementValue(a int) {
    a += 100
}
func IncrementPointer(a *int) {
    *a += 100
}
func main() {
    a := 100
    
    // pass-by-value: a와 incrementValue() 내부의 a는 다른 메모리 공간
    IncrementValue(a)
    fmt.Println(a)
    /* OUTPUT
    100
    */

    // pass-by-pointer: a의 포인터를 전달, IncrementPointer() 내부의 a와 같은 메모리 공간
    IncrementPointer(&a)
    fmt.Println(a)
    /* OUTPUT
    200
    */
}
```

## 가변인자 함수
```go
func ExceptForZero(numbers ...int) (rst []int) {
    for _, n := range numbers {
        if n != 0 {
            rst = append(rst, n)
        }
    }
    return
}

func main() {
    nonZeroes := ExceptForZero(0, 1, 2, 0, 3, 4, 0, 5, 0)
    fmt.Println(nonZeroes)
    /* OUTPUT
    [1 2 3 4 5]
    */
}

```

## 함수 리턴 구문과 지연 실행
```go
// 다중 리턴 구문
// C#의 튜플과 비슷한듯 다름
func divide(a, b int) (c int, err error) {
    if b == 0 {
        err = fmt.Errorf("`%d` can't divide by 0", a)
        return  // c는 int의 zero-value: 0
    }

    if a == 0 {
        return 0, nil
    }

    c = a / b
    return  // err은 error의 zero-value: nil
}

func main() {
    defer fmt.Println("나중 실행") // 함수가 종료될 때 실행, 주로 자원 반납 목적으로 사용
    defer fmt.Println("먼저 실행") // 지연함수는 스택구조, 나중에 선언된 실행이 먼저 실행

    c, err := divide(5, 3)
    fmt.Println(c, err)
    c, err = divide(5, 0)
    fmt.Println(c, err)
    x, _ := divide(5, 0)
    fmt.Println(x)
    /* OUTPUT
    1 <nil>
    0 `5` can't divide by 0
    0
    먼저 실행
    나중 실행    
    */
}
```


## 함수리터럴(익명함수/람다)

```go
// 의존성 주입(dependency injection)
func ExceptForZero(numbers []int, f func(int) bool) (rst []int) {
    for _, n := range numbers {
        if !f(n) {
            rst = append(rst, n)
        }
    }
    return
}

func main() {
    numbers := []int{0, 1, 2, 0, 3, 4, 0, 5, 0}
    nonZeroes := ExceptForZero(numbers, func(n int) bool {
        return n == 0
    })
    fmt.Println(nonZeroes)
    /* OUTPUT
    [1 2 3 4 5]
    */

    go func(i int) {
        fmt.Println(i)  // 메인 함수 종료로 출력 안됨
    }(100)
}
```