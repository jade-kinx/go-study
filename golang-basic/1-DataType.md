# 데이터 타입(built-in/primitive type)

* Golang은 강타입 언어  

## 내장 타입

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

    // 4. 문자열을 문자 배열로 변환
    var str string = "hello, world"
    runes := []rune(str)
}
```

## 오버플로우와 언더플로우
- not ready yet

## 배열
* 배열(array)은 같은 타입을 가지는 변수들의 집합이며, 연속된 메모리 공간에 할당  
* zero-based index로 접근 
* 배열의 크기는 `불변(immutable)`하며, 크기가 다른 배열은 다른 타입으로 인식  
  
```go
var a [5]int    // 크기가 5인 int 타입의 배열. int의 기본값인 0으로 모두 초기화
var b [6]int = [6]int{1, 2, 3, 4, 5, 6} // 크기가 6인 int 배열 초기화
c := [...]string{"one", "two", "three"} // 크기가 3인 string 배열

var d [5]int = a    // 배열의 복사(a와 d는 다른 메모리 공간)
// b = a            // (X) 크기가 다른 배열은 다른 타입

// 다차원 배열의 선언
var mul [2][3]int = [2][3]int{
    {1, 2, 3},
    {4, 5, 6},  // 끝에 콤마(,)가 필요함
}
```

## 슬라이스
* 슬라이스는 `동적 배열(mutable array)`  

```go
type SliceHeader struct {
	Data uintptr    // 데이터의 위치 포인터
	Len  int        // 데이터 현재 길이
	Cap  int        // 할당된 메모리 크기
}
```

### 배열과 슬라이스의 메모리 공간 할당
```go
	// 배열
    arr1 := [5]int{1, 2, 3, 4, 5}
    arr2 := arr1 // 배열의 복사(DeepCopy) (arr1, arr2는 다른 공간을 가리킴)
    arr2[0] = 100
    fmt.Println(arr1)
    fmt.Println(arr2)
    /* OUTPUT
    [1 2 3 4 5]
    [100 2 3 4 5]
    */

	// 슬라이스
	slice1 := []int{1, 2, 3, 4, 5}
	slice2 := slice1 // 헤더만 복사(ShallowCopy)(slice1, slice2는 동일한 공간을 가리킴)
	slice2[0] = 100
	fmt.Println(slice1)
	fmt.Println(slice2)
    /* OUTPUT
    [100 2 3 4 5]
    [100 2 3 4 5]
    */
```

### 슬라이싱과 구조분해할당(javascript)
```go
    // 슬라이싱: 배열 또는 슬라이스에서 일부를 잘라내 슬라이스로 만드는 것
    arr := [5]int{1, 2, 3, 4, 5}
    slice := arr[4:5]
    slice = append(slice, 10, 15, 20, 25)
    fmt.Println(arr)
    fmt.Println(slice)
    /* OUTPUT
    [1 2 3 4 5]
    [5 10 15 20 25]
    */

    // 구조분해할당
    slice2 := append(slice, arr[1:]...) // ... => 전개연산자(spread operator)
    fmt.Println(slice2)
    /* OUTPUT
    [5 10 15 20 25 2 3 4 5]
    */
```

## 문자열

* 문자열(`string`)은 문자의 `배열(array)` (슬라이스가 아님을 주의)  
* 배열이므로 `불변(immutable)`  
* GO는 `UTF8` 인코딩을 문자의 기본으로 사용  

```go
type StringHeader struct {
	Data uintptr    // 데이터의 위치 포인터
	Len  int        // 데이터의 길이
}
```

### 문자열 순회
```go

// C스타일 for문 순회
str := "헬로 월드!!"
for i := 0; i < len(str); i++ {
    fmt.Printf("type: %T, character: %c, code: %d\n", str[i], str[i], str[i])
}
/* OUTPUT
type: uint8, character: í, code: 237
type: uint8, character: , code: 151
type: uint8, character: ¬, code: 172
type: uint8, character: ë, code: 235
type: uint8, character: ¡, code: 161
type: uint8, character: , code: 156
type: uint8, character:  , code: 32
type: uint8, character: ì, code: 236
type: uint8, character: , code: 155
type: uint8, character: , code: 148
type: uint8, character: ë, code: 235
type: uint8, character: , code: 147
type: uint8, character: , code: 156
type: uint8, character: !, code: 33
type: uint8, character: !, code: 33
(1바이트씩 읽으면 한글이 깨진다)
*/

// for range 순회
str := "헬로 월드!!"
for _, c := range str {
    fmt.Printf("type: %T, character: %c, code: %d\n", c, c, c)
}
/* OUTPUT
type: int32, character: 헬, code: 54764
type: int32, character: 로, code: 47196
type: int32, character:  , code: 32
type: int32, character: 월, code: 50900
type: int32, character: 드, code: 46300
type: int32, character: !, code: 33
type: int32, character: !, code: 33
(for-range를 이용하여 순회하면 문자를 rune 타입으로 순회하여 한글 정상 출력)
*/

```

### 문자열 연산

| 연산자 | 설명 |
| :---: | --- |
| + | 더하기 연산자 |
| = | 대입 연산자 |
| == | 문자의 배열이 같음 |
| != | 문자의 배열이 같지 않음 |
| >, >= | 크다, 크거나 같다 |
| <, <= | 작다, 작거나 같다 |

* 문자열은 배열이므로 immutable 하고, 문자열에 다른 문자열을 추가하는 등의 연산이 일어날 경우, 새로운 객체가 생성

## 타입 정의
```go
// 별칭(alias)
type MyInt int

// 구조체 정의
type MyStruct struct {
    Id int
    Name string
}

// 인터페이스 정의
type MyError interface {
    error   // 내장 에러 임베딩
    Code() int
}

```
