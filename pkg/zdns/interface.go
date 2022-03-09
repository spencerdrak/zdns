/*
 * ZDNS Copyright 2016 Regents of the University of Michigan
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not
 * use this file except in compliance with the License. You may obtain a copy
 * of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
 * implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

package zdns

import (
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/zmap/dns"
	"github.com/zmap/go-iptree/blacklist"
)

/* Each lookup module registers a single GlobalLookupFactory, which is
 * instantiated once.  This global factory is responsible for providing command
 * line arguments and performing any configuration that should only occur once.
 * For each thread in the worker pool, the framework calls
 * MakePerRoutineFactory(), which should return a second factory, which should
 * perform any "per-thread" initialization. Within each "thread", the framework
 * will then call MakeLookup() for each connection it will make, on which it
 * will call DoLookup().  While two layers of factories is a bit... obnoxious,
 * this allows each module to maintain global, per-thread, and per-connection
 * state.
 *
 * Each layer has access to one proceeding layer (e.g., RoutineLookupFactory
 * knows the GlobalLookupFactory, from which it was created. Therefore, modules
 * should refer to this configuration instead of copying all configuration
 * values for every connection. The Base structs implement these basic
 * pieces of functionality and should be inherited in most situations.
 */

type Trace []interface{}

type Module interface {
	// NewLookupClient is called by the client to get a new LookupModule
	NewLookupClient() LookupClient
}

type IsTraced bool
type IsInternallyRecursive bool
type ModuleOptions map[string]string

type Response struct {
	//TODO(spencer): revisit Result handling
	Result interface{}
	Trace  Trace
	Status Status
	// Return an ID linked to the Question, so that distinct queries can be linked.
	Id uuid.UUID
	// Define an additional field such that modules can return extra data as needed.
	Additional interface{}
}

type Question struct {
	// The DNS type to query for
	// TODO(spencer): make this a dns.Type?
	Type uint16
	// The class to query for
	// TODO(spencer): make this a dns.Class?
	Class uint16
	// The Domain name in question
	Name string
	// Set an ID to associate distinct queries together, for easier aggregation
	// ID will be passed along through the answer.
	Id uuid.UUID
	// Timeout for individual name resolution
	Timeout int
}

type IterativeOptions struct {
	// Cache to use if the IsInternallyRecursiveFlag is set
	Cache               Cache
	IterativeTimeout    time.Duration
	IterativeResolution bool
	// Max depth of recursion. Only useful for iterative lookup
	MaxDepth int
}

type ClientOptions struct {
	// Reuse socket between requests
	ReuseSockets bool
	// Return a trace
	IsTraced
	// Logging Verbosity
	Verbosity int
	TCPOnly   bool
	UDPOnly   bool
	// Nanosecond timestamp resolution in output
	NsResolution bool
	// Local Address to use for requests
	LocalAddr net.IP
	// Local interface to use for requests
	LocalIF net.Interface
	// Nameserver to use if not internally recursive
	Nameserver string
	// How many times to retry a lookup
	Retries int
	// Connection to use for lookups
	Conn *dns.Conn
	// Set a blacklist of nameservers to not use
	Blacklist *blacklist.Blacklist
	// Protect this blacklist from concurrent access
	BlackListMutex sync.Mutex
	// Allow modules to specify their own options if needed.
	// Modules will be responsible for parsing/validating these options.
	// The raw ZDNS lookups will leave this empty
	ModuleOptions
	// IsInternallyRecursive tells DoLookup to do internal recursion. If true, uses cache. If false, uses nameserver.
	IsInternallyRecursive
	IterativeOptions
}

type LookupClient interface {
	Initialize(options ClientOptions) error
	SetOptions(options ClientOptions) error
	DoLookup(question Question) (Response, error)
}
