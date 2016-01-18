package conf

import (
	"flag"
)

type ExampleType struct {
	Port *int
	Address *string
}

var ExampleConf = ExampleType{
	Port: flag.Int("example_port", 9001, "example port"),
	Address:     flag.String("api_address", "127.0.0.1:9100;", "rest address of example api"),
}
