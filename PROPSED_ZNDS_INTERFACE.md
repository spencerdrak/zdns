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
	Initialize(name string, options ClientOptions) error
    SetOptions(options ClientOptions) error
	DoLookup() error
	DoIterativeLookup(cache Cache) error
}
```

The Module interface will also be standardized and made to be more like ZGrab2. See below for a proposed interface:

```go

// LookupModule is an interface which represents some higher-level functionality above a 
type LookupModule interface {
	// NewLookupClient is called by the framework for each time an individual scan is specified in the config or on
	// the command-line. The framework will then call scanner.Initialize(name, flags).
	NewLookupClient() zdns.LookupClient

	// Description returns a string suitable for use as an overview of this
	// module within usage text.
	Description() string
}

var modules ModuleSet

```

The ModuleSet interface (almost entirely borrowed from ZGrab2) is below:

```go
// ModuleSet is a map of name (string) -> Module, one per module.
type ModuleSet map[string]LookupModule

// CopyInto copies the modules in s to destination. The sets will be unique, but
// the underlying ScanModule instances will be the same.
func (s ModuleSet) CopyInto(destination ModuleSet) {
	for name, m := range s {
		if _, ok := destination[strings.ToUpper(name)]; ok {
			log.Warnf("overwriting module %s", name)
		}
		destination[strings.ToUpper(name)] = m
	}
}

// AddModule adds m to the ModuleSet, accessible via the given name. If the name
// is already in the ModuleSet, it is overwritten.
func (s ModuleSet) AddModule(name string, m GlobalLookupFactory) {
	if _, ok := s[strings.ToUpper(name)]; ok {
		log.Warnf("overwriting module %s", name)
	}
	s[strings.ToUpper(name)] = m
}

// RemoveModule removes the module at the specified name. If the name is not in
// the module set, nothing happens.
func (s ModuleSet) RemoveModule(name string) {
	delete(s, strings.ToUpper(name))
}

// ModuleSet returns an empty ModuleSet.
func NewModuleSet() ModuleSet {
	return make(ModuleSet)
}
```