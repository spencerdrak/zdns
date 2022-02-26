# Proposed ZDNS Interface 

The interface defined below exposes the ZDNS library to the modules and other programs that wish to use it. This also assumes that the caching interface will remain nearly the same, so those interfaces have been omitted for brevity.

The idea here is that the modules will need to be updated to use this new interface, but will _also_ provide the same interface to any client programs. In this way, we can allow the modules to remain lightweight and only require the end-user client programs to manage the state around input/output, caching, and concurrency. The library will be responsible for the raw lookups, and an example of a module might be something the MXLOOKUP module, which performs more complex operations on top of this. 

This allows clients to either implement all their logic in their client program, or, if their logic is more generally applicable, to write a module that meets their needs.

Finally, The global and routine LookupFactories may not be necessary anymore. I believe that we are moving that specific logic to the modules or clients themselves - so they'll be responsible for managing all of that state.


```go
type IsCached bool
type IsTraced bool

type Response struct {
    Result zdns.Result
    Trace  zdns.Trace
    Status zdns.Status
    // Return an ID linked to the Question, so that distinct queries can be linked.
    Id     UUID
}

type Question struct {
    // The DNS type to query for
	Type        uint16
    // The class to query for
	Class       uint16
    // The Domain name in question
	Name        string
    // The nameserver to use. Leave blank for default (system?) resolver
    Nameserver  string
    // Set an ID to associate distinct queries together, for easier aggregation
    // ID will be passed along through the answer.
    Id          UUID
    // Timeout for individual name resolution
    Timeout     int
}

type ClientOptions struct {
    // Reuse socket between requests
    ReuseSockets bool
    // Pass in a cache, shared between threads
    Cache        Cache
    // Use the above cache
    IsCached     IsCached
    // Return a trace
    IsTraced     IsTraced
    // Logging Verbosity
    Verbosity    int
    // Max depth of recursion. Only useful for iterative lookup
    MaxDepth     int
    TCPOnly      bool
    UDPOnly      bool
    // Nanosecond timestamp resolution in output
    NsResolution bool
    // Local Address to use for requests
    LocalAddr    net.IP
    // Local interface to use for requests
    LocalIF      net.Interface
}

type Cache struct {
	IterativeCache cachehash.ShardedCacheHash
}

type LookupClient interface {
	Initialize(options ClientOptions) error
    SetOptions(options ClientOptions) error
	DoLookup() error
	DoIterativeLookup(cache Cache) error
}
```

The Module interface will also be standardized and made to be more like ZGrab2. See below for a propsed interface:

```go

// Module is an interface that represents all functions necessary to run a lookup
type Module interface {
	// Init runs once for this module at library init time
	Init(flags ScanFlags) error

	// Returns the name passed at init
	GetName() string

	// Returns the trigger passed at init
	GetTrigger() string

	// Protocol returns the protocol identifier for the scan.
	Protocol() string

	// Scan connects to a host. The result should be JSON-serializable
	Scan(t ScanTarget) (ScanStatus, interface{}, error)
}

// ScanModule is an interface which represents a module that the framework can
// manipulate
type ScanModule interface {
	// NewFlags is called by the framework to pass to the argument parser. The parsed flags will be passed
	// to the scanner created by NewScanner().
	NewFlags() interface{}

	// NewScanner is called by the framework for each time an individual scan is specified in the config or on
	// the command-line. The framework will then call scanner.Init(name, flags).
	NewScanner() Scanner

	// Description returns a string suitable for use as an overview of this
	// module within usage text.
	Description() string
}

// ScanFlags is an interface which must be implemented by all types sent to
// the flag parser
type ScanFlags interface {
	// Help optionally returns any additional help text, e.g. specifying what empty defaults
	// are interpreted as.
	Help() string

	// Validate enforces all command-line flags and positional arguments have valid values.
	Validate(args []string) error
}

// BaseFlags contains the options that every flags type must embed
type BaseFlags struct {
	Port           uint          `short:"p" long:"port" description:"Specify port to grab on"`
	Name           string        `short:"n" long:"name" description:"Specify name for output json, only necessary if scanning multiple modules"`
	Timeout        time.Duration `short:"t" long:"timeout" description:"Set connection timeout (0 = no timeout)" default:"10s"`
	Trigger        string        `short:"g" long:"trigger" description:"Invoke only on targets with specified tag"`
	BytesReadLimit int           `short:"m" long:"maxbytes" description:"Maximum byte read limit per scan (0 = defaults)"`
}

// UDPFlags contains the common options used for all UDP scans
type UDPFlags struct {
	LocalPort    uint   `long:"local-port" description:"Set an explicit local port for UDP traffic"`
	LocalAddress string `long:"local-addr" description:"Set an explicit local address for UDP traffic"`
}

// GetName returns the name of the respective scanner
func (b *BaseFlags) GetName() string {
	return b.Name
}

// GetModule returns the registered module that corresponds to the given name
// or nil otherwise
func GetModule(name string) ScanModule {
	return modules[name]
}

var modules map[string]ScanModule

```