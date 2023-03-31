package main

import (
	"fmt"
	"runtime"
	"sync"
)

// N senders -> 1 receiver
func ChannelForSendersToReceiver() {
	channel := make(chan string) // 채널
	wgs := sync.WaitGroup{}      // 송신자 wg
	wgr := sync.WaitGroup{}      // 수신자 wg

	// 채널 송신자 N
	for i := 0; i < runtime.NumCPU(); i++ {
		wgs.Add(1)
		go func(ch chan<- string, id int) { // write-only 채널
			defer fmt.Printf("sender[%d] terminated\n", id)
			defer wgs.Done()
			for counter := 0; counter < 10; counter++ {
				// 채널에 메세지 입력
				ch <- fmt.Sprintf("sender[%d]: Hello, Go! (%d)", id, counter)
			}
		}(channel, i)
	}

	// 채널 수신자 1
	wgr.Add(1)
	go func(ch <-chan string) { // read-only 채널
		defer fmt.Println("receiver terminated")
		defer wgr.Done()
		// 채널에서 메세지를 수신 (채널이 닫히면 for loop 종료)
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

// 1 sender -> N receivers
func ChannelForSenderToReceivers() {
	channel := make(chan string)
	wgs := sync.WaitGroup{} // 송신자 wg
	wgr := sync.WaitGroup{} // 수신자 wg

	// 채널 송신자 1
	wgs.Add(1)
	go func(ch chan<- string) { // write-only 채널
		defer fmt.Println("sender terminated")
		defer wgs.Done()
		for counter := 0; counter < 10; counter++ {
			// 채널에 메세지 입력
			ch <- fmt.Sprintf("Hello, Go! (%d)", counter)
		}
	}(channel)

	// 채널 수신자 N (어떤 수신자가 메세지를 받을지는 고런타임 스케쥴러가 결정)
	for i := 0; i < runtime.NumCPU(); i++ {
		wgr.Add(1)
		go func(ch <-chan string, id int) { // read-only 채널
			defer fmt.Printf("receiver[%d] terminated\n", id)
			defer wgr.Done()
			// 채널에서 메세지를 수신 (채널이 닫히면 for loop 종료)
			for msg := range ch {
				fmt.Printf("receiver[%d]: %s\n", id, msg)
			}
		}(channel, i)
	}

	// wait for channel senders
	wgs.Wait()

	// close channel
	close(channel)

	// wait for channel receiver
	wgr.Wait()
}

func main() {

	fmt.Println("ChannelForSendersToReceiver() started")
	ChannelForSendersToReceiver()
	fmt.Println("ChannelForSendersToReceiver() completed")

	fmt.Println("ChannelForSenderToReceivers() started")
	ChannelForSenderToReceivers()
	fmt.Println("ChannelForSenderToReceivers() completed")
}
