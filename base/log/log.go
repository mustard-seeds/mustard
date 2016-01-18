package log

import (
	"io"
	"log"
	"os"
	"mustard/base/conf"
)

type logS struct {
	logI  *log.Logger
}

func (l *logS) Info(v ...interface{}) {
	l.logI.SetPrefix("[Info]")
	l.logI.Println(v)
}

var Log logS
func init() {
	var writer io.Writer
	if *conf.Conf.UseStdout {
		writer = os.Stdout
	} else if logDir := *conf.Conf.LogDir; logDir != "" {
		f, err := os.OpenFile(logDir, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		writer = f
	} else {
		writer = os.Stdout
	}
	if writer != os.Stdout {
		writer = io.MultiWriter(writer, os.Stdout)
	}
	Log.logI = log.New(writer, "", log.LstdFlags|log.Lshortfile)
}
