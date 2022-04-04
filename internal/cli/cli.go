package cli

import (
	"io/ioutil"
	"net"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/zmap/dns"
	"github.com/zmap/zdns/internal/util"
	"github.com/zmap/zdns/iohandlers"
	"github.com/zmap/zdns/pkg/zdns"

	log "github.com/sirupsen/logrus"
)

func Run(run ZdnsRun, flags *pflag.FlagSet) {

	logger := log.WithFields(log.Fields{
		"Module": "cli",
	})

	modSet := GenerateModSet()

	if !modSet.HasModule(run.GlobalConf.Module) {
		logger.Fatal("Invalid lookup module specified. Valid modules: ", modSet.ValidModulesString())
	}

	run.GlobalConf.RequestedModule = modSet[run.GlobalConf.Module]

	// TODO(spencer) - set module-specific flags

	if run.GlobalConf.LogFilePath != "" {
		f, err := os.OpenFile(run.GlobalConf.LogFilePath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			logger.Fatalf("Unable to open log file (%s): %s", run.GlobalConf.LogFilePath, err.Error())
		}
		log.SetOutput(f)
	}

	// Translate the assigned verbosity level to a logrus log level.
	switch run.GlobalConf.Verbosity {
	case 1: // Fatal
		log.SetLevel(log.FatalLevel)
	case 2: // Error
		log.SetLevel(log.ErrorLevel)
	case 3: // Warnings  (default)
		log.SetLevel(log.WarnLevel)
	case 4: // Information
		log.SetLevel(log.InfoLevel)
	case 5: // Debugging
		log.SetLevel(log.DebugLevel)
	default:
		logger.Fatal("Unknown verbosity level specified. Must be between 1 (lowest)--5 (highest)")
	}

	// complete post facto global initialization based on command line arguments
	run.GlobalConf.Timeout = time.Duration(time.Second * time.Duration(run.Timeout))
	run.GlobalConf.IterationTimeout = time.Duration(time.Second * time.Duration(run.IterationTimeout))

	// class initialization
	switch strings.ToUpper(run.Class) {
	case "INET", "IN":
		run.GlobalConf.Class = dns.ClassINET
	case "CSNET", "CS":
		run.GlobalConf.Class = dns.ClassCSNET
	case "CHAOS", "CH":
		run.GlobalConf.Class = dns.ClassCHAOS
	case "HESIOD", "HS":
		run.GlobalConf.Class = dns.ClassHESIOD
	case "NONE":
		run.GlobalConf.Class = dns.ClassNONE
	case "ANY":
		run.GlobalConf.Class = dns.ClassANY
	default:
		logger.Fatal("Unknown record class specified. Valid valued are INET (default), CSNET, CHAOS, HESIOD, NONE, ANY")
	}

	if run.Servers == "" {
		// if we're doing recursive resolution, figure out default OS name servers
		// otherwise, use the set of 13 root name servers
		if run.GlobalConf.IterativeResolution {
			run.GlobalConf.NameServers = RootServers[:]
		} else {
			ns, err := zdns.GetDNSServers(run.ConfigFile)
			if err != nil {
				ns = util.GetDefaultResolvers()
				logger.Warn("Unable to parse resolvers file. Using ZDNS defaults: ", strings.Join(ns, ", "))
			}
			run.GlobalConf.NameServers = ns
		}
		run.GlobalConf.NameServersSpecified = false
		logger.Info("No name servers specified. will use: ", strings.Join(run.GlobalConf.NameServers, ", "))
	} else {
		if run.GlobalConf.NameServerMode {
			logger.Fatal("name servers cannot be specified on command line in --name-server-mode")
		}
		var ns []string
		if (run.Servers)[0] == '@' {
			filepath := (run.Servers)[1:]
			f, err := ioutil.ReadFile(filepath)
			if err != nil {
				logger.Fatalf("Unable to read file (%s): %s", filepath, err.Error())
			}
			if len(f) == 0 {
				logger.Fatalf("Empty file (%s)", filepath)
			}
			ns = strings.Split(strings.Trim(string(f), "\n"), "\n")
		} else {
			ns = strings.Split(run.Servers, ",")
		}
		for i, s := range ns {
			ns[i] = util.AddDefaultPortToDNSServerName(s)
		}
		run.GlobalConf.NameServers = ns
		run.GlobalConf.NameServersSpecified = true
	}

	if run.LocalAddr != "" {
		for _, la := range strings.Split(run.LocalAddr, ",") {
			ip := net.ParseIP(la)
			if ip != nil {
				run.GlobalConf.LocalAddrs = append(run.GlobalConf.LocalAddrs, ip)
			} else {
				logger.Fatal("Invalid argument for --local-addr (", la, "). Must be a comma-separated list of valid IP addresses.")
			}
		}
		logger.Info("using local address: ", run.LocalAddr)
		run.GlobalConf.LocalAddrSpecified = true
	}

	if run.LocalIF != "" {
		if run.GlobalConf.LocalAddrSpecified {
			logger.Fatal("Both --local-addr and --local-interface specified.")
		} else {
			li, err := net.InterfaceByName(run.LocalIF)
			if err != nil {
				logger.Fatal("Invalid local interface specified: ", err)
			}
			addrs, err := li.Addrs()
			if err != nil {
				logger.Fatal("Unable to detect addresses of local interface: ", err)
			}
			for _, la := range addrs {
				run.GlobalConf.LocalAddrs = append(run.GlobalConf.LocalAddrs, la.(*net.IPNet).IP)
				run.GlobalConf.LocalAddrSpecified = true
			}
			logger.Info("using local interface: ", run.LocalIF)
		}
	}
	if !run.GlobalConf.LocalAddrSpecified {
		// Find local address for use in unbound UDP sockets
		if conn, err := net.Dial("udp", "8.8.8.8:53"); err != nil {
			logger.Fatal("Unable to find default IP address: ", err)
		} else {
			run.GlobalConf.LocalAddrs = append(run.GlobalConf.LocalAddrs, conn.LocalAddr().(*net.UDPAddr).IP)
		}
	}
	if run.NanoSeconds {
		run.GlobalConf.TimeFormat = time.RFC3339Nano
	} else {
		run.GlobalConf.TimeFormat = time.RFC3339
	}
	if run.GlobalConf.GoMaxProcs < 0 {
		logger.Fatal("Invalid argument for --go-processes. Must be >1.")
	}
	if run.GlobalConf.GoMaxProcs != 0 {
		runtime.GOMAXPROCS(run.GlobalConf.GoMaxProcs)
	}
	if run.GlobalConf.UDPOnly && run.GlobalConf.TCPOnly {
		logger.Fatal("TCP Only and UDP Only are conflicting")
	}
	if run.GlobalConf.NameServerMode && run.GlobalConf.AlexaFormat {
		logger.Fatal("Alexa mode is incompatible with name server mode")
	}
	if run.GlobalConf.NameServerMode && run.GlobalConf.MetadataFormat {
		logger.Fatal("Metadata mode is incompatible with name server mode")
	}
	if run.GlobalConf.NameServerMode && run.GlobalConf.NameOverride == "" && run.GlobalConf.Module != "BINDVERSION" {
		logger.Fatal("Static Name must be defined with --override-name in --name-server-mode unless DNS module does not expect names (e.g., BINDVERSION).")
	}
	// Output Groups are defined by a base + any additional fields that the user wants
	groups := strings.Split(run.GlobalConf.IncludeInOutput, ",")
	if run.GlobalConf.ResultVerbosity != "short" && run.GlobalConf.ResultVerbosity != "normal" && run.GlobalConf.ResultVerbosity != "long" && run.GlobalConf.ResultVerbosity != "trace" {
		logger.Fatal("Invalid result verbosity. Options: short, normal, long, trace")
	}

	run.GlobalConf.OutputGroups = append(run.GlobalConf.OutputGroups, run.GlobalConf.ResultVerbosity)
	run.GlobalConf.OutputGroups = append(run.GlobalConf.OutputGroups, groups...)

	// some modules require multiple passes over a file (this is really just the case for zone files)
	if !run.GlobalConf.RequestedModule.Module.AllowStdIn() && run.GlobalConf.InputFilePath == "-" {
		logger.Fatal("Specified module does not allow reading from stdin")
	}

	// setup i/o
	run.GlobalConf.InputHandler = iohandlers.NewFileInputHandler(run.GlobalConf.InputFilePath)
	run.GlobalConf.OutputHandler = iohandlers.NewFileOutputHandler(run.GlobalConf.OutputFilePath)

	err := RunLookups(&run.GlobalConf)

	if err != nil {
		logger.Fatal(err)
	}
}
