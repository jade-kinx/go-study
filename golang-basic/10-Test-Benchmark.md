# 테스트/벤치마크

## TDD (테스트 주도 개발)
테스트 주도 개발(TDD)은 소프트웨어를 개발하는 여러 방법론 중 하나이다. 제품이 오류 없이 정상 작동하는지 확인하기 위해 모든 코드는 프로그래머가 작성하고 나서 테스트를 거치게 된다. TDD에서는 제품의 기능 구현을 위한 코드와 별개로, 해당 기능이 정상적으로 움직이는지 검증하기 위한 테스트 코드를 작성한다. 이를 통해 테스트가 실패할 경우, 테스트를 통과하기 위한 최소한으로 코드를 개선한다. 최종적으로 테스트에 성공한 코드를 리팩토링 하는 과정을 거친다.

### TDD 효과
* 코드가 내 손을 벗어나기 전에 가장 빠르게 피드백을 받을 수 있다.  
* 작성한 코드가 가지는 불안정성을 개선하여 신뢰성/생산성/품질을 높일 수 있다.  
* 개발자의 심리적 안정과 성취를 보장한다.  
* 잘 짜여진 테스트 코드는 그 자체로 코드에 대한 문서로 대체할 수 있다.  
* 복잡한 프로젝트에서 코드의 변경이 다른 모듈에 미치는 영향을 빠르게 감지할 수 있다.  
* 일부, 개발 기간이 길어지고 생산성이 떨어진다는 의견도 있지만, 디버깅 포함 프로젝트 완성까지 고려하면 오히려 기간이 단축되는 경우가 많다.  
  * 테스트가 필요한 단위를 잘 선정하여 과도한 테스트를 방지해야 한다.  
  * 보통, 내가 잘 아는 부분이 아닌 새로 접하게 되는 영역에서 확신을 가지는 차원에서 테스트를 추가  

### 단위 테스트와 통합 테스트

### 좋은 테스트의 특징(FIRST)
* Fast: 테스트는 빠르게 동작하여 자주 돌릴 수 있어야 한다.  
* Independent: 각각의 테스트는 독립적이며 서로 의존해서는 안된다.  
* Repeatable: 어느 환경에서도 반복 가능해야 한다.  
* Self-Validating: 테스트는 성공 또는 실패로 결과를 내어 자체적으로 검증할 수 있어야 한다.  
* Timely: 테스트는 적시에 즉, 테스트하려는 실제 코드를 구현하기 직전에 구현해야 한다.  

## 테스트/벤치마크 과정 예제

* 테스트 파일은 `*_test.go` 이어야 한다.  
* 테스트 파일은 패키지 빌드시 제외된다.  
* 테스트 함수 이름은 `TestXXX` 여야 한다.
* 벤치마크 함수 이름은 `BenchmarkXXX` 여야 한다.  

```bash
# 테스트 실행
$ go test -v -run ^TestSumToNSanityCheck$
```

### 0. 예제 작업 명세

> 임의의 정수 `n`이 주어지면 1부터 `n`까지의 합을 반환한다.  
> `n`이 음의 정수이면 `error`를 반환한다.

### 1. Stub 작성
```go
// SumToN은 1부터 n까지 정수를 더한 값을 반환한다.
// n < 0 이면 에러를 반환한다.
func SumToN(n int) (int, error) {
	return 0, fmt.Errorf("not implemented")
}
```

### 2. 테스트 작성
```go
func TestSumToNSanityCheck(t *testing.T) {
	// https://www.mathsisfun.com/numbers/sigma-calculator.html
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
```

```bash
# 테스트 실행 결과(실패)
Running tool: C:\Program Files\Go\bin\go.exe test -timeout 30m -run ^TestSumToNSanityCheck$ ex9 -v

=== RUN   TestSumToNSanityCheck
    d:\gitworks\go-study\golang-basic\examples\ex9\ex9_test.go:73: SumToN(0) returns error: not implemented
--- FAIL: TestSumToNSanityCheck (0.00s)
FAIL
FAIL	ex9	0.030s
```


### 3. working 코드 작성
```go
// SumToN은 1부터 n까지 정수를 더한 값을 반환한다.
// n < 0 이면 에러를 반환한다.
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
```
```bash
### 테스트 실행 결과(성공)
Running tool: C:\Program Files\Go\bin\go.exe test -timeout 30m -run ^TestSumToNSanityCheck$ ex9 -v

=== RUN   TestSumToNSanityCheck
--- PASS: TestSumToNSanityCheck (0.00s)
PASS
ok  	ex9	0.031s
```

