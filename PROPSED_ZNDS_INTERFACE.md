# Proposed ZDNS Interface 

The interface defined below exposes the ZDNS library to the modules and other programs that wish to use it.

```go
type IsCached bool
type QuestionId 

type Response struct {
    Result zdns.Result
    Trace  zdns.Trace
    Status zdns.Status
    Id     UUID
}

type Question struct {
	Type        uint16
	Class       uint16
	Name        string
    Nameserver  string
    // Set an ID to associate distinct queries together, for easier aggregation
    Id          UUID
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
	DoRecursiveLookup(input InputHandler, oh OutputHandler, options LookupOptions) error
	DoIterativeLookup(input InputHandler, oh OutputHandler, options LookupOptions) error
    RandomNameServer() string
}

type LookupOptions struct {
    Traced   bool
    Cached   bool
    Retrying bool
    Cache    Cache
    IsCached IsCached
}

type ClientOptions struct {
	NameServer   string
	DnsType      uint16
	DnsClass     uint16
    ReuseSockets bool
}
```