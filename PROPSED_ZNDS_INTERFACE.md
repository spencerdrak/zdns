# Proposed ZDNS Interface 

The interface defined below exposes the ZDNS library to the modules and other programs that wish to use it.

This also assumes that the caching interface will remain nearly the same, so those interfaces have been omitted for brevity.

There will be some additional structure around how modules are created and registered in the

Finally, The global and routine LookupFactories may not be necessary anymore. I believe that we 

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
}

type Cache struct {
	IterativeCache cachehash.ShardedCacheHash
}

type InputHandler interface {
	FeedChannel(in chan<- Question, wg *sync.WaitGroup) error
}

type OutputHandler interface {
	WriteResults(results <-chan Response, wg *sync.WaitGroup) error
}

type LookupClient interface {
	Initialize(options ClientOptions) error
    SetOptions(options ClientOptions) error
    // maybe put the input/output handlers in the client, not per-query
	DoLookup() error
    //TODO: Caching is only useful in this case
	DoIterativeLookup(input InputHandler, oh OutputHandler, options LookupOptions) error
}
```