package rng

import (
	"math"
	"reflect"
	"regexp"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

func TestNextBytes(t *testing.T) {
	t.Logf("TestNextBytes() started")
	for i := 0; i < 1024*1024; i++ {
		actual := NextBytes(i)
		expected := i
		if len(actual) != expected {
			t.Errorf("TestNextBytes(): expected=%d, actual=%d", expected, actual)
		}
	}
	t.Logf("TestNextBytes() completed")
}

// s is alphabet string?
func IsAlphabet(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}

	return true
}

// s is numeric string?
func IsNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsNumber(r) {
			return false
		}
	}

	return true
}

// s is alphabet or numeric string?
func IsAlphaNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return false
		}
	}

	return true
}

// s is hex string? (except for prefix '0x')
func IsHex(s string) bool {
	matched, _ := regexp.MatchString(`^[0-9a-fA-F]+$`, s)
	return matched
}

// s is posix safe filename?
func IsFileName(s string) bool {
	matched, _ := regexp.MatchString(`^[0-9a-zA-Z._-]+$`, s)
	return matched
}

// test NextString()
func TestNextString(t *testing.T) {
	assert := assert.New(t)
	repeats := 10000

	// 0. NextString(x, 0) should return ""
	assert.Equal("", NextString(Alphabet, 0))
	assert.Equal("", NextString(Numeric, 0))
	assert.Equal("", NextString(AlphaNumeric, 0))
	assert.Equal("", NextString(Hex, 0))
	assert.Equal("", NextString(FileName, 0))

	// 1. alphabet only
	for i := 0; i < repeats; i++ {
		expectedLength := NextInRange(1, 1024)
		actual := NextString(Alphabet, expectedLength)

		// expected length? is alphabet?
		assert.Equalf(expectedLength, len(actual), "expected=%d, actual=%d", expectedLength, actual)
		assert.Truef(IsAlphabet(actual), "actual=%s, expected=%s", actual, Alphabet)
	}

	// 2. numeric only
	for i := 0; i < repeats; i++ {
		expectedLength := NextInRange(1, 1024)
		actual := NextString(Numeric, expectedLength)

		// expected length? is numeric string?
		assert.Equalf(expectedLength, len(actual), "expected=%d, actual=%d", expectedLength, actual)
		assert.Truef(IsNumeric(actual), "actual=%s, expected=%s", actual, Numeric)
	}

	// 3. alphanumeric
	for i := 0; i < repeats; i++ {
		expectedLength := NextInRange(1, 1024)
		actual := NextString(AlphaNumeric, expectedLength)

		// expected length? is alphabet or numeric string?
		assert.Equalf(expectedLength, len(actual), "expected=%d, actual=%d", expectedLength, actual)
		assert.Truef(IsAlphaNumeric(actual), "actual=%s, expected=%s", actual, AlphaNumeric)
	}

	// 4. hex string
	for i := 0; i < repeats; i++ {
		length := NextInRange(1, 1024)
		actual := NextString(Hex, length)
		expected := length * 2

		// expected length? is hex string?
		assert.Equalf(expected, len(actual), "expected=%d, actual=%d", expected, actual)
		assert.Truef(IsHex(actual), "actual=%s, expected=%s", actual, Hex)
	}

	// 5. filename string
	for i := 0; i < repeats; i++ {
		expectedLength := NextInRange(1, 1024)
		actual := NextString(FileName, expectedLength)

		// expected length? is filename?
		assert.Equalf(expectedLength, len(actual), "expected=%d, actual=%d", expectedLength, actual)
		assert.Truef(IsFileName(actual), "actual=%s, expected=%s", actual, FileName)
	}
}

func BenchmarkNextStringAlphabet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NextString(Alphabet, 1024)
	}
}

func BenchmarkNextStringNumeric(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NextString(Numeric, 1024)
	}
}

func BenchmarkNextStringAlphaNumeric(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NextString(AlphaNumeric, 1024)
	}
}

