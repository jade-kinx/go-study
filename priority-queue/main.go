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

const (
	PATIENT_COUNT = 10000 // 총 환자 수
)

// 우선순위 큐를 사용한 예
func runWithPriorityQueue() {
	// 응급 환자 큐
	patients := NewChannel[Patient](PATIENT_COUNT)

	wg := sync.WaitGroup{}
	doctors := runtime.NumCPU() / 2

	// 시작 시간
	begin := time.Now()

	// 환자 발생!!!
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < PATIENT_COUNT; i++ {

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
					fmt.Printf("critical!!! patient dead!!! %+v\n", patient)
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
}

// Go 채널을 사용한 예
func runWithGoChannel() {
	// 응급 환자 큐
	patients := make(chan Patient, PATIENT_COUNT)

	wg := sync.WaitGroup{}
	doctors := runtime.NumCPU() / 2

	// 시작 시간
	begin := time.Now()

	// 환자 발생!!!
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < PATIENT_COUNT; i++ {

			// 랜덤한 환자 생성( hp: 10-90 )
			patient := Patient{id: i, hp: rng.NextInRange(10, 90), visitAt: time.Now()}

			// 환자 대기열에 추가
			patients <- patient
		}

		// 채널을 닫는다.
		close(patients)
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

			// 대기열에서 환자를 호출
			for patient := range patients {
				// 죽었나?
				if patient.IsDead() {
					fmt.Printf("critical!!! patient dead!!! %+v\n", patient)
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
}

func main() {
	fmt.Println("Hello, Go!")
	defer fmt.Println("Bye, Go!")

	// 채널을 이용한 예를 실행
	runWithGoChannel()
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
	critical!!! patient dead!!! {id:7169 age:0 sex:false hp:29 visitAt:{wall:13909285466389033824 ext:5396801 loc:0x2d03c0}}
	critical!!! patient dead!!! {id:7174 age:0 sex:false hp:15 visitAt:{wall:13909285466389033824 ext:5396801 loc:0x2d03c0}}
	critical!!! patient dead!!! {id:7181 age:0 sex:false hp:22 visitAt:{wall:13909285466389033824 ext:5396801 loc:0x2d03c0}}
	critical!!! patient dead!!! {id:7182 age:0 sex:false hp:34 visitAt:{wall:13909285466389033824 ext:5396801 loc:0x2d03c0}}
	critical!!! patient dead!!! {id:7187 age:0 sex:false hp:25 visitAt:{wall:13909285466389033824 ext:5396801 loc:0x2d03c0}}
	critical!!! patient dead!!! {id:7191 age:0 sex:false hp:33 visitAt:{wall:13909285466389033824 ext:5396801 loc:0x2d03c0}}
	critical!!! patient dead!!! {id:7195 age:0 sex:false hp:33 visitAt:{wall:13909285466389033824 ext:5396801 loc:0x2d03c0}}
	critical!!! patient dead!!! {id:7200 age:0 sex:false hp:21 visitAt:{wall:13909285466389033824 ext:5396801 loc:0x2d03c0}}
	critical!!! patient dead!!! {id:7201 age:0 sex:false hp:22 visitAt:{wall:13909285466389033824 ext:5396801 loc:0x2d03c0}}
	...(생략)...
	critical!!! patient dead!!! {id:9949 age:0 sex:false hp:26 visitAt:{wall:13909285466389592424 ext:5955401 loc:0x2d03c0}}
	critical!!! patient dead!!! {id:9952 age:0 sex:false hp:17 visitAt:{wall:13909285466389592424 ext:5955401 loc:0x2d03c0}}
	critical!!! patient dead!!! {id:9953 age:0 sex:false hp:25 visitAt:{wall:13909285466389592424 ext:5955401 loc:0x2d03c0}}
	critical!!! patient dead!!! {id:9957 age:0 sex:false hp:38 visitAt:{wall:13909285466389592424 ext:5955401 loc:0x2d03c0}}
	critical!!! patient dead!!! {id:9958 age:0 sex:false hp:37 visitAt:{wall:13909285466389592424 ext:5955401 loc:0x2d03c0}}
	critical!!! patient dead!!! {id:9961 age:0 sex:false hp:19 visitAt:{wall:13909285466389592424 ext:5955401 loc:0x2d03c0}}
	critical!!! patient dead!!! {id:9963 age:0 sex:false hp:17 visitAt:{wall:13909285466389592424 ext:5955401 loc:0x2d03c0}}
	critical!!! patient dead!!! {id:9962 age:0 sex:false hp:34 visitAt:{wall:13909285466389592424 ext:5955401 loc:0x2d03c0}}
	critical!!! patient dead!!! {id:9965 age:0 sex:false hp:10 visitAt:{wall:13909285466389592424 ext:5955401 loc:0x2d03c0}}
	critical!!! patient dead!!! {id:9973 age:0 sex:false hp:11 visitAt:{wall:13909285466389592424 ext:5955401 loc:0x2d03c0}}
	critical!!! patient dead!!! {id:9974 age:0 sex:false hp:31 visitAt:{wall:13909285466389592424 ext:5955401 loc:0x2d03c0}}
	critical!!! patient dead!!! {id:9979 age:0 sex:false hp:21 visitAt:{wall:13909285466389592424 ext:5955401 loc:0x2d03c0}}
	critical!!! patient dead!!! {id:9982 age:0 sex:false hp:29 visitAt:{wall:13909285466389592424 ext:5955401 loc:0x2d03c0}}
	critical!!! patient dead!!! {id:9985 age:0 sex:false hp:38 visitAt:{wall:13909285466389592424 ext:5955401 loc:0x2d03c0}}
	critical!!! patient dead!!! {id:9989 age:0 sex:false hp:16 visitAt:{wall:13909285466389592424 ext:5955401 loc:0x2d03c0}}
	critical!!! patient dead!!! {id:9993 age:0 sex:false hp:39 visitAt:{wall:13909285466389592424 ext:5955401 loc:0x2d03c0}}
	critical!!! patient dead!!! {id:9994 age:0 sex:false hp:39 visitAt:{wall:13909285466389592424 ext:5955401 loc:0x2d03c0}}
	critical!!! patient dead!!! {id:9995 age:0 sex:false hp:20 visitAt:{wall:13909285466389592424 ext:5955401 loc:0x2d03c0}}
	critical!!! patient dead!!! {id:9999 age:0 sex:false hp:11 visitAt:{wall:13909285466389592424 ext:5955401 loc:0x2d03c0}}
	doctor[9] says: I'm done!
	doctor[0] says: I'm done!
	doctor[6] says: I'm done!
	doctor[2] says: I'm done!
	doctor[1] says: I'm done!
	doctor[3] says: I'm done!
	doctor[4] says: I'm done!
	doctor[5] says: I'm done!
	doctor[8] says: I'm done!
	doctor[7] says: I'm done!
	dead: 1923, cured: 8077, elapsed: 42.00(s)
	doctor treated: [809 812 795 800 802 814 804 812 825 804]
	Bye, Go!
	*/

	// 우선순위 큐를 이용한 실행
	runWithPriorityQueue()
	/* OUTPUT
	D:\gitworks\go-study\priority-queue>go run .
	Hello, Go!
	doctor[3] started to cure patients...
	doctor[0] started to cure patients...
	doctor[9] started to cure patients...
	doctor[4] started to cure patients...
	doctor[5] started to cure patients...
	doctor[6] started to cure patients...
	doctor[7] started to cure patients...
	doctor[8] started to cure patients...
	doctor[1] started to cure patients...
	doctor[2] started to cure patients...
	doctor[8] says: I'm done!
	doctor[4] says: I'm done!
	doctor[3] says: I'm done!
	doctor[1] says: I'm done!
	doctor[7] says: I'm done!
	doctor[6] says: I'm done!
	doctor[2] says: I'm done!
	doctor[0] says: I'm done!
	doctor[9] says: I'm done!
	doctor[5] says: I'm done!
	dead: 0, cured: 10000, elapsed: 58.91(s)
	doctor treated: [1000 1000 1001 1000 999 1000 1000 1001 999 1000]
	Bye, Go!
	*/
}
