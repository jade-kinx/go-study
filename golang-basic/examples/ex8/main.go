package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func ChannelBasicUsage() {
	channel := make(chan string) // 채널

	// 채널 송신자 N
	wgs := sync.WaitGroup{}
	for i := 0; i < runtime.NumCPU(); i++ {
		wgs.Add(1)
		go func(ch chan<- string, sender int) { // write-only 채널
			defer fmt.Printf("sender[%d] terminated\n", sender)
			defer wgs.Done()
			for i := 0; i < 10; i++ {
				time.Sleep(time.Second * 1)
				// 채널에 메세지 입력
				ch <- fmt.Sprintf("sender[%d]: Hello, Go! (%d)", sender, i)
			}
		}(channel, i)
	}

	// 채널 수신자 1
	wgr := sync.WaitGroup{}
	wgr.Add(1)
	go func(ch <-chan string) { // read-only 채널
		defer fmt.Println("receiver terminated")
		defer wgr.Done()
		// 채널에서 메세지를 수신 ( 채널이 닫히면 for loop 종료)
		for msg := range ch {
			fmt.Println(msg)
		}
	}(channel)

	// wait for channel senders
	wgs.Wait()

	// close channel
	close(channel)

	// wait for channel receiver
	wgr.Wait()
}

func main() {

	ChannelBasicUsage()
}