func BenchmarkNextStringHex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NextString(Hex, 1024)
	}
}

func BenchmarkNextStringFileName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NextString(FileName, 1024)
	}
}

func TestNextUUID(t *testing.T) {
	t.Logf("TestNextUUID() started")
	for i := 0; i < 1024*128; i++ {
		actual := NextUUID()

		// length check
		if len(actual) != 36 {
			t.Errorf("TestNextUUID(): unexpected length: %d", len(actual))
		}

		// is UUID?
		if match, err := regexp.MatchString("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$", actual); err != nil || !match {
			t.Errorf("TestNextUUID(): match=%v, err=%s", match, err)
		}
	}
	t.Logf("TestNextUUID() completed")
}

func TestNextInt8(t *testing.T) {
	t.Logf("TestNextInt8() started")
	for r := 0; r < 100; r++ {
		positive := make([]int, math.MaxInt8+1)
		negative := make([]int, math.MaxInt8+1)

		// no-hit 확률은 0.1%
		for count := 0; count < math.MaxUint8*1000; count++ {
			// next int8
			r := Next[int8]()
			if r >= 0 {
				positive[r]++
			} else {
				negative[uint(r)%uint(len(negative))]++
			}
		}

		// hit test
		for i, c := range positive {
			if c == 0 {
				t.Errorf("NextInt8(): missing hit value for positive %d", i)
			}
		}
		for i, c := range negative {
			if c == 0 {
				t.Errorf("NextInt8(): missing hit value for negative %d", i)
			}
		}
	}
	t.Logf("TestNextInt8() completed")
}

func TestNextUint8(t *testing.T) {
	t.Logf("TestNextUint8() started")
	for r := 0; r < 100; r++ {
		positive := make([]int, math.MaxUint8+1)

		// no-hit 확률은 0.1%
		for count := 0; count < math.MaxUint8*1000; count++ {
			// next uint8
			r := Next[uint8]()
			positive[r]++
		}

		// hit test
		for i, c := range positive {
			if c == 0 {
				t.Errorf("NextUint8(): missing hit value for positive %d", i)
			}
		}
	}
	t.Logf("TestNextUint8() completed")
}

func TestNextInt16(t *testing.T) {
	t.Logf("TestNextInt16() started")
	for r := 0; r < 10; r++ {
		positive := make([]int, math.MaxInt16+1)
		negative := make([]int, math.MaxInt16+1)

		// no-hit 확률은 0.1%
		for count := 0; count < math.MaxUint16*1000; count++ {
			// next int16
			r := Next[int16]()
			if r >= 0 {
				positive[int(r)%len(positive)]++
			} else {
				negative[uint(r)%uint(len(negative))]++
			}
		}

		// hit test
		for i, c := range positive {
			if c == 0 {
				t.Errorf("NextInt16(): missing hit value for positive %d", i)
			}
		}
		for i, c := range negative {
			if c == 0 {
				t.Errorf("NextInt16(): missing hit value for negative %d", i)
			}
		}
	}
	t.Logf("TestNextInt16() completed")
}

func TestNextUint16(t *testing.T) {
	t.Logf("TestNextUint16() started")
	for r := 0; r < 10; r++ {
		positive := make([]int, math.MaxInt16+1)

		// no-hit 확률은 0.1%
		for count := 0; count < math.MaxUint16*1000; count++ {
			// next uint16
			r := Next[uint16]()
			positive[int(r)%len(positive)]++
		}

		// hit test
		for i, c := range positive {
			if c == 0 {
				t.Errorf("NextUint16(): missing hit value for positive %d", i)
			}
		}
	}
	t.Logf("TestNextUint16() completed")
}

