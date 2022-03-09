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
	"errors"
	"fmt"
	"time"

	"github.com/zmap/dns"
	"github.com/zmap/zdns/internal/util"

	log "github.com/sirupsen/logrus"
)

// Keep track of some state internal to ZDNS/Raw Module
type RawOptions struct {
	IterativeStop time.Time
	Logger        *log.Entry
}

// Provide a LookupClient for the "raw" modules, e.g., the ZDNS library itself.
type RawLookupClient struct {
	ClientOptions
	RawOptions
}

// Create the module wrapper around the RawLookupClient
type RawModule struct{}

type ConfigError struct {
	Field string
	Msg   string
}

type TraceStep struct {
	Result     Result                `json:"results" groups:"trace"`
	DnsType    uint16                `json:"type" groups:"trace"`
	DnsClass   uint16                `json:"class" groups:"trace"`
	Name       string                `json:"name" groups:"trace"`
	NameServer string                `json:"name_server" groups:"trace"`
	Depth      int                   `json:"depth" groups:"trace"`
	Layer      string                `json:"layer" groups:"trace"`
	Cached     IsInternallyRecursive `json:"cached" groups:"trace"`
}

type DNSFlags struct {
	Response           bool `json:"response" groups:"flags,long,trace"`
	Opcode             int  `json:"opcode" groups:"flags,long,trace"`
	Authoritative      bool `json:"authoritative" groups:"flags,long,trace"`
	Truncated          bool `json:"truncated" groups:"flags,long,trace"`
	RecursionDesired   bool `json:"recursion_desired" groups:"flags,long,trace"`
	RecursionAvailable bool `json:"recursion_available" groups:"flags,long,trace"`
	Authenticated      bool `json:"authenticated" groups:"flags,long,trace"`
	CheckingDisabled   bool `json:"checking_disabled" groups:"flags,long,trace"`
	ErrorCode          int  `json:"error_code" groups:"flags,long,trace"`
}

type RawResult struct {
	Answers     []interface{} `json:"answers,omitempty" groups:"short,normal,long,trace"`
	Additional  []interface{} `json:"additionals,omitempty" groups:"short,normal,long,trace"`
	Authorities []interface{} `json:"authorities,omitempty" groups:"short,normal,long,trace"`
	Protocol    string        `json:"protocol" groups:"protocol,normal,long,trace"`
	Resolver    string        `json:"resolver" groups:"resolver,normal,long,trace"`
	Flags       DNSFlags      `json:"flags" groups:"flags,long,trace"`
}

// TODO(spencer): use logging package, but my solution is a half-and-half that is bad
func (lc RawLookupClient) VerboseLog(depth int, args ...interface{}) {
	lc.RawOptions.Logger.Debug(util.MakeVerbosePrefix(depth), args)
}

func (e ConfigError) Error() string {
	return fmt.Sprintf("Invalid ZDNS Config in field %s - %s", e.Field, e.Msg)
}

func (m RawModule) NewLookupClient() LookupClient {
	return RawLookupClient{}
}

func (lc RawLookupClient) Initialize(option ClientOptions) error {
	// Do args validation on input
	// set fields on RawLookupClient

	lc.RawOptions.Logger = log.WithFields(log.Fields{
		"Module": "RawLookupClient",
	})

	return errors.New("not implemented")
}

func (lc RawLookupClient) SetOptions(options ClientOptions) error {
	// Do args validation on input
	// set fields on RawLookupClient
	return errors.New("not implemented")
}

