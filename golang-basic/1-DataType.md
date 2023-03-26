# 데이터 타입(built-in/primitive type)

## 데이터 타입

| 타입 | 설명 | 범위 |
| --- | --- | --- |
| bool | 참 또는 거짓 | {true, false} |
| uint8 | 1바이트 부호 없는 정수 | 0 ~ (2^8)-1 |
| uint16 | 2바이트 부호 없는 정수 | 0 ~ (2^16)-1 |
| uint32 | 4바이트 부호 없는 정수 | 0 ~ (2^32)-1 |
| uint64 | 8바이트 부호 없는 정수 | 0 ~ (2^64)-1 |
| int8 | 1바이트 부호 있는 정수 | -(2^7) ~ (2^7)-1 |
| int16 | 2바이트 부호 있는 정수 | -(2^15) ~ (2^15)-1 |
| int32 | 4바이트 부호 있는 정수 | -(2^31) ~ (2^31)-1 |
| int64 | 8바이트 부호 있는 정수 | -(2^63) ~ (2^63)-1 |
| float32 | 4바이트 실수 | IEEE-754 32비트 실수 |
| float64 | 8바이트 실수 | IEEE-754 64비트 실수 |
| complex64 | 8바이트 복소수(진수,가수) | 진수와 가수 범위는 float32 범위와 동일 |
| complex128 | 16바이트 복소수(진수,가수) | 진수와 가수 범위는 float64 범위와 동일 |
| byte | uint8 별칭 | uint8 범위와 동일 |
| rune | int32 별칭 (UTF-8 문자 하나를 표현) | int32 범위와 동일 |
| int | 32비트 컴퓨터 int32, 64비트 컴퓨터 int64 | |
| uint | 32비트 컴퓨터 uint32, 64비트 컴퓨터 uint64 | |
| uintptr | 메모리 포인터 타입 | |
| string | 문자열 | `|
| struct | 구조체 | `type User struct { id string }` |
| [n]array | 배열 | `var numbers [10]int` |
| []slice | 슬라이스 | `var names []string` |
| map[k]v | 맵 | `var phones map[string]string` |
| chan | 채널 | `var ch chan string` |

## 변수 선언

```go
package main

import "fmt"

func main() {
    var a int = 3   // 1. 변수 선언 기본 형태
    var b int       // 2. 변수 선언 초기값 생략 (타입의 기본값: 0, 0.0, false, "", nil)
    var c = a       // 3. 타입 생략 (c의 타입은 우변의 타입)
    d := 3          // 4. 선언 대입문 (for convinience)
    fmt.Println(a, b, c, d) // output: 3 0 3 3

    // var d                // X : d의 타입을 알 수 없음
    // d := 100             // X : d는 이미 선언되어 있으므로 선언 대입할 수 없음
    d, e := 100, "hello"    // d재할당, e는 선언 대입
    fmt.Println(d, e)       // output: 100 hello
}
```

## 타입 캐스팅
타입명(변수) 형태로 타입 변환
```go
package main

import "fmt"

func main() {

    // 1. 타입 크기가 작은 타입으로 변환시 바이트 버림
    var i16 int16 = 3456    // 00001101 10000000
    var i8 int8 = int8(i16) // -------- 10000000
    fmt.Println(i8, i16)    // output: -128 3456

    // 2. 실수에서 정수로 변환시 소수점 이하 버림
    var pi float64 = 3.141592
    var intpi int64 = int64(pi)
    fmt.Println(pi, intpi)  // output: 3.141592 3

    // 3. 실수에서 정수로 크기가 작은 타입으로 변환 (2+1)
    var x float64 = 3456.012345789
    var y int8 = int8(x)
    fmt.Println(x, y)       // output: 3456.0123456789 -128
}
```
