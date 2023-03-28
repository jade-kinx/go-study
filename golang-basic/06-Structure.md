# 구조체

## 구조체
* 구조체는 필드의 집합  
* 구조체는 함수를 가지지 않지만 리시버를 붙여 메소드를 선언하여 클래스처럼 사용할 수 있음  

### 구조체 정의와 초기화
```go
type User struct {
    Name string
    Id  string
    Age int
}

func main() {
    var empty User  // 각 필드의 zero-value로 초기화
    john := User{"john doe", "john", 18}
    jane := User{Name: "jade done", Id: "jane", Age: 20}

    fmt.Println(empty)
    fmt.Println(john)
    fmt.Println(jane)
    /* OUTPUT
    {  0}
    {john doe john 18}
    {jade done jane 20}
    */
}

```

### 메모리 할당 (패딩)

```go
type Padding struct {
    a int32             // 4 바이트
    b int16             // 2 바이트
    //_ padding         // 2 바이트
    c int32             // 4 바이트
    d byte              // 1 바이트
    //_ padding         // 3 바이트
}

```
* 크기가 가장 큰 변수 기준으로 패딩  
* `string` 타입의 크기는 내용에 관계 없이 항상 `16`(32비트 컴퓨터에서는 `12`, `StringHeader`)  
* `slice` 타입의 크기는 `SliceHeader{uintptr, int, int}`  
* 크기가 큰 변수부터 차례로 써주면 메모리 낭비를 줄일 수 있지만, 굳이?  

## 메소드
* 메소드는 객체에 붙여서(리시버 정의) 사용하는 함수  
* 구조체 등 객체를 클래스처럼 사용할 수 있도록 함  
* 메소드로 선언된 함수는 리시버 객체를 통해서만 호출할 수 있음  

```go

type MyInt int

// (i *MyInt)가 리시버
func (i *MyInt) Increment(count int) {
	*i = MyInt(int(*i) + count)
}

type User struct {
	name string
	age  int
}

// 객체를 복사(pass-by-value)하기 때문에 caller객체와 user 객체는 다른 객체
func (user User) Name() string {
	return user.name
}
func (user User) Age() int {
	return user.age
}
// pass-by-pointer: 내용을 변경하면 caller에도 반영
func (user *User) SetName(name string) {
	user.name = name
}
func (user *User) SetAge(age int) {
	user.age = age
}

func main() {
	var i MyInt = 100
	i.Increment(100)
	fmt.Println(i)

	user := User{"john.doe", 18}
	fmt.Println(user.Name(), user.Age())
	user.SetName("jane.doe")
	user.SetAge(20)
	fmt.Println(user.Name(), user.Age())
    /* OUTPUT
    200
    john.doe 18
    jane.doe 20
    */
}
```
* C#의 확장메소드(Method Extension)와 비슷한 개념으로 생각할 수 있음  

```go
// SetName() 메소드는 결국 SetUserName()과 같이 동작
func (user *User) SetName(name string) {
    user.name = name
}
func SetUserName(user *User, name string) {
    user.name = name
}
```

  
## 인터페이스
* 구조체는 필드의 집합, 인터페이스는 메소드의 집합  
* 메소드의 구현이 아닌, 추상화된 메소드 원형만 정의  
* 관례적으로 인터페이스는 `-er` 어미를 사용(ex: `io.Reader`, `io.Writer`, `io.Closer`)  

```go
// 로깅 인터페이스
type Logger interface {
	WriteLine(message string) error
}

// 콘솔 로거
type ConsoleLog struct {
}
func (_ ConsoleLog) WriteLine(message string) error {
	fmt.Println(message)
	return nil
}

// 파일 로거
type FileLog struct {
	logFilePath string
}
func (log *FileLog) WriteLine(message string) error {
	f, err := os.Open(log.logFilePath)
	if err != nil {
		return fmt.Errorf("could not open logFilePath(%s)", log.logFilePath)
	}
	defer f.Close()
	_, err = f.WriteString(message + "\n")
	return err
}

// 원격 로그 (나중에 구현해야할 필요가 생겼다고 가정)
type RemoteLog struct {
	remoteUrl string
}
func (log *RemoteLog) WriteLine(message string) error {
	return fmt.Errorf("not implemented yet! message=%s", message)
}

// 정보 출력
func Info(log Logger, message string) error {
	return log.WriteLine(message)
}
// 에러 출력
func Error(log Logger, message string) error {
	return log.WriteLine(fmt.Sprintf("[Error]: %s", message))
}

var log Logger

func main() {
    log = ConsoleLog{}
    // 상황에 따라 로거를 바꿔 끼거나
	// log = FileLog{"./mylog.log"}
	// log = RemoteLog{"http://xxxx/xxxx"}
	// 또는 로그 체인에 추가

	Info(log, "hello, world!")  
}

```

### 덕타이핑(duck typing)

> 만약, 어떤 새가 오리처럼 걷고, 오리처럼 헤엄치고, 오리처럼 날면 나는 그 새를 오리라고 부를 것이다.

* 만약, 어떤 객체가 인터페이스에 정의된 메소드 목록을 모두 구현하고 있다면 그 객체의 타입과 관계없이 그 인터페이스로 사용할 수 있음  
* C++, Java, C# 등은 객체가 해당 인터페이스를 구현하고 있음을 명시적으로 선언해야 함  
* 의존성 주입(Dependency Injection)에 사용할 수 있음  
