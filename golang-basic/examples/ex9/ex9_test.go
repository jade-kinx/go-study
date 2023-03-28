package main

import (
	"fmt"
	"math"

	// "math/big"
	// "strconv"
	"testing"
)

// sum to n stub code
// func SumToN(n int) (int, error) {
// 	return 0, fmt.Errorf("not implemented")
// }

// SumToN은 0부터 n까지 정수를 더한 값을 반환
func SumToN(n int) (int, error) {
	if n < 0 {
		return 0, fmt.Errorf("n(%d) should be greater or equal than 0", n)
	}

	sum := 0
	for i := 1; i <= n; i++ {
		sum += i
	}
	return sum, nil
}

// SumToN은 0부터 n까지 정수를 더한 값을 반환
// func SumToN(n int) (int, error) {
// 	if n < 0 {
// 		return 0, fmt.Errorf("n(%d) should be greater or equal than 0", n)
// 	}

// 	return n * (n + 1) / 2, nil
// }

// SumToN은 0부터 n까지 정수를 더한 값을 에러와 함께 반환
// func SumToN(n int) (int, error) {
// 	if n < 0 {
// 		return 0, fmt.Errorf("n(%d) should be greater or equal than 0", n)
// 	}

// 	// f(n) = n * (n + 1) / 2
// 	v := big.NewInt(0)
// 	v = v.Mul(big.NewInt(int64(n)), big.NewInt(int64(n+1)))
// 	v = v.Div(v, big.NewInt(2))

// 	// sum overflow?
// 	if !v.IsInt64() {
// 		return 0, fmt.Errorf("sum overflow int64")
// 	}

// 	sum := v.Int64()

// 	// check overflow for int32 in 32bit system
// 	if strconv.IntSize == 32 {
// 		if sum < math.MinInt32 || sum > math.MaxInt32 {
// 			return 0, fmt.Errorf("sum overflow int32")
// 		}
// 	}

// 	return int(sum), nil
// }

// SumToN 함수의 기본 동작을 테스트한다.
func TestSumToNSanityCheck(t *testing.T) {
	// expected: https://www.mathsisfun.com/numbers/sigma-calculator.html
	expected := []int{0, 1, 3, 6, 10, 15, 21, 28, 36, 45, 55}
	for n, want := range expected {
		have, err := SumToN(n)
		if err != nil {
			t.Fatalf("SumToN(%d) returns error: %s", n, err)
		}

		if want != have {
			t.Fatalf("SumToN(%d) unexpected! want: %d, have: %d", n, want, have)
		}
	}

	// should return error if n < 0
	_, err := SumToN(-1)
	if err == nil {
		t.Fatalf("SumToN(-1) should return error")
	}
}

// SumToN 함수의 성능을 측정한다.
func BenchmarkSumToN(b *testing.B) {
	for n := 0; n < b.N; n++ {
		SumToN(n)
	}
}

// SumToN 함수의 int overflow를 테스트한다.
func TestSumToNShouldReturnOverflowError(t *testing.T) {
	n := math.MaxInt // 9223372036854775807
	_, err := SumToN(n)
	if err == nil {
		t.Fatalf("SumToN(%d) should return overflow error", n)
	}
}
