package main

import (
	"fmt"
	"github.com/jade-kinx/go-study/golang-basic/examples/common"
)

func TypeDefExamples() {
	var a int = 3 // 1. 변수 선언 기본 형태
	var b int     // 2. 변수 선언 초기값 생략 (타입의 기본값(zero-value): 0, 0.0, false, "", nil)
	var c = a     // 3. 타입 생략 (c의 타입은 우변의 타입)
	d := 3        // 4. 선언 대입문(d의 타입은 3이 int 리터럴이므로 int 타입)
	fmt.Printf("a=%d, b=%d, c=%d, d=%d(%p)\n", a, b, c, d, &d)
	/* OUTPUT
	a=3, b=0, c=3, d=3(0xc0000a6058)
	*/

	// var d                // X : d의 타입을 알 수 없음
	// d := 100             // X : d는 이미 선언되어 있으므로 선언 대입할 수 없음
	d, e := 100, "hello" // d는 기존 선언된 변수의 공간을 재사용, e는 선언 대입
	fmt.Printf("d=%d(%p), e=%s\n", d, &d, e)
	/* OUTPUT
	d=100(0xc0000a6058), e=hello
	*/
}

func PrintFormatExamples() {
	var i int64 = 3456
	fmt.Printf("binary format: %08b\n", i)

	bin, err := common.ToBinaryString(i)
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}
	fmt.Printf("binary format: %s\n", bin)

	x := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	bs, err := common.ToBinaryString(x)
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}
	fmt.Printf("binary format: %s\n", bs)
}

func main() {
	fmt.Println("Hello, Go!")

	TypeDefExamples()

	PrintFormatExamples()
}
