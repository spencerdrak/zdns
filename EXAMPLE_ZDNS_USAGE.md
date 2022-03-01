## Example ZDNS Usage

This is a sample program to show how a client or driver program might make use of the new ZDNS interface.

```go
package main

import (
    "github.com/zmap/zdns"
    "github.com/zmap/zdns/modules/mxlookup"
    "github.com/google/uuid"
)

// You'll likely want to replace this with some other type of datastructure in a real application, but for a sample, this works.
hosts := [2]string{"censys.io", "google.com"}

rawAnswers = make([]zdns.Response, 0)
mxAnswers = make([]mxlookup.Response, 0)

options := ClientOptions{
    ReuseSockets: true,
    IsTraced: true,
    Verbosity: 3,
    MaxDepth: 10,
    TCPOnly: false
    UDPOnly: false
    NsResolution: false
    LocalAddr: nil
    LocalIF: nil
    Nameserver: "1.1.1.1"
}

rawClient := zdns.NewLookupClient()
rawClient.Initialize(options)

mxClient := mxlookup.NewLookupClient()
mxClient.Initialize(options)

for _, host := range hosts {
    q := {
        Type: 1
        Class: 1
        Name: host,
        Id: uuid.New()
        Timeout: 15
    }
    rawAnswers.append(rawAnswers, rawClient.DoLookup())
}

for _, host := range hosts {
    q := {
        Type: 1
        Class: 1
        Name: host,
        Id: uuid.New()
        Timeout: 15
    }
    mxAnswers.append(mxAnswers, mxLookup.DoLookup())
}
```