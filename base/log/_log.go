package log

import (
	"fmt"
	"io"
	"log"
	"mustard/base/conf"
	"os"
	"unsafe"
)

/*
* Debug(Log Level) - Info - Warning - Error -
 */
// inheritance version
type logS struct {
	log.Logger
	level int
}

func (l *logS) Debugf(format string, v ...interface{}) {
	if l.level <= *conf.Conf.LogV {
		l.SetPrefix("[Debug]")
		l.Output(2, fmt.Sprintf(format, v...))
	}
}

func (l *logS) Debug(v ...interface{}) {
	if l.level <= *conf.Conf.LogV {
		l.SetPrefix("[Debug]")
		l.Output(2, fmt.Sprintln(v...))
	}
}

func (l *logS) Info(v ...interface{}) {
	l.SetPrefix("[Info]")
	l.Output(2, fmt.Sprintln(v...))
}

func (l *logS) Warning(v ...interface{}) {
	l.SetPrefix("[Warning]")
	l.Output(2, fmt.Sprintln(v...))
}

func (l *logS) Error(v ...interface{}) {
	l.SetPrefix("[Error]")
	l.Output(2, fmt.Sprintln(v...))
}

// VLog user pairly with Debug
// chain function call
func (l *logS) VLog(level int) *logS {
	l.level = level
	return l
}

var Log *logS

func init() {
	var writer io.Writer
	if logFile := *conf.Conf.LogFile; logFile != "" {
		f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		writer = f
	} else {
		writer = os.Stdout
	}
	if *conf.Conf.Stdout && writer != os.Stdout {
		writer = io.MultiWriter(writer, os.Stdout)
	}
	Log = (*logS)(unsafe.Pointer(log.New(writer, "", log.LstdFlags|log.Lshortfile)))
}
