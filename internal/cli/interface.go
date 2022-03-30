package cli

import (
	"sync"
)

// handle domain input
type InputHandler interface {
	// FeedChannel takes a channel to write domains to, the WaitGroup managing them, and if it's a zonefile input
	FeedChannel(in chan<- interface{}, wg *sync.WaitGroup) error
}

// handle output results
type OutputHandler interface {
	// takes a channel (results) to write the query results to, and the WaitGroup managing the handlers
	WriteResults(results <-chan string, wg *sync.WaitGroup) error
}
