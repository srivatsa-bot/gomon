package logger

import (
	"bufio"
	"fmt"
	"io"
	"log"
)

// to print output in user terminal, takes in piped output
func LogOutput(reader io.ReadCloser, prefix string) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fmt.Printf("[%s] %s\n", prefix, scanner.Text())
	}
}

// function to log something takes string and multiple parameters as input using variadic
func Info(format string, p ...interface{}) { //using empty inteface because it cna store value of any type int,srting,float etc
	log.Printf(format+"\n", p...)
}

//to log erros

func Error(format string, p ...interface{}) {
	log.Printf("ERROR: "+format+"\n", p...)
}
