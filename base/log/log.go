package log

import (
	"io"
	"log"
	"fmt"
	"os"
	"mustard/base/conf"
)

/*
* Debug(Log Level) - Info - Warning - Error -
*/

// combination version, it hide log.Logger method.
type logS struct {
	logI  *log.Logger
	level int
}

var _log logS

func Debugf(format string, v ...interface{}) {
	if _log.level <= *conf.Conf.LogV {
		_log.logI.SetPrefix("[Debug]")
		_log.logI.Output(2, fmt.Sprintf(format, v...))
	}
}

func Debug(v ...interface{}) {
	if _log.level <= *conf.Conf.LogV {
		_log.logI.SetPrefix("[Debug]")
		_log.logI.Output(2, fmt.Sprintln(v...))
	}
}

func Info(v ...interface{}) {
	_log.logI.SetPrefix("[Info]")
	_log.logI.Output(2, fmt.Sprintln(v...))
}

func Infof(format string, v ...interface{}) {
	_log.logI.SetPrefix("[Info]")
	_log.logI.Output(2, fmt.Sprintf(format, v...))
}

func Warning(v ...interface{}) {
	_log.logI.SetPrefix("[Warning]")
	_log.logI.Output(2, fmt.Sprintln(v...))
}

func Warningf(format string,v ...interface{}) {
	_log.logI.SetPrefix("[Warning]")
	_log.logI.Output(2, fmt.Sprintf(format,v...))
}

func Error(v ...interface{}) {
	_log.logI.SetPrefix("[Error]")
	_log.logI.Output(2, fmt.Sprintln(v...))
}

func Errorf(format string, v ...interface{}) {
	_log.logI.SetPrefix("[Error]")
	_log.logI.Output(2, fmt.Sprintf(format,v...))
}

func Fatal(v ...interface{}) {
	_log.logI.SetPrefix("[Fatal]")
	_log.logI.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	_log.logI.SetPrefix("[Fatal]")
	_log.logI.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// VLog user pairly with Debug
// chain function call
func VLog(level int) *logS {
	_log.level = level
	return &_log
}


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
	_log.logI = log.New(writer, "", log.LstdFlags|log.Lshortfile)
}
// test
/*
func main() {
	LOG.Info("hello world")
	LOG.VLog(1).Debug("debug message1")
	LOG.VLog(2).Debug("debug message2")
	LOG.VLog(3).Debug("debug message3")
}
//*/