package logging

import (
	"log"
	"os"
)

type Loggers struct {
	Info      *log.Logger
	Error     *log.Logger
	Operation *log.Logger
}

var Log *Loggers

func InitLoggers() {
	operationLogFile, _ := os.OpenFile("/etc/webdav/operation.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	Log = &Loggers{
		Info:      log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		Error:     log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		Operation: log.New(operationLogFile, "OPERATION: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
