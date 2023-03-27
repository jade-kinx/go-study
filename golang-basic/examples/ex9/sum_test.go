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
