package conf

import (
	"flag"
	"mustard/internal/github.com/vharitonsky/iniflags"
)

type ConfType struct {
	LogFile   	*string
	LogV 		*int
	Stdout		*bool

	Example  	ExampleType
}

var _conf = ConfType{
	LogFile		:	flag.String("log_file", "", "log to file"),
	LogV		:	flag.Int("v", 1, "log level for debug"),
	Stdout		:	flag.Bool("stdout", true, "output stdout or not"),
	Example		:      	ExampleConf,
}

var Conf *ConfType
func init() {
	Conf = &_conf
	iniflags.Parse()
}
