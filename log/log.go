package log

import (
	"fmt"
	"log"
	"os"
)

func Info(format string, args ...interface{}) {
	l := log.New(os.Stdout, "INFO:", log.LstdFlags)
	info := fmt.Sprintf(format, args...)
	l.Println(info)
}

func Error(message string, err error) {
	l := log.New(os.Stdout, "ERROR:", log.LstdFlags|log.Lshortfile)
	info := message + ": " + err.Error()
	l.Println(info)
}
