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
