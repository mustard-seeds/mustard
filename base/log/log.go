package log

import (
    "io"
    "log"
    "fmt"
    "os"
    "flag"
    "strings"
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

func (l *logS)Debugf(format string, v ...interface{}) {
    if l.level <= *conf.Conf.LogV {
        l.logI.SetPrefix("[Debug]")
        l.logI.Output(2, fmt.Sprintf(format, v...))
    }
}

func (l *logS)Debug(v ...interface{}) {
    if l.level <= *conf.Conf.LogV {
        l.logI.SetPrefix("[Debug]")
        l.logI.Output(2, fmt.Sprintln(v...))
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
    // dump flags in log, because of dependency cycle
    dumpFlags()
}


func escapeUsage(s string) string {
    return strings.Replace(s, "\n", "\n    # ", -1)
}

func quoteValue(v string) string {
    if !strings.ContainsAny(v, "\n#;") && strings.TrimSpace(v) == v {
        return v
    }
    v = strings.Replace(v, "\\", "\\\\", -1)
    v = strings.Replace(v, "\n", "\\n", -1)
    v = strings.Replace(v, "\"", "\\\"", -1)
    return fmt.Sprintf("\"%s\"", v)
}
func dumpFlags() {
    Info("=================Dump Flags=========================================")
    flag.VisitAll(func(f *flag.Flag) {
        if f.Name != "config" && f.Name != "dumpflags" {
            Infof("%s = %s # %s\n", f.Name, quoteValue(f.Value.String()), escapeUsage(f.Usage))
        }
    })
    Info("=================Dump Flags Finish===================================")
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