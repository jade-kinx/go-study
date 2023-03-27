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
