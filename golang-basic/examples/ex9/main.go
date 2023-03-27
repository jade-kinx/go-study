package main

import (
	"fmt"
	"os"
)

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

func Info(log Logger, message string) error {
	return log.WriteLine(message)
}

func Error(log Logger, message string) error {
	return log.WriteLine(fmt.Sprintf("[Error]: %s", message))
}

var log Logger

func main() {
	log = ConsoleLog{}
	// log = FileLog{"./mylog.log"}
	// log = RemoteLog{"http://xxxx/xxxx"}
	// 또는 로그 체인에 추가

	Info(log, "hello, world!")
}
