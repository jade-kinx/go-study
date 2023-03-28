# 키워드/연산자/내장함수
## 키워드
* GO 언어의 키워드는 모두 25개 (참고: https://go.dev/ref/spec#Keywords)  
* C++(17)의 키워드는 84개, C#의 키워드는 80개 이상

```
break        default      func         interface    select
case         defer        go           map          struct
chan         else         goto         package      switch
const        fallthrough  if           range        type
continue     for          import       return       var
```

## 연산자

#### 산술연산자
| 연산자 | 설명 | 피연산자 |
| :---: | --- | --- |
| + | 덧셈 | 정수,실수,복소수,문자열 |
| - | 뺄셈 | 정수,실수,복소수 |
| * | 곱셈 | 정수,실수,복소수 |
| / | 나눗셈 | 정수,실수,복소수 |
| % | 나머지 | 정수,실수,복소수 |

#### 비트연산자
| 연산자 | 설명 | 피연산자 |
| :---: | --- | --- |
| & | AND 비트 연산 | 정수 |
| &#124; | OR 비트 연산 | 정수 |
| ^ | XOR 비트 연산 | 정수 |
| &^ | 비트 클리어 | 정수 |
| << | 레프트 시프트 | 정수 |
| >> | 라이트 시프트 | 정수 |

#### 비교연산자
| 연산자 | 설명 | 피연산자 |
| :---: | --- | --- |
| == | 같다 | any |
| != | 다르다 | any |
| > | 크다 | Number, string |
| >= | 크거나 같다 | Number, string |
| < | 작다 | Number, string |
| <= | 작거나 같다 | Number, string |

#### 논리연산자
| 연산자 | 설명 | 반환값 |
| :---: | --- | --- |
| && | AND 논리 연산 | 양변이 모두 true면 true |
| &#124;&#124; | OR 논리 연산 | 양변중 하나라도 true면 true |
| ! | NOT 논리 연산 | 피연산자가 true면 false 아니면 true |

#### 기타연산자
| 연산자 | 설명 | 
| :---: | --- |
| [] | 배열/맵의 요소에 접근 |
| . | 구조체/패키지 요소에 접근 |
| & | 변수의 메모리 주소값 반환 |
| * | 포인터 변수가 가리키는 메모리에 접근 |
| ... | 가변인자, 배열/슬라이스 구조분할 |
| <- | 채널에 값을 입출력 |


## 내장함수

| 내장함수 | 설명 | 인자
| :---: | --- | --- |
| close(ch) | 채널을 닫는다 | chan |
| len(s) | s 의 길이를 얻는다 | string, 배열, 슬라이스, 맵, 채널 등 |
| cap(s) | s의 할당된 크기를 얻는다 | 배열, 슬라이스, 채널 등 |
| new(T) | T 타입의 객체를 생성 (초기화 X) | |
| make(T, [n], [m]) | T 타입의 객체를 생성 (초기화 O) | |
| append(s, ...e) | s에 e를 추가하여 반환 ||
| copy(dst, src) | dst에 src를 복사 ||
| delete(m, k) | remove element m[k] from map m ||
| complex, real, img | complex 타입의 객체 생성, 소수부, 지수부 반환 ||
| panic, recover | 런타임 오류 발생, 런타임 오류 복구 ||


## 내장인터페이스

```go
// The error built-in interface type is the conventional interface for
// representing an error condition, with the nil value representing no error.
type error interface {
	Error() string
}
```
