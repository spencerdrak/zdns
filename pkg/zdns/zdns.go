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
	"net"
	"strings"
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
	*ClientOptions
	RawOptions
	Client    *dns.Client
	TCPClient *dns.Client
}

// Create the module wrapper around the RawLookupClient
type RawModule struct{}

type RawResult struct {
	Answers     []interface{} `json:"answers,omitempty" groups:"short,normal,long,trace"`
	Additional  []interface{} `json:"additionals,omitempty" groups:"short,normal,long,trace"`
	Authorities []interface{} `json:"authorities,omitempty" groups:"short,normal,long,trace"`
	Protocol    string        `json:"protocol" groups:"protocol,normal,long,trace"`
	Resolver    string        `json:"resolver" groups:"resolver,normal,long,trace"`
	Flags       DNSFlags      `json:"flags" groups:"flags,long,trace"`
}

func (lc *RawLookupClient) VerboseLog(depth int, args ...interface{}) {
	lc.RawOptions.Logger.Debug(util.MakeVerbosePrefix(depth), args)
}

func (e ConfigError) Error() string {
	return fmt.Sprintf("Invalid ZDNS Config in field %s - %s", e.Field, e.Msg)
}

func (m RawModule) NewLookupClient() LookupClient {
	return &RawLookupClient{}
}

func (m RawModule) NewReusableUDPConn(localAddr net.IP) (dns.Conn, net.IP, error) {
	// create PacketConn for use throughout life

	if localAddr == nil {
		// Find local address for use in unbound UDP sockets
		if conn, err := net.Dial("udp", "8.8.8.8:53"); err != nil {
			log.Fatal("Unable to find default IP address: ", err)
		} else {
			localAddr = conn.LocalAddr().(*net.UDPAddr).IP
		}
	}

	conn, err := net.ListenUDP("udp", &net.UDPAddr{localAddr, 0, ""})
	if err != nil {
		log.Fatal("unable to create socket", err)
		return dns.Conn{}, net.IP{}, err
	}
	udpConn := new(dns.Conn)
	udpConn.Conn = conn

	return *udpConn, localAddr, nil
}

func (m RawModule) NewSingleUseUDPConn(localAddr net.IP, remoteAddr net.IP) (dns.Conn, net.IP, error) {
	// create PacketConn for use throughout life

	if localAddr == nil {
		// Find local address for use in unbound UDP sockets
		if conn, err := net.Dial("udp", "8.8.8.8:53"); err != nil {
			log.Fatal("Unable to find default IP address: ", err)
		} else {
			localAddr = conn.LocalAddr().(*net.UDPAddr).IP
		}
	}

	if remoteAddr == nil {
		return dns.Conn{}, net.IP{}, errors.New("remoteAddr was set to nil")
	}

	conn, err := net.Dial("udp", remoteAddr.String())
	if err != nil {
		log.Fatal("unable to create connection", err)
		return dns.Conn{}, net.IP{}, err
	}
	udpConn := new(dns.Conn)
	udpConn.Conn = conn

	return *udpConn, localAddr, nil
}

func (m RawModule) AllowStdIn() bool {
	return true
}

func (lc *RawLookupClient) Initialize(options *ClientOptions) error {
	// Do args validation on input
	// set fields on RawLookupClient

	lc.RawOptions.Logger = log.WithFields(log.Fields{
		"Module": "RawLookupClient",
	})

	lc.ClientOptions = options

	//TODO(spencer): handle timeouts

	if lc.LocalAddr == nil {
		return ConfigError{"LocalAddr", "unset (set to nil)"}
	}

	if !lc.TCPOnly {
		lc.Client = new(dns.Client)
		lc.Client.Timeout = lc.Timeout
		lc.Client.Dialer = &net.Dialer{
			Timeout:   lc.Timeout,
			LocalAddr: &net.UDPAddr{IP: lc.LocalAddr},
		}
	}
	if !lc.UDPOnly {
		lc.TCPClient = new(dns.Client)
		lc.TCPClient.Net = "tcp"
		lc.TCPClient.Timeout = lc.Timeout
		lc.TCPClient.Dialer = &net.Dialer{
			Timeout:   lc.Timeout,
			LocalAddr: &net.TCPAddr{IP: lc.LocalAddr},
		}
	}

	lc.IterativeTimeout = options.Timeout
	lc.Retries = options.Retries
	lc.MaxDepth = options.MaxDepth
	lc.IterativeResolution = options.IterativeResolution

	return nil
}

