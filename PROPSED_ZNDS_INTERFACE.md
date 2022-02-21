# Proposed ZDNS Interface 

The interface defined below exposes the ZDNS library to the modules and other programs that wish to use it.

```go
// LookupClient configured once per module or use case.
type LookupClient interface {
	Initialize(options ClientOptions) error
    SetOptions(options ClientOptions) error
	DoRecursiveLookup(options LookupOptions) (zdns.Result, zdns.Trace, zdns.Status, error)
	DoIterativeLookup(options LookupOptions) (zdns.Result, zdns.Trace, zdns.Status, error)

    Help() string

}

type LookupOptions struct {
    Traced   bool
    Cached   bool
    Retrying bool
}

type ClientOptions struct {
	NameServer   string
	DnsType      uint16
	DnsClass     uint16
	Factory      zdns.RoutineLookupFactory
    ReuseSockets bool
}
```