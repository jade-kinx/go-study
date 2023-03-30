package main

import (
	"fmt"
)

func test() (int, error) {
	return 0, nil
}

func main() {
	fmt.Println("Hello, Go!")

	i := 100
	fmt.Printf("i=%d(%p)\n", i, &i)

	i, str := 200, "hello"
	fmt.Printf("i=%d(%p)\n", i, &i)
	fmt.Printf("str=%s\n", str)

	if i, err := test(); err == nil {
		fmt.Printf("i=%d(%p)\n", i, &i)
	}

	fmt.Printf("i=%d(%p)\n", i, &i)
}