func (lc RawLookupClient) DoLookup(question Question) (Response, error) {
	if question.Type == 0 {
		return Response{}, ConfigError{"Type", "unset (set to 0)"}
	}
	if question.Class == 0 {
		return Response{}, ConfigError{"Class", "unset (set to 0)"}
	}
	if question.Type == dns.TypePTR {
		var err error
		question.Name, err = dns.ReverseAddr(question.Name)
		if err != nil {
			resp := Response{
				Result: Result{},
				Trace:  Trace{},
				Status: STATUS_ILLEGAL_INPUT,
				Id:     question.Id,
			}
			return resp, err
		}
		question.Name = question.Name[:len(question.Name)-1]
	}
	if lc.ClientOptions.IsInternallyRecursive {
		lc.VerboseLog(0, "MIEKG-IN: iterative lookup for ", question.Name, " (", question.Type, ")")
		lc.RawOptions.IterativeStop = time.Now().Add(time.Duration(lc.ClientOptions.IterativeOptions.IterativeTimeout))
		response, err := lc.iterativeLookup(question, lc.ClientOptions.Nameserver, 1, ".", make([]interface{}, 0))
		lc.VerboseLog(0, "MIEKG-OUT: iterative lookup for ", question.Name, " (", question.Type, "): status: ", response.Status, " , err: ", err)

		// TODO(spencer): confirm tracing behavior
		if lc.ClientOptions.IsTraced {
			return response, err
		}
		// TODO(spencer): unsure if the if-block above does anything.
		return response, err
	} else {
		return tracedRetryingLookup(question, lc.ClientOptions.Nameserver, true)
	}
}

func (lc RawLookupClient) iterativeLookup(question Question, nameServer string, depth int, layer string, trace []interface{}) (Response, error) {
	if log.GetLevel() == log.DebugLevel {
		lc.VerboseLog((depth), "iterative lookup for ", q.Name, " (", q.Type, ") against ", nameServer, " layer ", layer)
	}
	if depth > lc.ClientOptions.MaxDepth {
		lc.VerboseLog((depth + 1), "-> Max recursion depth reached")
		return Response{Result{}, trace, STATUS_ERROR, question.Id, nil}, errors.New("Max recursion depth reached")
	}
	response, err := lc.cachedRetryingLookup(question, nameServer, layer, depth)
	if lc.IsTraced && response.Status == STATUS_NOERROR {
		var t TraceStep
		t.Result = response.Result
		t.DnsType = question.Type
		t.DnsClass = question.Class
		t.Name = question.Name
		t.NameServer = nameServer
		t.Layer = layer
		t.Depth = depth
		t.Cached = lc.IsInternallyRecursive
		trace = append(trace, t)

	}
	if response.Status != STATUS_NOERROR {
		lc.VerboseLog((depth + 1), "-> error occurred during lookup")
		return response, err
	} else if len(response.Result.Answers) != 0 || result.Flags.Authoritative == true {
		if len(result.Answers) != 0 {
			s.VerboseLog((depth + 1), "-> answers found")
			if len(result.Authorities) > 0 {
				s.VerboseLog((depth + 2), "Dropping ", len(result.Authorities), " authority answers from output")
				result.Authorities = make([]interface{}, 0)
			}
			if len(result.Additional) > 0 {
				s.VerboseLog((depth + 2), "Dropping ", len(result.Additional), " additional answers from output")
				result.Additional = make([]interface{}, 0)
			}
		} else {
			s.VerboseLog((depth + 1), "-> authoritative response found")
		}
		return result, trace, status, err
	} else if len(result.Authorities) != 0 {
		s.VerboseLog((depth + 1), "-> Authority found, iterating")
		return s.iterateOnAuthorities(q, depth, result, layer, trace)
	} else {
		s.VerboseLog((depth + 1), "-> No Authority found, error")
		return result, trace, zdns.STATUS_ERROR, errors.New("NOERROR record without any answers or authorities")
	}
}

func tracedRetryingLookup(question Question, nameServer string, recursive bool) (Response, error) {
	return Response{Result{}, nil, STATUS_ERROR, question.Id, nil}, errors.New("not implemented")
}

func (lc RawLookupClient) cachedRetryingLookup(question Question, nameServer, layer string, depth int) (Response, error) {
	return Response{Result{}, nil, STATUS_ERROR, question.Id, nil}, errors.New("not implemented")
}
