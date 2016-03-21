package log

import (
	"fmt"
	"io"
	"log"
	"mustard/base/conf"
	"os"
)

/*
* Debug(Log Level) - Info - Warning - Error -
 */

// combination version, it hide log.Logger method.
type logS struct {
	logI  *log.Logger
	level int
}

func (l *logS) Debugf(format string, v ...interface{}) {
	if l.level <= *conf.Conf.LogV {
		l.logI.SetPrefix("[Debug]")
		l.logI.Output(2, fmt.Sprintf(format, v...))
	}
}

func (l *logS) Debug(v ...interface{}) {
	if l.level <= *conf.Conf.LogV {
		l.logI.SetPrefix("[Debug]")
		l.logI.Output(2, fmt.Sprintln(v...))
	}
}

func (l *logS) Info(v ...interface{}) {
	l.logI.SetPrefix("[Info]")
	l.logI.Output(2, fmt.Sprintln(v...))
}

func (l *logS) Infof(format string, v ...interface{}) {
	l.logI.SetPrefix("[Info]")
	l.logI.Output(2, fmt.Sprintf(format, v...))
}

func (l *logS) Warning(v ...interface{}) {
	l.logI.SetPrefix("[Warning]")
	l.logI.Output(2, fmt.Sprintln(v...))
}

func (l *logS) Warningf(format string, v ...interface{}) {
	l.logI.SetPrefix("[Warning]")
	l.logI.Output(2, fmt.Sprintf(format, v...))
}

func (l *logS) Error(v ...interface{}) {
	l.logI.SetPrefix("[Error]")
	l.logI.Output(2, fmt.Sprintln(v...))
}

func (l *logS) Errorf(format string, v ...interface{}) {
	l.logI.SetPrefix("[Error]")
	l.logI.Output(2, fmt.Sprintf(format, v...))
}

func (l *logS) Fatal(v ...interface{}) {
	l.logI.SetPrefix("[Fatal]")
	l.logI.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

func (l *logS) Fatalf(format string, v ...interface{}) {
	l.logI.SetPrefix("[Fatal]")
	l.logI.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// VLog user pairly with Debug
// chain function call
func (l *logS) VLog(level int) *logS {
	l.level = level
	return l
}

var Log logS

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
	Log.logI = log.New(writer, "", log.LstdFlags|log.Lshortfile)
}

// test
/*
import "mustard/base/log
var LOG log.Log
func main() {
	LOG.Info("hello world")
	LOG.VLog(1).Debug("debug message1")
	LOG.VLog(2).Debug("debug message2")
	LOG.VLog(3).Debug("debug message3")
}
//*/