### 4. 벤치마크 코드 작성
```go
func BenchmarkSumToN(b *testing.B) {
	for n := 0; n < b.N; n++ {
		SumToN(n)	// 원래는 b.N이 테스트에 관여하면 안됨
	}
}
```
```bash
# 벤치마크 결과
Running tool: C:\Program Files\Go\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkSumToN$ ex9 -v

goos: windows
goarch: amd64
pkg: ex9
cpu: 12th Gen Intel(R) Core(TM) i7-12700
BenchmarkSumToN
BenchmarkSumToN-20
 1000000	    105678 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	ex9	105.730s
```

### 5. 리팩토링 코드
```go
// SumToN은 1부터 n까지 정수를 더한 값을 반환한다.
// n < 0 이면 에러를 반환한다.
func SumToN(n int) (int, error) {
	if n < 0 {
		return 0, fmt.Errorf("n(%d) should be greater or equal than 0", n)
	}
	// f(n) = n * (n + 1) / 2
	return n * (n + 1) / 2, nil
}
```
코드 리팩토링 후 테스트를 통과하는지 확인한 후에 다시 벤치마크를 수행한다.

```bash
# 리팩토링 후 벤치마크 결과
Running tool: C:\Program Files\Go\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkSumToN$ ex9 -v

goos: windows
goarch: amd64
pkg: ex9
cpu: 12th Gen Intel(R) Core(TM) i7-12700
BenchmarkSumToN
BenchmarkSumToN-20
1000000000	         0.1078 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	ex9	0.153s
```

### 6. 이걸로 끝???
* `int` 범위를 초과하면?
* `32bit` or `64bit` ??
* 그 외 모든 경우의 수를 더 고려하고 테스트 케이스를 작성하고 리팩토링을 반복  

```go
// SumToN 함수의 int overflow를 테스트한다.
func TestSumToNShouldReturnOverflowError(t *testing.T) {
	n := math.MaxInt // 9223372036854775807
	_, err := SumToN(n)
	if err == nil {
		t.Fatalf("SumToN(%d) should return overflow error", n)
	}
}
```

```bash
# 테스트 실패(int 내부에서 overflow발생하지만 error를 반환하지는 않음)
Running tool: C:\Program Files\Go\bin\go.exe test -timeout 30m -run ^TestSumToNShouldReturnErrorOnOverflow$ ex9 -v

=== RUN   TestSumToNShouldReturnErrorOnOverflow
    d:\gitworks\go-study\golang-basic\examples\ex9\ex9_test.go:70: SumToN(9223372036854775807) should return overflow error
--- FAIL: TestSumToNShouldReturnErrorOnOverflow (0.00s)
FAIL
FAIL	ex9	0.033s
```

다시 리팩토링 (`math/big.Int` 사용)
```go
// SumToN은 1부터 n까지 정수를 더한 값을 반환한다.
// n < 0 이면 에러를 반환한다.
// 결과값이 int 범위를 벗어나면 overflow error를 반환한다.
func SumToN(n int) (int, error) {
	if n < 0 {
		return 0, fmt.Errorf("n(%d) should be greater or equal than 0", n)
	}

	// f(n) = n * (n + 1) / 2
	v := big.NewInt(0)
	v = v.Mul(big.NewInt(int64(n)), big.NewInt(int64(n)+1))
	v = v.Div(v, big.NewInt(2))

	// overflow?
	if !v.IsInt64() {
		return 0, fmt.Errorf("sum overflow int64")
	}
	sum := v.Int64()

	// check overflow for int32 in 32bit system
	if strconv.IntSize == 32 {
		if sum < math.MinInt32 || sum > math.MaxInt32 {
			return 0, fmt.Errorf("sum overflow int32")
		}
	}

	return int(sum), nil
}
```
```bash
# 벤치마크 결과
Running tool: C:\Program Files\Go\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkSumToN$ ex9 -v

goos: windows
goarch: amd64
pkg: ex9
cpu: 12th Gen Intel(R) Core(TM) i7-12700
BenchmarkSumToN
BenchmarkSumToN-20
30295455	        35.17 ns/op	      48 B/op	       1 allocs/op
PASS
ok  	ex9	1.139s
```

### 그 외
* 가급적이면 입력 범위의 모든 경우를 테스트 케이스로 작성하는 것이 좋지만, 시간이 오래 걸리므로 적절히 트레이드 오프가 필요(주로 엣지 테스트)  
* 내 코드에 대한 신뢰성 확보가 문제 발생시 동료의 코드에 대한 의심으로 번지면 안됨!  
