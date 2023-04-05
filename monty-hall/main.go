package main

import (
	"math/rand"
	"time"

	"golang.org/x/text/message"
)

const (
	DoorCount = 3
	repeats   = 100000000
)

// 몬티홀 문제에 대해 repeats횟수만큼 시행하여,
// 선택을 바꿨을 때와 바꾸지 않았을 때의 당첨 확률을 구해서 출력한다.
func main() {

	p := message.NewPrinter(message.MatchLanguage("ko"))
	p.Println("Hello, Monty-Hall Problem!")
	defer p.Println("Bye, Monty-Hall Problem!")

	winsWhenUnchanged := 0
	winsWhenChanged := 0
	start := time.Now()

	// repeats 횟수만큼 반복 시행
	for i := 0; i < repeats; i++ {
		// 0. 새로운 문제 생성
		mh := NewMontyHall()

		// 1. 플레이어가 임의의 문을 선택한다.
		choice := selectDoorRandomly()
		mh.Choice(choice)

		// 2. 사회자가 미리보기 문을 열어준다.
		preview := mh.Preview()

		// 3. 처음 선택을 고수했을 때 당첨 여부 확인
		if mh.IsWinPrize(choice) {
			winsWhenUnchanged++
		}

		// 4. 선택을 바꿨을 때 당첨 여부 확인
		changed := mh.ChangeFrom(preview)
		if mh.IsWinPrize(changed) {
			winsWhenChanged++
		}
	}

	// 선택을 바꾸지 않은 경우와 바꾼 경우의 당첨 확률 결과 출력
	p.Printf("unchanged: %d/%d(%.2f%%)\n", winsWhenUnchanged, repeats, float64(winsWhenUnchanged)*100/repeats)
	p.Printf("changed: %d/%d(%.2f%%)\n", winsWhenChanged, repeats, float64(winsWhenChanged)*100/repeats)
	p.Printf("repeats: %d, elapsed: %.2f(s)\n", repeats, time.Since(start).Seconds())
}

/* OUTPUT
Hello, Monty-Hall Problem!
unchanged: 33,330,098/100,000,000(33.33%)
changed: 66,669,902/100,000,000(66.67%)
repeats: 100,000,000, elapsed: 7.18(s)
Bye, Monty-Hall Problem!
*/

// 0, 1, 2 문 중 하나의 번호를 랜덤하게 선택한다
func selectDoorRandomly() int {
	return rand.Intn(DoorCount)
}

type MontyHall struct {
	doors  [DoorCount]bool // 상품이 포함되어 있는 문들 ( true: 스포츠카(당첨), false: 염소(꽝) )
	choice int             // 플레이어가 최초로 선택한 문
}

// 새로운 몬티홀 문제를 생성한다.
func NewMontyHall() *MontyHall {
	mh := MontyHall{choice: -1}
	mh.doors[selectDoorRandomly()] = true // 0, 1, 2 문중 하나에 당첨을 설정
	return &mh
}

// 플레이어가 문을 선택한다.
func (mh *MontyHall) Choice(n int) {
	mh.choice = n
}

// 사회자가 플레이어가 선택한 문 또는 당첨이 아닌 문을 연다.
func (mh *MontyHall) Preview() int {
	n := selectDoorRandomly()
	for mh.doors[n] || n == mh.choice {
		n = selectDoorRandomly()
	}
	return n
}

// 당첨인지 확인한다.
func (mh *MontyHall) IsWinPrize(n int) bool {
	return mh.doors[n]
}

// 선택을 바꾼다.
func (mh *MontyHall) ChangeFrom(preview int) int {
	for i := 0; i < DoorCount; i++ {
		if i != mh.choice && i != preview {
			return i
		}
	}
	panic("wtf! this should not happen!")
}
