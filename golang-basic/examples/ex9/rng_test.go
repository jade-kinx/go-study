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

	// all covered?
	for i, c := range append(positives, negatives...) {
		assert.NotZerof(c, "uncovered value: %d", i)
	}
}

func TestNextUint8ShouldCoverAll(t *testing.T) {
	assert := assert.New(t)

	repeat := 100
	covered := make([]int, math.MaxUint8+1)
	for i := 0; i < math.MaxUint8*repeat; i++ {
		b := Next[uint8]()
		covered[int(b)]++ // increase hit count
	}

	// all covered?
	for i, c := range covered {
		assert.NotZerof(c, "uncovered value: %d", i)
	}
}

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
