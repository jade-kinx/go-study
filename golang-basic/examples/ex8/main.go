package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	ch := make(chan string, 10000)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT)

	for i := 0; i < 16; i++ {
		go receiver(i, ch)
	}

	for i, quit := 0, false; !quit; i++ {
		select {
		case <-done:
			quit = true
		default:
		}

		ch <- fmt.Sprintf("ping(%d)", i)
	}

	fmt.Println("program terminated")
}

func receiver(i int, ch <-chan string) {
	for data := range ch {
		fmt.Printf("[%02d] received: %s\n", i, data)
	}
}
