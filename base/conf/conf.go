package conf

import (
    "flag"
    "mustard/internal/github.com/vharitonsky/iniflags"
)

type ConfType struct {
    LogFile     *string
    LogV        *int
    Stdout      *bool

    UseTLS   *bool
    CertFile *string // server
    KeyFile  *string // server
    CaFile  *string  // client
    ServerHostOverride  *string // client

    Example     ExampleType
    Crawler     CrawlerType
}

var _conf = ConfType{
    LogFile     :   flag.String("log_file", "", "log to file"),
    LogV        :   flag.Int("v", 3, "log level for debug"),
    Stdout      :   flag.Bool("stdout", true, "output stdout or not"),
    UseTLS      :   flag.Bool("use_tls", false, "use tls or not"),
    CertFile    :   flag.String("cert_file", "", "TLS cert file"),
    KeyFile     :   flag.String("key_file", "", "TLS key file"),
    CaFile      :   flag.String("ca_file", "","The file containning the CA root cert file"),
    ServerHostOverride: flag.String("server_host_override","x.a.com", "The server name use to verify the hostname returned by TLS handshake"),

    Example     :       ExampleConf,
    Crawler     :       CrawlerConf,
}

var Conf *ConfType
func init() {
    Conf = &_conf
    iniflags.Parse()
}