func (lc *RawLookupClient) SetOptions(options *ClientOptions) error {
	// Do args validation on input
	// set fields on RawLookupClient
	return errors.New("not implemented")
}

func (lc *RawLookupClient) DoLookup(question Question) (Response, error) {
	if question.Type == 0 {
		return Response{}, ConfigError{"Type", "unset (set to 0)"}
	}
	if question.Class == 0 {
		return Response{}, ConfigError{"Class", "unset (set to 0)"}
	}

	timeFormat := time.RFC3339

	if lc.NsResolution {
		timeFormat = time.RFC3339Nano
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

	if lc.IsInternallyRecursive {
		lc.VerboseLog(0, "MIEKG-IN: iterative lookup for ", question.Name, " (", question.Type, ")")
		lc.RawOptions.IterativeStop = time.Now().Add(time.Duration(lc.IterativeOptions.IterativeTimeout))
		result, trace, status, err := lc.iterativeLookup(question, 1, ".", make([]interface{}, 0))
		lc.VerboseLog(0, "MIEKG-OUT: iterative lookup for ", question.Name, " (", question.Type, "): status: ", status, " , err: ", err)

		return Response{
			Result:     Result{Data: result},
			Timestamp:  time.Now().Format(timeFormat),
			Name:       question.Name,
			Trace:      trace,
			Status:     status,
			Id:         question.Id,
			Additional: nil}, err
	} else {
		result, trace, status, err := lc.tracedRetryingLookup(question, question.NameServer, true)
		return Response{
			Result:     Result{Data: result},
			Timestamp:  time.Now().Format(timeFormat),
			Name:       question.Name,
			Trace:      trace,
			Status:     status,
			Id:         question.Id,
			Additional: nil}, err
	}
}

func (lc *RawLookupClient) iterativeLookup(question Question, depth int, layer string, trace []interface{}) (RawResult, []interface{}, Status, error) {
	if log.GetLevel() == log.DebugLevel {
		lc.VerboseLog((depth), "iterative lookup for ", question.Name, " (", question.Type, ") against ", question.NameServer, " layer ", layer)
	}
	if depth > lc.MaxDepth {
		lc.VerboseLog((depth + 1), "-> Max recursion depth reached")
		return RawResult{}, trace, STATUS_ERROR, errors.New("max recursion depth reached")
	}
	result, isCached, status, err := lc.cachedRetryingLookup(question, layer, depth)
	if lc.IsTraced && status == STATUS_NOERROR {
		var t TraceStep
		t.RawResult = result
		t.DnsType = question.Type
		t.DnsClass = question.Class
		t.Name = question.Name
		t.NameServer = question.NameServer
		t.Layer = layer
		t.Depth = depth
		t.Cached = isCached
		trace = append(trace, t)

	}
	if status != STATUS_NOERROR {
		lc.VerboseLog((depth + 1), "-> error occurred during lookup")
		return result, trace, status, err
	} else if len(result.Answers) != 0 || result.Flags.Authoritative == true {
		if len(result.Answers) != 0 {
			lc.VerboseLog((depth + 1), "-> answers found")
			if len(result.Authorities) > 0 {
				lc.VerboseLog((depth + 2), "Dropping ", len(result.Authorities), " authority answers from output")
				result.Authorities = make([]interface{}, 0)
			}
			if len(result.Additional) > 0 {
				lc.VerboseLog((depth + 2), "Dropping ", len(result.Additional), " additional answers from output")
				result.Additional = make([]interface{}, 0)
			}
		} else {
			lc.VerboseLog((depth + 1), "-> authoritative response found")
		}
		return result, trace, status, err
	} else if len(result.Authorities) != 0 {
		lc.VerboseLog((depth + 1), "-> Authority found, iterating")
		rawRes, trace, status, err := lc.iterateOnAuthorities(question, depth, result, layer, trace)
		return rawRes, trace, status, err
	} else {
		lc.VerboseLog((depth + 1), "-> No Authority found, error")
		return result, trace, STATUS_ERROR, errors.New("NOERROR record without any answers or authorities")
	}
}

func (lc *RawLookupClient) cachedRetryingLookup(question Question, layer string, depth int) (RawResult, IsCached, Status, error) {
	var isCached IsCached
	isCached = false
	lc.VerboseLog(depth+1, "Cached retrying lookup. Name: ", question, ", Layer: ", layer, ", Nameserver: ", question.NameServer)
	if lc.RawOptions.IterativeStop.Before(time.Now()) {
		lc.VerboseLog(depth+2, "ITERATIVE_TIMEOUT ", question, ", Layer: ", layer, ", Nameserver: ", question.NameServer)
		return RawResult{}, isCached, STATUS_ITER_TIMEOUT, nil
	}
	// First, we check the answer
	cachedResult, ok := lc.Cache.GetCachedResult(question, false, depth+1, question.Id)
	if ok {
		isCached = true
		return cachedResult, isCached, STATUS_NOERROR, nil
	}

	// Stop if we hit a nameserver we don't want to hit
	if lc.Blacklist != nil {
		lc.BlackListMutex.Lock()
		if blacklisted, err := lc.Blacklist.IsBlacklisted(question.NameServer); err != nil {
			lc.BlackListMutex.Unlock()
			lc.VerboseLog(depth+2, "Blacklist error!", err)
			return RawResult{}, isCached, STATUS_ERROR, err
		} else if blacklisted {
			lc.BlackListMutex.Unlock()
			lc.VerboseLog(depth+2, "Hit blacklisted nameserver ", question.Name, ", Layer: ", layer, ", Nameserver: ", question.NameServer)
			return RawResult{}, isCached, STATUS_BLACKLIST, nil
		}
		lc.BlackListMutex.Unlock()
	}

	// Now, we check the authoritative:
	name := strings.ToLower(question.Name)
	layer = strings.ToLower(layer)
	authName, err := nextAuthority(name, layer)
	if err != nil {
		lc.VerboseLog(depth+2, err)
		return RawResult{}, isCached, STATUS_AUTHFAIL, err
	}
	if name != layer && authName != layer {
		if authName == "" {
			lc.VerboseLog(depth+2, "Can't parse name to authority properly. name: ", name, ", layer: ", layer)
			return RawResult{}, isCached, STATUS_AUTHFAIL, nil
		}
		lc.VerboseLog(depth+2, "Cache auth check for ", authName)
		var qAuth Question
		qAuth.Name = authName
		qAuth.Type = dns.TypeNS
		qAuth.Class = dns.ClassINET
		cachedResult, ok = lc.Cache.GetCachedResult(qAuth, true, depth+2, question.Id)
		if ok {
			isCached = true
			return cachedResult, isCached, STATUS_NOERROR, nil
		}
	}

	// Alright, we're not sure what to do, go to the wire.
	lc.VerboseLog(depth+2, "Wire lookup for name: ", question.Name, " (", question.Type, ") at nameserver: ", question.NameServer)
	result, status, err := lc.retryingLookup(question, false)

	lc.Cache.CacheUpdate(layer, result, depth+2, question.Id)
	return result, isCached, status, err
}

func (lc *RawLookupClient) tracedRetryingLookup(question Question, nameServer string, recursive bool) (RawResult, []interface{}, Status, error) {
	res, status, err := lc.retryingLookup(question, recursive)

	trace := make([]interface{}, 0)

	if lc.IsTraced {
		var t TraceStep
		t.RawResult = res
		t.DnsType = question.Type
		t.DnsClass = question.Class
		t.Name = question.Name
		t.NameServer = nameServer
		t.Layer = question.Name
		t.Depth = 1
		t.Cached = false
		trace = append(trace, t)
	}

	return res, trace, status, err
}

func (lc *RawLookupClient) iterateOnAuthorities(question Question, depth int, result RawResult,
	layer string, trace []interface{}) (RawResult, []interface{}, Status, error) {
	if len(result.Authorities) == 0 {
		return RawResult{}, trace, STATUS_NOAUTH, nil
	}
	for i, elem := range result.Authorities {
		lc.VerboseLog(depth+1, "Trying Authority: ", elem)
		ns, ns_status, layer, trace := lc.extractAuthority(elem, layer, depth, result, trace)
		lc.VerboseLog((depth + 1), "Output from extract authorities: ", ns)
		if ns_status == STATUS_ITER_TIMEOUT {
			lc.VerboseLog((depth + 2), "--> Hit iterative timeout: ")
			return RawResult{}, trace, STATUS_ITER_TIMEOUT, nil
		}
		if ns_status != STATUS_NOERROR {
			var err error
			new_status, err := handleStatus(&ns_status, err)
			// default case we continue
			if new_status == nil && err == nil {
				if i+1 == len(result.Authorities) {
					lc.VerboseLog((depth + 2), "--> Auth find Failed. Unknown error. No more authorities to try, terminating: ", ns_status)
					return RawResult{}, trace, ns_status, err
				} else {
					lc.VerboseLog((depth + 2), "--> Auth find Failed. Unknown error. Continue: ", ns_status)
					continue
				}
			} else {
				// otherwise we hit a status we know
				if i+1 == len(result.Authorities) {
					// We don't allow the continue fall through in order to report the last auth falure code, not STATUS_EROR
					lc.VerboseLog((depth + 2), "--> Final auth find non-success. Last auth. Terminating: ", ns_status)
					return RawResult{}, trace, *new_status, err
				} else {
					lc.VerboseLog((depth + 2), "--> Auth find non-success. Trying next: ", ns_status)
					continue
				}
			}
		}
		r, trace, status, err := lc.iterativeLookup(question, depth+1, layer, trace)
		if isStatusAnswer(status) {
			lc.VerboseLog((depth + 1), "--> Auth Resolution success: ", status)
			return r, trace, status, err
		} else if i+1 < len(result.Authorities) {
			lc.VerboseLog((depth + 2), "--> Auth resolution of ", ns, " Failed: ", status, ". Will try next authority")
			continue
		} else {
			// We don't allow the continue fall through in order to report the last auth falure code, not STATUS_EROR
			lc.VerboseLog((depth + 2), "--> Iterative resolution of ", question.Name, " at ", ns, " Failed. Last auth. Terminating: ", status)
			return r, trace, status, err
		}
	}
	panic("should not be able to reach here")
}

func (lc *RawLookupClient) retryingLookup(question Question, recursive bool) (RawResult, Status, error) {
	lc.VerboseLog(1, "****WIRE LOOKUP*** ", dns.TypeToString[question.Type], " ", question.Name, " ", question.NameServer)

	var origTimeout time.Duration
	if lc.Client != nil {
		origTimeout = lc.Client.Timeout
	} else {
		origTimeout = lc.TCPClient.Timeout
	}
	for i := 0; i <= lc.Retries; i++ {
		result, status, err := lc.doLookup(question, recursive)
		if (status != STATUS_TIMEOUT && status != STATUS_TEMPORARY) || i == lc.Retries {
			if lc.Client != nil {
				lc.Client.Timeout = origTimeout
			}
			if lc.TCPClient != nil {
				lc.TCPClient.Timeout = origTimeout
			}
			return result, status, err
		}
		if lc.Client != nil {
			lc.Client.Timeout = 2 * lc.Client.Timeout
		}
		if lc.TCPClient != nil {
			lc.TCPClient.Timeout = 2 * lc.TCPClient.Timeout
		}
	}
	panic("loop must return")
}

func (lc *RawLookupClient) extractAuthority(authority interface{}, layer string,
	depth int, result RawResult, trace []interface{}) (string, Status, string, []interface{}) {
	// Is it an answer
	ans, ok := authority.(Answer)
	if !ok {
		return "", STATUS_FORMERR, layer, trace
	}

	// Is the layering correct
	ok, layer = nameIsBeneath(ans.Name, layer)
	if !ok {
		return "", STATUS_AUTHFAIL, layer, trace
	}

	server := strings.TrimSuffix(ans.Answer, ".")

	// Short circuit a lookup from the glue
	// Normally this would be handled by caching, but we want to support following glue
	// that would normally be cache poison. Because it's "ok" and quite common
	res, status := checkGlue(server, depth, result)
	if status != STATUS_NOERROR {
		// Fall through to normal query
		var q Question
		q.Name = server
		q.Type = dns.TypeA
		q.Class = dns.ClassINET
		res, trace, status, _ = lc.iterativeLookup(q, depth+1, ".", trace)
	}
	if status == STATUS_ITER_TIMEOUT {
		return "", status, "", trace
	}
	if status == STATUS_NOERROR {
		// XXX we don't actually check the question here
		for _, inner_a := range res.Answers {
			inner_ans, ok := inner_a.(Answer)
			if !ok {
				continue
			}
			if inner_ans.Type == "A" {
				server := strings.TrimSuffix(inner_ans.Answer, ".") + ":53"
				return server, STATUS_NOERROR, layer, trace
			}
		}
	}
	return "", STATUS_SERVFAIL, layer, trace
}

func (lc *RawLookupClient) doLookup(question Question, recursive bool) (RawResult, Status, error) {
	return doLookupInternal(lc.Client, lc.TCPClient, lc.Conn, question, recursive, lc.ReuseSockets)
}

func doLookupInternal(udp *dns.Client, tcp *dns.Client, conn *dns.Conn, question Question,
	recursive bool, reuseSockets bool) (RawResult, Status, error) {
	res := RawResult{Answers: []interface{}{}, Authorities: []interface{}{}, Additional: []interface{}{}}
	res.Resolver = question.NameServer

	m := new(dns.Msg)
	m.SetQuestion(dotName(question.Name), question.Type)
	m.Question[0].Qclass = question.Class
	m.RecursionDesired = recursive

	var r *dns.Msg
	var err error
	if udp != nil {
		res.Protocol = "udp"

		dst, _ := net.ResolveUDPAddr("udp", question.NameServer)
		if reuseSockets {
			r, _, err = udp.ExchangeWithConnTo(m, conn, dst)
		} else {
			r, _, err = udp.ExchangeWithConn(m, conn)
		}
		// if record comes back truncated, but we have a TCP connection, try again with that
		if r != nil && (r.Truncated || r.Rcode == dns.RcodeBadTrunc) {
			if tcp != nil {
				return doLookupInternal(nil, tcp, conn, question, recursive, reuseSockets)
			} else {
				return res, STATUS_TRUNCATED, err
			}
		}
	} else {
		res.Protocol = "tcp"
		r, _, err = tcp.Exchange(m, question.NameServer)
	}
	if err != nil || r == nil {
		if nerr, ok := err.(net.Error); ok {
			if nerr.Timeout() {
				return res, STATUS_TIMEOUT, nil
			} else if nerr.Temporary() {
				return res, STATUS_TEMPORARY, err
			}
		}
		return res, STATUS_ERROR, err
	}

	if err != nil || r == nil {
		return res, STATUS_ERROR, err
	}
	if r.Rcode != dns.RcodeSuccess {
		return res, TranslateMiekgErrorCode(r.Rcode), nil
	}

	res.Flags.Response = r.Response
	res.Flags.Opcode = r.Opcode
	res.Flags.Authoritative = r.Authoritative
	res.Flags.Truncated = r.Truncated
	res.Flags.RecursionDesired = r.RecursionDesired
	res.Flags.RecursionAvailable = r.RecursionAvailable
	res.Flags.Authenticated = r.AuthenticatedData
	res.Flags.CheckingDisabled = r.CheckingDisabled
	res.Flags.ErrorCode = r.Rcode

	for _, ans := range r.Answer {
		inner := ParseAnswer(ans)
		if inner != nil {
			res.Answers = append(res.Answers, inner)
		}
	}
	for _, ans := range r.Extra {
		inner := ParseAnswer(ans)
		if inner != nil {
			res.Additional = append(res.Additional, inner)
		}
	}
	for _, ans := range r.Ns {
		inner := ParseAnswer(ans)
		if inner != nil {
			res.Authorities = append(res.Authorities, inner)
		}
	}
	return res, STATUS_NOERROR, nil
}
