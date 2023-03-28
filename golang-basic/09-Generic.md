# 제너릭

* Go v1.18 부터 지원  
* v1.18 이전에는 빈 인터페이스`interface{}`를 이용하여 제너릭처럼 사용했지만, `boxing/unboxing` 과정에서 번거로움과 성능 저하 발생  
* Go의 제너릭 타입은 `함수`와 `구조체`에서 사용 가능하고, `메소드`에서는 아직 미지원  
  * 구조체에 정의된 타입 T는 메소드에서 사용 가능하지만, 메소드에서 새로운 타입 U를 정의하여 사용할 수 없음  

```go
import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"reflect"
)

// NextBytes 주어진 길이 만큼의 랜덤 바이트 슬라이스를 반환
func NextBytes(length int) (b []byte) {
	b = make([]byte, length)
	rand.Read(b)
	return
}

// number type 타입 제한자
type number interface {
	~int8 | ~int16 | ~int32 | ~int64 | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~int | ~uint
}

// Next returns random number in Number type
// Usage: i32, i64 := rng.Next[int32](), rng.Next[int64]()
func Next[T number]() (ret T) {
	b := NextBytes(int(reflect.TypeOf(ret).Size()))
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

// NextInRange returns the random number in range of (min, max)
// output range: min <= value < max
func NextInRange[T number](min, max T) T {

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