func TestNextInt32(t *testing.T) {
	t.Logf("TestNextInt32() started")
	for r := 0; r < 10; r++ {
		positive := make([]int, math.MaxInt16+1)
		negative := make([]int, math.MaxInt16+1)

		// no-hit 확률은 0.1%
		for count := 0; count < math.MaxUint16*1000; count++ {
			// next int32
			r := Next[int32]()
			if r >= 0 {
				positive[int(r)%len(positive)]++
			} else {
				negative[uint32(r)%uint32(len(negative))]++
			}
		}

		// hit test
		for i, c := range positive {
			if c == 0 {
				t.Errorf("NextInt32(): missing hit value for positive %d", i)
			}
		}
		for i, c := range negative {
			if c == 0 {
				t.Errorf("NextInt32(): missing hit value for negative %d", i)
			}
		}
	}
	t.Logf("TestNextInt32() completed")
}

func TestNextUint32(t *testing.T) {
	t.Logf("TestNextUint32() started")
	for r := 0; r < 10; r++ {
		positive := make([]int, math.MaxInt16+1)

		// no-hit 확률은 0.1%
		for count := 0; count < math.MaxUint16*1000; count++ {
			// next uint32
			r := Next[uint32]()
			positive[r%uint32(len(positive))]++
		}

		// hit test
		for i, c := range positive {
			if c == 0 {
				t.Errorf("NextUint32(): missing hit value for positive %d", i)
			}
		}
	}
	t.Logf("TestNextUint32() completed")
}

func TestNextInt64(t *testing.T) {
	t.Logf("TestNextInt64() started")
	for r := 0; r < 10; r++ {
		positive := make([]int, math.MaxInt16+1)
		negative := make([]int, math.MaxInt16+1)

		// no-hit 확률은 0.1%
		for count := 0; count < math.MaxUint16*1000; count++ {
			// next int64
			r := Next[int64]()
			if r >= 0 {
				positive[r%int64(len(positive))]++
			} else {
				negative[uint64(r)%uint64(len(negative))]++
			}
		}

		// hit test
		for i, c := range positive {
			if c == 0 {
				t.Errorf("NextInt64(): missing hit value for positive %d", i)
			}
		}
		for i, c := range negative {
			if c == 0 {
				t.Errorf("NextInt64(): missing hit value for negative %d", i)
			}
		}
	}
	t.Logf("TestNextInt64() completed")
}

func TestNextUint64(t *testing.T) {
	t.Logf("TestNextUint64() started")
	for r := 0; r < 10; r++ {
		positive := make([]int, math.MaxInt16+1)

		// no-hit 확률은 0.1%
		for count := 0; count < math.MaxUint16*1000; count++ {
			// next uint64
			r := Next[uint64]()
			positive[r%uint64(len(positive))]++
		}

		// hit test
		for i, c := range positive {
			if c == 0 {
				t.Errorf("NextUint64(): missing hit value for positive %d", i)
			}
		}
	}
	t.Logf("TestNextUint64() completed")
}

func TestNextInt(t *testing.T) {
	t.Logf("TestNextInt() started")
	for r := 0; r < 10; r++ {
		positive := make([]int, math.MaxUint16)
		negative := make([]int, math.MaxUint16)

		// no-hit 확률은 0.1%
		for count := 0; count < math.MaxUint16*1000; count++ {
			// next int
			r := Next[int]()
			if r >= 0 {
				positive[r%math.MaxUint16]++
			} else {
				r = -r
				negative[r%math.MaxUint16]++
			}
		}

		// hit test
		for i, c := range positive {
			if c == 0 {
				t.Errorf("NextInt(): missing hit value for positive %d", i)
			}
		}
		for i, c := range negative {
			if c == 0 {
				t.Errorf("NextInt(): missing hit value for negative %d", i)
			}
		}
	}
	t.Logf("TestNextInt() completed")
}

