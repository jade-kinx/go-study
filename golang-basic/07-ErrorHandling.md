# 에러 처리
* Go에서는 try-catch 와 같은 예외 처리 구문이 없음

## 에러/오류
* ~~컴파일 에러~~
* error: 런타임 에러(error를 반환: 'no such file or directory')
* panic: 런타임 오류(divide by 0, index out of range, etcs, ...)

### 에러 타입 (내장 타입)
```go
type error interface {
    Error() string
}
```

```go
// 커스텀 에러 타입
type MyError struct {
	code int
}
func (e *MyError) Error() string {
	return fmt.Sprintf("MyError: %d", e.code)
}

var (
	ErrMyError2 = errors.New("MyError type 2")
	ErrMyError3 = errors.New("MyError type 3")
)

func getSomeError() error {
	return &MyError{100}
}

func wrapErrorFunction() error {
	if err := getSomeError(); err != nil {
		return fmt.Errorf("wrapErrorFunction() err=%w", err)
	}

	// ...
	return nil
}
```

## 패닉 & 리커버

* 패닉 발생시 죽을 것인가? 처리하고 살릴 것인가?
* 개발 단계에서는 빠르게 죽어서, 원인을 파악/수정하는 편이 좋음
* 서비스 단계에서는 죽지 않고, 원인을 logging/notify 해주어서 수정할 수 있도록 정보를 제공해야 한다.  
* 어디에서 복구할 것인가? (애초에 패닉이 발생하지 않도록 작성해야 하고, 주로 콜스택의 최상단에서 복구...)  

```go
func ReadFromUrl(url string) ([]byte, error) {
    // ReadFromUrl() 함수 내 어딘가에서 panic이 발생할 경우 복구
    defer func() {
        if r := recover(); r != nil {
            fmt.Println(r)
        }
    }()

    // 패닉 발생
    panic(fmt.Sprintf("intended panic: %s", url))

    // 또는 여기 안에서 패닉 발생
    return DoSomeWork(url)
}
```
