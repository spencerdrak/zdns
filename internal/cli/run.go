package cli

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/liip/sheriff"
	log "github.com/sirupsen/logrus"
	"github.com/zmap/dns"
	"github.com/zmap/zdns/internal/util"
	"github.com/zmap/zdns/pkg/zdns"
)

type routineMetadata struct {
	Names  int
	Status map[zdns.Status]int
}

func aggregateMetadata(c <-chan routineMetadata) Metadata {
	var meta Metadata
	meta.Status = make(map[string]int)
	for m := range c {
		meta.Names += m.Names
		for k, v := range m.Status {
			meta.Status[string(k)] += v
		}
	}
	return meta
}

func RunLookups(c *GlobalConf) error {

	logger := log.WithFields(log.Fields{
		"Module": "cli",
	})

	// DoLookup:
	//	- n threads that do processing from in and place results in out
	//	- process until inChan closes, then wg.done()
	// Once we processing threads have all finished, wait until the
	// output and metadata threads have completed
	inChan := make(chan interface{})
	outChan := make(chan string)
	metaChan := make(chan routineMetadata, c.Threads)
	var routineWG sync.WaitGroup

	inHandler := c.InputHandler
	if inHandler == nil {
		logger.Fatal("Input handler is nil")
	}

	outHandler := c.OutputHandler
	if outHandler == nil {
		logger.Fatal("Output handler is nil")
	}

	// Use handlers to populate the input and output/results channel
	go inHandler.FeedChannel(inChan, &routineWG)
	go outHandler.WriteResults(outChan, &routineWG)
	routineWG.Add(2)

	conn, localAddr, err := c.RequestedModule.NewReusableUDPConn(nil)

	if err != nil {
		panic(err)
	}

	client := c.RequestedModule.NewLookupClient()

	//TODO(spencer) - populate from GlobalConf
	clientOptions := zdns.ClientOptions{
		ReuseSockets:          false,
		IsTraced:              true,
		Verbosity:             3,
		TCPOnly:               false,
		UDPOnly:               false,
		NsResolution:          false,
		LocalAddr:             localAddr,
		Conn:                  &conn,
		Nameserver:            "1.1.1.1:53",
		ModuleOptions:         map[string]string{},
		IsInternallyRecursive: false,
		IterativeOptions:      zdns.IterativeOptions{},
	}

	err = client.Initialize(&clientOptions)

	if err != nil {
		logger.Fatal(err)
	}

	// create pool of worker goroutines
	var lookupWG sync.WaitGroup
	lookupWG.Add(c.Threads)
	startTime := time.Now().Format(c.TimeFormat)
	for i := 0; i < c.Threads; i++ {

	}
	lookupWG.Wait()
	close(outChan)
	close(metaChan)
	routineWG.Wait()
	if c.MetadataFilePath != "" {
		// we're done processing data. aggregate all the data from individual routines
		metaData := aggregateMetadata(metaChan)
		metaData.StartTime = startTime
		metaData.EndTime = time.Now().Format(c.TimeFormat)
		metaData.NameServers = c.NameServers
		metaData.Retries = c.Retries
		// Seconds() returns a float. However, timeout is passed in as an integer
		// command line argument, so there should be no loss of data when casting
		// back to an integer here.
		metaData.Timeout = int(c.Timeout.Seconds())
		metaData.Conf = c
		// add global lookup-related metadata
		// write out metadata
		var f *os.File
		if c.MetadataFilePath == "-" {
			f = os.Stderr
		} else {
			var err error
			f, err = os.OpenFile(c.MetadataFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				logger.Fatal("unable to open metadata file:", err.Error())
			}
			defer f.Close()
		}
		j, err := json.Marshal(metaData)
		if err != nil {
			logger.Fatal("unable to JSON encode metadata:", err.Error())
		}
		f.WriteString(string(j))
	}
	return nil
}

func runRoutineLookup(gc *GlobalConf, input <-chan interface{}, output chan<- string, metaChan chan<- routineMetadata, wg *sync.WaitGroup, lc zdns.LookupClient) error {
	logger := log.WithFields(log.Fields{
		"Module": "cli",
	})

	var metadata routineMetadata
	metadata.Status = make(map[zdns.Status]int)
	for genericInput := range input {
		var res zdns.Response
		var innerRes interface{}
		var trace []interface{}
		var status zdns.Status
		var err error
		l, err := f.MakeLookup()
		if err != nil {
			logger.Fatal("Unable to build lookup instance", err)
		}
		line := genericInput.(string)
		var changed bool
		var lookupName string
		rawName := ""
		nameServer := ""
		var rank int
		var entryMetadata string
		if gc.AlexaFormat == true {
			rawName, rank = parseAlexa(line)
			res.AlexaRank = rank
		} else if gc.MetadataFormat {
			rawName, entryMetadata = parseMetadataInputLine(line)
			res.Metadata = entryMetadata
		} else if gc.NameServerMode {
			nameServer = util.AddDefaultPortToDNSServerName(line)
		} else {
			rawName, nameServer = parseNormalInputLine(line)
		}
		lookupName, changed = makeName(rawName, gc.NamePrefix, gc.NameOverride)
		if changed {
			res.AlteredName = lookupName
		}
		res.Name = rawName
		res.Class = dns.Class(gc.Class).String()
		innerRes, trace, status, err = l.DoLookup(lookupName, nameServer)
		res.Timestamp = time.Now().Format(gc.TimeFormat)
		if status != STATUS_NO_OUTPUT {
			res.Status = string(status)
			res.Data = innerRes
			res.Trace = trace
			if err != nil {
				res.Error = err.Error()
			}
			v, _ := version.NewVersion("0.0.0")
			o := &sheriff.Options{
				Groups:     gc.OutputGroups,
				ApiVersion: v,
			}
			data, err := sheriff.Marshal(o, res)
			jsonRes, err := json.Marshal(data)
			if err != nil {
				logger.Fatal("Unable to marshal JSON result", err)
			}
			output <- string(jsonRes)
		}
		metadata.Names++
		metadata.Status[status]++
	}
	metaChan <- metadata
	wg.Done()
	return nil
}