func TestNextUint(t *testing.T) {
	t.Logf("TestNextUint() started")
	for r := 0; r < 10; r++ {
		positive := make([]int, math.MaxUint16)

		// no-hit 확률은 0.1%
		for count := 0; count < math.MaxUint16*1000; count++ {
			// next int
			r := Next[uint]()
			positive[r%math.MaxUint16]++
		}

		// hit test
		for i, c := range positive {
			if c == 0 {
				t.Errorf("NextInt(): missing hit value for positive %d", i)
			}
		}
	}
	t.Logf("TestNextUint() completed")
}

func TestNextInRange(t *testing.T) {

	repeats := 1000000
	t.Logf("TestNextInRange() started. repeats=%d", repeats)

	// byte
	for x := math.MinInt8; x <= math.MaxInt8; x++ {
		for y := math.MinInt8; y <= math.MaxInt8; y++ {
			min, max := byte(x), byte(y)

			// should panic if min > max
			if min > max {
				assert.Panics(t, func() { NextInRange(min, max) })
				continue
			}

			actual := NextInRange(min, max)

			// should return min if min == max
			if min == max {
				assert.Equal(t, min, actual)
				continue
			}

			// should be min <= actual < max
			assert.Truef(t, actual >= min && actual < max,
				"NextInRange[%s](): actual(%d) is not in range(%d, %d)",
				reflect.TypeOf(actual).Name(), actual, min, max)
		}
	}
	t.Logf("byte completed")

	// int
	for x := 0; x <= repeats; x++ {
		min, max := Next[int](), Next[int]()

		// should panic if min > max
		if min > max {
			assert.Panics(t, func() { NextInRange(min, max) })
			continue
		}

		actual := NextInRange(min, max)

		// should return min if min == max
		if min == max {
			assert.Equal(t, min, actual)
			continue
		}

		// should be min <= actual < max
		assert.Truef(t, actual >= min && actual < max,
			"NextInRange[%s](): actual(%d) is not in range(%d, %d)",
			reflect.TypeOf(actual).Name(), actual, min, max)
	}
	t.Logf("int completed")

	// uint
	for x := 0; x <= repeats; x++ {
		min, max := Next[uint](), Next[uint]()

		// should panic if min > max
		if min > max {
			assert.Panics(t, func() { NextInRange(min, max) })
			continue
		}

		actual := NextInRange(min, max)

		// should return min if min == max
		if min == max {
			assert.Equal(t, min, actual)
			continue
		}

		// should be min <= actual < max
		assert.Truef(t, actual >= min && actual < max,
			"NextInRange[%s](): actual(%d) is not in range(%d, %d)",
			reflect.TypeOf(actual).Name(), actual, min, max)
	}
	t.Logf("uint completed")

	// int8
	for x := math.MinInt8; x <= math.MaxInt8; x++ {
		for y := math.MinInt8; y <= math.MaxInt8; y++ {
			min, max := int8(x), int8(y)

			// should panic if min > max
			if min > max {
				assert.Panics(t, func() { NextInRange(min, max) })
				continue
			}

			actual := NextInRange(min, max)

			// should return min if min == max
			if min == max {
				assert.Equal(t, min, actual)
				continue
			}

			// should be min <= actual < max
			assert.Truef(t, actual >= min && actual < max,
				"NextInRange[%s](): actual(%d) is not in range(%d, %d)",
				reflect.TypeOf(actual).Name(), actual, min, max)
		}
	}
	t.Logf("int8 completed")

	// uint8
	for x := 0; x <= math.MaxUint8; x++ {
		for y := 0; y <= math.MaxUint8; y++ {
			min, max := uint8(x), uint8(y)

			// should panic if min > max
			if min > max {
				assert.Panics(t, func() { NextInRange(min, max) })
				continue
			}

			actual := NextInRange(min, max)

			// should return min if min == max
			if min == max {
				assert.Equal(t, min, actual)
				continue
			}

			// should be min <= actual < max
			assert.Truef(t, actual >= min && actual < max,
				"NextInRange[%s](): actual(%d) is not in range(%d, %d)",
				reflect.TypeOf(actual).Name(), actual, min, max)
		}
	}
	t.Logf("uint8 completed")

	// int16
	for x := 0; x <= repeats; x++ {
		min, max := Next[int16](), Next[int16]()

		// should panic if min > max
		if min > max {
			assert.Panics(t, func() { NextInRange(min, max) })
			continue
		}

		actual := NextInRange(min, max)

		// should return min if min == max
		if min == max {
			assert.Equal(t, min, actual)
			continue
		}

		// should be min <= actual < max
		assert.Truef(t, actual >= min && actual < max,
			"NextInRange[%s](): actual(%d) is not in range(%d, %d)",
			reflect.TypeOf(actual).Name(), actual, min, max)
	}
	t.Logf("int16 completed")

	// uint16
	for x := 0; x <= repeats; x++ {
		min, max := Next[uint16](), Next[uint16]()

		// should panic if min > max
		if min > max {
			assert.Panics(t, func() { NextInRange(min, max) })
			continue
		}

		actual := NextInRange(min, max)

		// should return min if min == max
		if min == max {
			assert.Equal(t, min, actual)
			continue
		}

		// should be min <= actual < max
		assert.Truef(t, actual >= min && actual < max,
			"NextInRange[%s](): actual(%d) is not in range(%d, %d)",
			reflect.TypeOf(actual).Name(), actual, min, max)
	}
	t.Logf("uint16 completed")

	// int32
	for x := 0; x <= repeats; x++ {
		min, max := Next[int32](), Next[int32]()

		// should panic if min > max
		if min > max {
			assert.Panics(t, func() { NextInRange(min, max) })
			continue
		}

		actual := NextInRange(min, max)

		// should return min if min == max
		if min == max {
			assert.Equal(t, min, actual)
			continue
		}

		// should be min <= actual < max
		assert.Truef(t, actual >= min && actual < max,
			"NextInRange[%s](): actual(%d) is not in range(%d, %d)",
			reflect.TypeOf(actual).Name(), actual, min, max)
	}
	t.Logf("int32 completed")

	// uint32
	for x := 0; x <= repeats; x++ {
		min, max := Next[uint32](), Next[uint32]()

		// should panic if min > max
		if min > max {
			assert.Panics(t, func() { NextInRange(min, max) })
			continue
		}

		actual := NextInRange(min, max)

		// should return min if min == max
		if min == max {
			assert.Equal(t, min, actual)
			continue
		}

		// should be min <= actual < max
		assert.Truef(t, actual >= min && actual < max,
			"NextInRange[%s](): actual(%d) is not in range(%d, %d)",
			reflect.TypeOf(actual).Name(), actual, min, max)
	}
	t.Logf("uint32 completed")

	// int64
	for x := 0; x <= repeats; x++ {
		min, max := Next[int64](), Next[int64]()

		// should panic if min > max
		if min > max {
			assert.Panics(t, func() { NextInRange(min, max) })
			continue
		}

		actual := NextInRange(min, max)

		// should return min if min == max
		if min == max {
			assert.Equal(t, min, actual)
			continue
		}

		// should be min <= actual < max
		assert.Truef(t, actual >= min && actual < max,
			"NextInRange[%s](): actual(%d) is not in range(%d, %d)",
			reflect.TypeOf(actual).Name(), actual, min, max)
	}
	t.Logf("int64 completed")

	// uint64
	for x := 0; x <= repeats; x++ {
		min, max := Next[uint64](), Next[uint64]()

		// should panic if min > max
		if min > max {
			assert.Panics(t, func() { NextInRange(min, max) })
			continue
		}

		actual := NextInRange(min, max)

		// should return min if min == max
		if min == max {
			assert.Equal(t, min, actual)
			continue
		}

		// should be min <= actual < max
		assert.Truef(t, actual >= min && actual < max,
			"NextInRange[%s](): actual(%d) is not in range(%d, %d)",
			reflect.TypeOf(actual).Name(), actual, min, max)
	}
	t.Logf("uint64 completed")
	t.Logf("TestNextInRange() completed")
}

