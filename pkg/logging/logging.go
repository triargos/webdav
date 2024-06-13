package logging

import (
	"log"
	"os"
)

type Loggers struct {
	Info  *log.Logger
	Error *log.Logger
}

var Log *Loggers

func InitLoggers() {
	Log = &Loggers{
		Info:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		Error: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
