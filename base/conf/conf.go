package conf

import (
	"flag"
	"mustard/internal/github.com/vharitonsky/iniflags"
)

type ConfType struct {
	UseStdout *bool
	LogDir   *string

	Example  ExampleType
}

var _conf = ConfType{
	UseStdout: flag.Bool("use_stdout", true, "output to stdout"),
	LogDir:   flag.String("log_dir", "", "log to file"),

	Example:      ExampleConf,
}

var Conf *ConfType
func init() {
	Conf = &_conf
	iniflags.Parse()
}