func TestNextMyNumberType(t *testing.T) {
	t.Logf("TestNextMyNumberType() started")

	repeats := 1000000

	// type byte = uint8
	for i := 0; i < repeats; i++ {
		actual := Next[byte]()
		assert.Truef(t, 0 <= actual && math.MaxUint8 >= actual, "Next[byte](): actual(%d) not in range(0, %d)", actual, math.MaxUint8)
	}

	// int/uint
	type MyInt = int
	for i := 0; i < repeats; i++ {
		actual := Next[MyInt]()
		assert.Truef(t, math.MinInt <= actual && math.MaxInt >= actual, "Next[MyInt](): actual(%d) not in range(%d, %d)", actual, math.MinInt, math.MaxInt)
	}
	type MyUint = uint
	for i := 0; i < repeats; i++ {
		actual := Next[MyUint]()
		assert.Truef(t, 0 <= actual && math.MaxUint >= actual, "Next[MyUint](): actual(%d) not in range(0, %d)", actual, uint64(math.MaxUint))
	}

	// int8/uint8
	type MyInt8 = int8
	for i := 0; i < repeats; i++ {
		actual := Next[MyInt8]()
		assert.Truef(t, math.MinInt8 <= actual && math.MaxInt8 >= actual, "Next[MyInt8](): actual(%d) not in range(%d, %d)", actual, math.MinInt8, math.MaxInt8)
	}
	type MyUint8 = uint8
	for i := 0; i < repeats; i++ {
		actual := Next[MyUint8]()
		assert.Truef(t, 0 <= actual && math.MaxUint8 >= actual, "Next[MyUint8](): actual(%d) not in range(0, %d)", actual, math.MaxUint8)
	}

	// int16/uint16
	type MyInt16 = int16
	for i := 0; i < repeats; i++ {
		actual := Next[MyInt16]()
		assert.Truef(t, math.MinInt16 <= actual && math.MaxInt16 >= actual, "Next[MyInt16](): actual(%d) not in range(%d, %d)", actual, math.MinInt16, math.MaxInt16)
	}
	type MyUint16 = uint16
	for i := 0; i < repeats; i++ {
		actual := Next[MyUint16]()
		assert.Truef(t, 0 <= actual && math.MaxUint16 >= actual, "Next[MyUint16](): actual(%d) not in range(0, %d)", actual, math.MaxUint16)
	}

	// int32/uint32
	type MyInt32 = int32
	for i := 0; i < repeats; i++ {
		actual := Next[MyInt32]()
		assert.Truef(t, math.MinInt32 <= actual && math.MaxInt32 >= actual, "Next[MyInt32](): actual(%d) not in range(%d, %d)", actual, math.MinInt32, math.MaxInt32)
	}
	type MyUint32 = uint32
	for i := 0; i < repeats; i++ {
		actual := Next[MyUint32]()
		assert.Truef(t, 0 <= actual && math.MaxUint32 >= actual, "Next[MyUint32](): actual(%d) not in range(0, %d)", actual, math.MaxUint32)
	}

	// int64/uint64
	type MyInt64 = int64
	for i := 0; i < repeats; i++ {
		actual := Next[MyInt64]()
		assert.Truef(t, math.MinInt64 <= actual && math.MaxInt64 >= actual, "Next[MyInt64](): actual(%d) not in range(%d, %d)", actual, math.MinInt64, math.MaxInt64)
	}
	type MyUint64 = uint64
	for i := 0; i < repeats; i++ {
		actual := Next[MyUint64]()
		assert.Truef(t, 0 <= actual && math.MaxUint64 >= actual, "Next[MyUint64](): actual(%d) not in range(0, %d)", actual, uint64(math.MaxUint64))
	}
	t.Logf("TestNextMyNumberType() completed")
}
