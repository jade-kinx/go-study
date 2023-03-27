# 모듈과 패키지

## 패키지
* 코드를 묶는 단위이며 모든 코드는 디렉토리 단위로 패키지에 속함  
* `main` 패키지는 실행파일 생성(`main()` 함수 포함), 그 외 패키지는 라이브러리 패키지  
* 패키지 내의 모든 자원은 패키지 내에서 공유됨 (함수, 전역 변수, 타입 정의 등)  
* 패키지 외부로 공개하고 싶은 자원은 첫글자를 대문자로...  

### 외부 패키지 임포트
```go
package mypackage   // 현재 패키지 이름 선언

import (
    "github.com/gin-gonic/gin"          // 외부 모듈 임포트(패키지 이름은 gin)
    log "github.com/sirupsen/logrus"    // 패키지 logrus를 log로 별칭
)

func main() {
    r := gin.Default()             // 외부 패키지 이름으로 사용
    log.Info("hello, world")       // 별칭 log로 사용
    // ...
}
```
```bash
# 사용하는 모듈을 다운로드(pkg.go.dev에 배포)하고, 사용하지 않는 모듈을 제거
$ go mod tidy
```

### 패키지 초기화 함수 `func init()`
* `init()` 함수는 패키지가 `import`될때, 가장 먼저 실행  

```go
package main

import "fmt"

// 패키지 초기화 함수
func init() {
    fmt.Println("하이!")
}

func main() {
    fmt.Println("Hello, World")
}

// 중복 정의도 가능
func init() {
    fmt.Println("헬로?")
}
/* OUTPUT
하이!
헬로?
Hello, World
*/
```

* init() 함수 실행 순서
  * 가장 상위 디렉토리 패키지 부터 하위 디렉토리 패키지 순으로
  * 패키지 파일명 정렬 순으로
  * 파일내 위에서 아래 순으로

## 모듈
* 모듈은 패키지의 모음(모듈:패키지=1:N)  
* 패키지 종속성 관리 목적  
* `Golang 1.16` 부터 기본 사양이므로 반드시 모듈 초기화 필요  

```bash
$ go mod init mymodule          # 실행 파일명은 mymodule
$ go mod init mymodule/test     # 실행 파일명은 test
$ go mod init github.com/{github_id}/mymodule # github를 통해 pkg.go.dev에 배포
```

