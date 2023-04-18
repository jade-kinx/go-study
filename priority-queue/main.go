package main

import (
	"fmt"
	"gostudy/pkg/rng"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// 응급 환자 정보
// 응급 환자는 방문한 시점부터 1초당 hp가 1씩 감소하며 0이 되면 사망한다.
// 어린이/노약자/여성 우선?
type Patient struct {
	id      int
	age     int
	sex     bool
	hp      int
	visitAt time.Time
}

// 환자가 죽었나?
func (p Patient) IsDead() bool {
	return int(time.Now().Sub(p.visitAt).Seconds()) >= p.hp
}

func main() {
	fmt.Println("Hello, Go!")
	defer fmt.Println("Bye, Go!")

	// 환자 수
	patients_count := 10000

	// 응급 환자 큐
	patients := NewChannel[Patient](100)

	wg := sync.WaitGroup{}
	doctors := runtime.NumCPU() / 2

	// 시작 시간
	begin := time.Now()

	// 환자 발생!!!
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < patients_count; i++ {

			// 랜덤한 환자 생성( hp: 10-90 )
			patient := Patient{id: i, hp: rng.NextInRange(10, 90), visitAt: time.Now()}

			// 환자 대기열에 추가
			if err := patients.Push(patient, patient.hp); err != nil {
				fmt.Printf("enque: err=%v", err)
				continue
			}
		}
	}()

	// 사망자수, 처치 수
	var dead, cured int64
	treated := make([]int, doctors)

	// 환자 진료
	wg.Add(doctors)
	for i := 0; i < doctors; i++ {
		go func(doctor int) {
			defer wg.Done()
			fmt.Printf("doctor[%d] started to cure patients...\n", doctor)
			defer fmt.Printf("doctor[%d] says: I'm done!\n", doctor)

			// 잠시 대기
			time.Sleep(time.Second * 1)

			for {
				// 환자 대기열에서 환자를 호출
				patient, err := patients.Deque()
				if err != nil {
					// fmt.Printf("doctor[%d] says: no more patients!\n", doctor)
					break
				}

				// 죽었나?
				if patient.IsDead() {
					fmt.Printf("patient dead!!! hp=%d, visitAt=%v\n", patient.hp, patient.visitAt)
					atomic.AddInt64(&dead, 1)
					continue
				}

				// 치료한다.
				treattime := 100 - patient.hp
				patient.hp = 100
				atomic.AddInt64(&cured, 1)
				treated[doctor]++
				// fmt.Printf("doctor[%d] cured patient[%+v]\n", doctor, patient)

				// 치료 시간만큼 대기
				time.Sleep(time.Millisecond * time.Duration(treattime))
			}
		}(i)
	}

	// wait for all go-routine done
	wg.Wait()

	// 결과 출력
	fmt.Printf("dead: %d, cured: %d, elapsed: %.2f(s)\n", dead, cured, time.Since(begin).Seconds())
	fmt.Printf("doctor treated: %+v\n", treated)

	/* OUTPUT
	D:\gitworks\go-study\priority-queue>go run .
	Hello, Go!
	doctor[9] started to cure patients...
	doctor[6] started to cure patients...
	doctor[1] started to cure patients...
	doctor[2] started to cure patients...
	doctor[0] started to cure patients...
	doctor[3] started to cure patients...
	doctor[7] started to cure patients...
	doctor[8] started to cure patients...
	doctor[4] started to cure patients...
	doctor[5] started to cure patients...
	doctor[1] says: I'm done!
	doctor[0] says: I'm done!
	doctor[7] says: I'm done!
	doctor[4] says: I'm done!
	doctor[6] says: I'm done!
	doctor[8] says: I'm done!
	doctor[2] says: I'm done!
	doctor[9] says: I'm done!
	doctor[3] says: I'm done!
	doctor[5] says: I'm done!
	dead: 0, cured: 10000, elapsed: 59.15(s)
	doctor treated: [1001 1011 1015 988 983 995 991 1020 1000 996]
	Bye, Go!
	*/
}
