# 테스트/벤치마크

## TDD (테스트 주도 개발)


## 테스트 예제

* 테스트 파일은 `*_test.go` 이어야 한다.  
* 테스트 파일은 패키지 빌드시 제외된다.  
* 테스트 함수 이름은 `TestXXX` 여야 한다.

```bash
$ go test -v .
$ go test -v -run ^TestGetBytes$
```

```go
// file: rng.go
package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"reflect"
)

// 랜덤 바이트 배열을 얻는다.
func GetBytes(length int) []byte {
	b := make([]byte, length)
	rand.Read(b)
	return b
}

// number type
type number interface {
	~int8 | ~int16 | ~int32 | ~int64 | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~int | ~uint
}

// Next returns random number in Number type
// Usage: i32, i64 := rng.Next[int32](), rng.Next[int64]()
func Next[T number]() (ret T) {
	b := GetBytes(int(reflect.TypeOf(ret).Size()))
	length := len(b)
	switch length {
	case 1:
		return T(b[0])
	case 2:
		return T(binary.BigEndian.Uint16(b))
	case 4:
		return T(binary.BigEndian.Uint32(b))
	case 8:
		return T(binary.BigEndian.Uint64(b))
	}
	panic(fmt.Errorf("unknown number type size(%d)", length))
}

// NextFromRange returns the random number in range of (min, max)
// output range: min <= value < max
func NextFromRange[T number](min, max T) T {
	// min > max : panic
	if min > max {
		panic(fmt.Errorf("NextInRange(min, max): min(%d) is greater than max(%d)", min, max))
	}
	// min == max : returns min
	if min == max {
		return min
	}

	n := Next[T]() % (max - min)
	if n < 0 {
		return min - n
	}
	return min + n
}
```

```go
// file: rng_test.go
package main

import (
	"math"
	"testing"
	"github.com/stretchr/testify/assert"
)

// GetBytes() 테스트
func TestGetBytes(t *testing.T) {
	assert := assert.New(t)

	// 0. lower 바운더리 체크 ( upper boundary는 없음 )
	assert.Panics(func() { GetBytes(-1) }, "GetBytes(-1) should panic")
	assert.Equal(0, len(GetBytes(0)), "GetBytes(0) should return empty slice")

	// 1. 범위 테스트
	for i := 1; i <= 1024; i++ {
		b := GetBytes(i)

		// expected byte length?
		assert.Equal(i, len(b))
	}
}

// uint8 커버리지 테스트
func TestNextUint8ShouldCoverAll(t *testing.T) {
	assert := assert.New(t)

	repeat := 100
	covered := make([]int, math.MaxUint8+1)
	for i := 0; i < math.MaxUint8*repeat; i++ {
		b := Next[uint8]()
		covered[int(b)]++ // increase hit count
	}

	// 범위내 모든 값이 고르게 분포하는가?
	for i, c := range covered {
		assert.NotZerof(c, "uncovered value: %d", i)
	}
}

// int8 커버리지 테스트
func TestNextInt8ShouldCoverAll(t *testing.T) {
	assert := assert.New(t)

	repeat := 100
	positives := make([]int, math.MaxInt8+1)
	negatives := make([]int, math.MaxInt8+1)
	for i := 0; i < math.MaxUint8*repeat; i++ {
		b := Next[int8]()
		if b >= 0 {
			positives[int(b)]++
		} else {
			negatives[int(-b)%len(negatives)]++ // -128 -> 0
		}
	}

	// 범위내 모든 값이 고르게 분포하는가?
	for i, c := range append(positives, negatives...) {
		assert.NotZerof(c, "uncovered value: %d", i)
	}
}

// uint8 범위 테스트
func TestNextFromRangeUint8(t *testing.T) {
	assert := assert.New(t)

	// 100만번 반복
	for i := 0; i < 1000000; i++ {
		min, max := Next[uint8](), Next[uint8]()

		// min > max면 panic
		if min > max {
			assert.Panicsf(func() { NextFromRange(min, max) }, "NextFromRange(%d, %d) should panic", min, max)
			continue
		}

		// min == max면 min
		if min == max {
			assert.Equalf(min, NextFromRange(min, max), "NextFromRange(%d, %d) should be %d", min, max, min)
			continue
		}

		// min > max면: min <= have < max
		have := NextFromRange(min, max)
		assert.Truef(have >= min && have < max, "NextFromRange(%d, %d) should be (%d <= %d < %d)", min, max, min, have, max)

		// max는 포함하지 않아야 한다
		assert.NotEqualf(max, have, "NextFromRange(%d) should not include max(%d)", have, max)
	}
}

// int8, int16, int32, int64, uint8, uint16, uint32 uint64 모두 테스트해야 함

```


## 벤치마크 예제

* 벤치마크는 테스트 파일(`*_test.go`)에 정의한다.  
* 벤치마크 함수 이름은 `BenchmarkXXX` 여야 한다.

```go
// sum.go
package main

// 0부터 n까지 정수를 더한 값을 반환
func Sum(n int) int {
	sum := 0
	for i := 0; i < n; i++ {
		sum += i
	}
	return sum
}

// 0부터 n까지 정수를 더한 값을 반환
func Sum2(n int) int {
	return n * (n - 1) / 2
}
```

```go
// sum_test.go
package main

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkSum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sum(i)
	}
}
/*
goos: windows
goarch: amd64
pkg: ex9
cpu: 12th Gen Intel(R) Core(TM) i7-12700
BenchmarkSum
BenchmarkSum-20

	1000000	     53803 ns/op	       0 B/op	       0 allocs/op

PASS
ok  	ex9	53.856s
*/

func BenchmarkSum2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sum2(i)
	}
}
/*
goos: windows
goarch: amd64
pkg: ex9
cpu: 12th Gen Intel(R) Core(TM) i7-12700
BenchmarkSum2
BenchmarkSum2-20
1000000000	         0.1063 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	ex9	0.163s
*/

// sum2를 math.MaxInt 까지 하면 overflow
func TestSum2Overflow(t *testing.T) {
	assert := assert.New(t)

	// expected Sum2(math.MaxInt32)
	max := math.MaxInt
	want := max * (max - 1) / 2
	have := Sum2(math.MaxInt32)

	assert.Equalf(want, have, "Sum2(math.MaxInt) should be %d", want)
}
/*
=== RUN   TestSum2Overflow
    d:\gitworks\go-study\golang-basic\examples\ex9\sum_test.go:65:
        	Error Trace:	d:/gitworks/go-study/golang-basic/examples/ex9/sum_test.go:65
        	Error:      	Not equal:
        	            	expected: -4611686018427387903
        	            	actual  : 2305843005992468481
        	Test:       	TestSum2Overflow
        	Messages:   	Sum2(math.MaxInt) should be -4611686018427387903
--- FAIL: TestSum2Overflow (0.00s)
FAIL
FAIL	ex9	0.038s
*/

// 이 경우는 big.Int 를 사용하도록 하거나,
// overflow 발생시 panic 또는 error를 리턴하도록 리팩토링이 필요하다.
```
