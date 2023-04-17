package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// 1. json.Unmarshaller 인터페이스 원형
// type Unmarshaler interface {
// 	UnmarshalJSON([]byte) error
// }

// 2. myTime 구조체
type myTime time.Time

// 3. myTime 구조체에 대해 UnmarshalJSON() 메소드를 구현하여 json.Unmarshaller 인터페이스를 구현
// 이제, json.Unmarshal() 메소드 내에서 myTime 구조체를 언마샬링 할때는 이 부분을 사용하게 됨
func (mt *myTime) UnmarshalJSON(bs []byte) error {
	var s string
	err := json.Unmarshal(bs, &s)
	if err != nil {
		return err
	}
	t, err := time.ParseInLocation("2006-01-02", s, time.UTC)
	if err != nil {
		return err
	}
	*mt = myTime(t)
	return nil
}

func (mt myTime) String() string {
	return time.Time(mt).Format("2006-01-02")
}

// 4. 인터페이스 변수에 구조체를 대입하는 예제 였던 것으로 추정되는 것
// 프로그램 실행 코드와는 관계 없음
// var _ json.Unmarshaler = &myTime{}

// StartAt은 myTime이란 구조체로 json.Unmarshal() 호출시 마샬링 할 수 없음!
// myTime 구조체에서 json.Unmarshaller 인터페이스를 구현(myTime.UnmarshalJSON)하여 json.Unmarshal() 호출시 마샬링 할 수 있음
type MyClass struct {
	StartAt     myTime `json:"start_at" binding:"required"`
	ChallengeID uint   `json:"challenge_id" gorm:"index" binding:"required"`
}

// json unmarshaling 할때, myTime 구조체에 대해서도 unmarshaling 할 수 있도록 json.Unmarshaller 인터페이스를 구현한 예제
func main() {
	fmt.Println("Hello, Go!")
	defer fmt.Println("Bye, Go!")

	str := `{"start_at": "2023-04-07", "challenge_id": 100}`

	var mc MyClass
	if err := json.Unmarshal([]byte(str), &mc); err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("mc=%+v\n", mc)
}

/* OUTPUT #1: 3. 항목을 주석처리하여 myTime.UnmarshalJSON() 메소드 구현이 없는 경우
Hello, Go!
error: json: cannot unmarshal string into Go struct field MyClass.start_at of type main.myTime
Bye, Go!
*/

/* OUTPUT #2: myTime.UnmarshalJSON() 구현으로 json.Unmarshaller 인터페이스로 사용할 수 있는 경우
Hello, Go!
mc={StartAt:2023-04-07 ChallengeID:100}
Bye, Go!
*/
