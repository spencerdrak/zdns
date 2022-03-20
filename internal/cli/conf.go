package cli

import (
	"net"
	"time"

	"github.com/zmap/zdns/pkg/zdns"
)

type ZdnsRun struct {
	GlobalConf       GlobalConf
	Servers          string
	LocalAddr        string
	LocalIF          string
	ConfigFile       string
	Timeout          int
	IterationTimeout int
	Class            string
	NanoSeconds      bool
	ModuleFlags      ModuleFlags
}

type ModuleFlags struct {
	Ipv4Lookup    bool
	Ipv6Lookup    bool
	BlacklistFile string
	MxCacheSize   int
}

type GlobalConf struct {
	Threads               int
	Timeout               time.Duration
	IterationTimeout      time.Duration
	Retries               int
	AlexaFormat           bool
	MetadataFormat        bool
	NameServerInputFormat bool
	IterativeResolution   bool

	ResultVerbosity string
	IncludeInOutput string
	OutputGroups    []string

	MaxDepth             int
	CacheSize            int
	GoMaxProcs           int
	Verbosity            int
	TimeFormat           string
	PassedName           string
	NameServersSpecified bool
	NameServers          []string
	TCPOnly              bool
	UDPOnly              bool
	LocalAddrSpecified   bool
	LocalAddrs           []net.IP

	InputHandler  InputHandler
	OutputHandler OutputHandler

	InputFilePath    string
	OutputFilePath   string
	LogFilePath      string
	MetadataFilePath string

	NamePrefix     string
	NameOverride   string
	NameServerMode bool

	Module          string
	RequestedModule zdns.Module
	Class           uint16
}

type Metadata struct {
	Names       int            `json:"names"`
	Status      map[string]int `json:"statuses"`
	StartTime   string         `json:"start_time"`
	EndTime     string         `json:"end_time"`
	NameServers []string       `json:"name_servers"`
	Timeout     int            `json:"timeout"`
	Retries     int            `json:"retries"`
	Conf        *GlobalConf    `json:"conf"`
}

var RootServers = [...]string{
	"198.41.0.4:53",
	"192.228.79.201:53",
	"192.33.4.12:53",
	"199.7.91.13:53",
	"192.203.230.10:53",
	"192.5.5.241:53",
	"192.112.36.4:53",
	"198.97.190.53:53",
	"192.36.148.17:53",
	"192.58.128.30:53",
	"193.0.14.129:53",
	"199.7.83.42:53",
	"202.12.27.33:53",
}
