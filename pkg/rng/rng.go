package rng

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"reflect"

	"github.com/google/uuid"
)

const (
	Alphabet     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Numeric      = "0123456789"
	AlphaNumeric = Alphabet + Numeric
	Hex          = Numeric + "abcdef"
	FileName     = Alphabet + Numeric + ".-_" // POSIX FILENAME
)

// NextBytes 주어진 길이 만큼의 랜덤 바이트 슬라이스를 반환
func NextBytes(length int) (b []byte) {
	b = make([]byte, length)
	rand.Read(b)
	return
}

// NextString length 길이의 임의의 문자열(characterset 문자중)을 반환한다.
// 단, Hex 타입의 경우, 주어진 길이*2 만큼의 문자열이 반환된다. 예: NextString(Hex, 4) = "ffffffff" (4바이트 hex string)
func NextString(characterset string, length int) (ret string) {
	// hex string? (seperated for performance)
	if characterset == Hex {
		// ret = hex.EncodeToString(NextBytes(length))[:length]
		ret = hex.EncodeToString(NextBytes(length))
		return
	}

	// character set length
	max := len(characterset)
	for i := 0; i < length; i++ {
		ret += string(characterset[NextInRange(0, max)])
	}
	return
}

// NextUUID 랜덤 UUID를 반환
func NextUUID() string {
	return uuid.NewString()
}

// number type
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
