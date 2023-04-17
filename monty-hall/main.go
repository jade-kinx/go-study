package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

const (
	DoorCount = 3
)

// 플레이어 인터페이스
type Player interface {
	PickDoor() int            // 플레이어가 3개의 문 중 하나의 문을 선택
	SwitchChoice() bool       // 사회자가 바꾸겠냐고 물어봤을 때의 대답 (true: 변경, false: 유지)
	ProcessWinOrNot(win bool) // 당첨 여부 확인
	Wins() int                // 당첨 횟수
}

// 플레이어
type player struct {
	wins int // 당첨 횟수
}

func (p player) PickDoor() int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("pick door(1/2/3) ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return p.PickDoor()
	}

	choice, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return p.PickDoor()
	}

	if choice < 1 || choice > DoorCount {
		fmt.Printf("door number(%d) out of range", choice)
		return p.PickDoor()
	}

	return choice
}

func (p player) SwitchChoice() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("change your mind? (y/n) ")
	answer, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return false
	}

	switch strings.ToLower(strings.TrimSpace(answer)) {
	case "y", "yes":
		fmt.Println("changing your mind")
		return true
	default:
		fmt.Println("staying your mind")
		return false
	}
}

func (p player) Wins() int {
	return p.wins
}

func (p *player) ProcessWinOrNot(win bool) {
	if win {
		p.wins++
		fmt.Println("congraturations! you win the prize!")
	} else {
		fmt.Println("you lost the game!")
	}
}

// 스테이 플레이어
type player1 player

func (p player1) PickDoor() int {
	return pickDoor()
}

// 항상 선택을 유지
func (p player1) SwitchChoice() bool {
	return false
}

func (p player1) Wins() int {
	return p.wins
}

func (p *player1) ProcessWinOrNot(win bool) {
	if win {
		p.wins++
	}
}

// 체인지 플레이어
type player2 player

func (p player2) PickDoor() int {
	return pickDoor()
}

// 항상 선택을 바꾼다
func (p player2) SwitchChoice() bool {
	return true
}

func (p player2) Wins() int {
	return p.wins
}

func (p *player2) ProcessWinOrNot(win bool) {
	if win {
		p.wins++
	}
}

// 몬티홀 문제
type MontyHall struct {
	doors [DoorCount + 1]bool // 상품이 포함되어 있는 문들 ( true: 스포츠카(당첨), false: 염소(꽝) )
	pick  int                 // 플레이어가 선택한 문 번호
}

// 새로운 몬티홀 문제를 생성한다.
func NewMontyHall() *MontyHall {
	mh := MontyHall{}
	mh.doors[pickDoor()] = true // 1, 2, 3 문중 하나에 당첨을 설정
	return &mh
}

// 1, 2, 3 문 중 하나의 번호를 랜덤하게 선택한다
func pickDoor() int {
	return 1 + rand.Intn(DoorCount)
}

// 사회자가 플레이어가 선택한 문 또는 당첨이 아닌 문을 연다.
func (mh MontyHall) Preview() int {
	n := pickDoor()
	for mh.doors[n] || n == mh.pick {
		n = pickDoor()
	}
	return n
}

// 당첨인지 확인한다.
func (mh MontyHall) IsWinPrize() bool {
	return mh.doors[mh.pick]
}

// 선택을 바꾼다.
func (mh MontyHall) SwitchDoor(preview int) int {
	for i := 1; i <= DoorCount; i++ {
		if i != mh.pick && i != preview {
			return i
		}
	}
	panic("wtf! this should not happen!")
}

// 몬티홀 문제에 대해 repeats횟수만큼 시행하여,
// 선택을 바꿨을 때와 바꾸지 않았을 때의 당첨 확률을 구해서 출력한다.
func main() {

	fmt.Println("Hello, Monty-Hall Problem!")
	defer fmt.Println("Bye, Monty-Hall Problem!")

	// 사용자 플레이어
	players := []Player{&player{}}
	repeats := 10

	// players (player1은 항상 선택을 유지, player2는 항상 선택을 바꾼다)
	// players := []Player{&player1{}, &player2{}}
	// repeats := 10000000

	// repeats 횟수 만큼 반복 시행
	for i := 0; i < repeats; i++ {
		// 0. 새로운 문제 생성
		mh := NewMontyHall()

		// 각 플레이어에 대해
		for _, p := range players {

			// 1. 플레이어가 임의의 문을 선택한다.
			mh.pick = p.PickDoor()

			// 2. 사회자가 미리보기 문을 열어준다.
			preview := mh.Preview()

			// 3. 선택을 바꿀까?
			if p.SwitchChoice() {
				mh.pick = mh.SwitchDoor(preview)
			}

			// 4. 당첨 여부 확인
			p.ProcessWinOrNot(mh.IsWinPrize())
		}
	}

	// 플레이어 별 결과 출력
	for i, p := range players {
		fmt.Printf("player[%d]: %d/%d(%.2f%%)\n", i, p.Wins(), repeats, float64(p.Wins()*100)/float64(repeats))
	}
}

/* OUTPUT for player
D:\gitworks\go-study\monty-hall>go run .
Hello, Monty-Hall Problem!
pick door(1/2/3) 1
change your mind? (y/n) n
staying your mind
you lost the game!
pick door(1/2/3) 1
change your mind? (y/n) n
staying your mind
congraturations! you win the prize!
pick door(1/2/3) 1
change your mind? (y/n) n
staying your mind
you lost the game!
pick door(1/2/3) 1
change your mind? (y/n) n
staying your mind
you lost the game!
pick door(1/2/3) 1
change your mind? (y/n) n
staying your mind
you lost the game!
pick door(1/2/3) 1
change your mind? (y/n) n
staying your mind
you lost the game!
pick door(1/2/3) 1
change your mind? (y/n) n
staying your mind
you lost the game!
pick door(1/2/3) 1
change your mind? (y/n) n
staying your mind
you lost the game!
pick door(1/2/3) 1
change your mind? (y/n) n
staying your mind
you lost the game!
pick door(1/2/3) 1
change your mind? (y/n) n
staying your mind
you lost the game!
player[0]: 1/10(10.00%)
Bye, Monty-Hall Problem!

D:\gitworks\go-study\monty-hall>go run .
Hello, Monty-Hall Problem!
pick door(1/2/3) 1
change your mind? (y/n) y
changing your mind
congraturations! you win the prize!
pick door(1/2/3) 1
change your mind? (y/n) y
changing your mind
congraturations! you win the prize!
pick door(1/2/3) 1
change your mind? (y/n) y
changing your mind
congraturations! you win the prize!
pick door(1/2/3) 1
change your mind? (y/n) y
changing your mind
you lost the game!
pick door(1/2/3) 1
change your mind? (y/n) y
changing your mind
you lost the game!
pick door(1/2/3) 1
change your mind? (y/n) y
changing your mind
congraturations! you win the prize!
pick door(1/2/3) 1
change your mind? (y/n) y
changing your mind
congraturations! you win the prize!
pick door(1/2/3) 1
change your mind? (y/n) y
changing your mind
you lost the game!
player[0]: 7/10(70.00%)
Bye, Monty-Hall Problem!
*/

/* OUTPUT for player1, player2
Hello, Monty-Hall Problem!
player[0]: 3333888/10000000(33.34%)
player[1]: 6666080/10000000(66.66%)
Bye, Monty-Hall Problem!
*/
