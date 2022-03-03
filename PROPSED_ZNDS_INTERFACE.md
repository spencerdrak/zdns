# Proposed ZDNS Interface 

The interface defined below exposes the ZDNS library to the modules and other programs that wish to use it. This also assumes that the caching interface will remain nearly the same, so those interfaces have been omitted for brevity.

The idea here is that the modules will need to be updated to use this new interface, but will _also_ provide the same interface to any client programs. In this way, we can allow the modules to remain lightweight and only require the end-user client programs to manage the state around input/output, caching, and concurrency. The library will be responsible for the raw lookups, and an example of a module might be something the MXLOOKUP module, which performs more complex operations on top of this. 

This allows clients to either implement all their logic in their client program, or, if their logic is more generally applicable, to write a module that meets their needs.

Finally, The global and routine LookupFactories are not necessary anymore. I believe that we are moving that specific logic to the modules or clients themselves - so they'll be responsible for managing all of that state.


```go
type IsTraced bool
type ModuleOptions map[string]string

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
    // Set an ID to associate distinct queries together, for easier aggregation
    // ID will be passed along through the answer.
    Id          UUID
    // Timeout for individual name resolution
    Timeout     int
}

type ClientOptions struct {
    // Reuse socket between requests
    ReuseSockets bool
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
    // The nameserver to use for all lookups on this interface. Leave blank for default (system?) resolver
    Nameserver  string
    // Allow modules to specify their own options if needed.
    // Modules will be responsible for parsing/validating these options.
    // The raw ZDNS lookups will leave this empty
    ModuleOptions ModuleOptions
}

type Cache struct {
	IterativeCache cachehash.ShardedCacheHash
}

type LookupClient interface {
	Initialize(options ClientOptions) error
    SetOptions(options ClientOptions) error
	DoLookup(question Question) (Response, error)
	DoInternallyRecurisveLookup(question Question, cache Cache) (Response, error)
}
```