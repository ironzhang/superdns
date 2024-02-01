package protocol

import "github.com/ironzhang/superlib/timeutil"

// SubscribeDomainsReq is a request to subscribe domains.
type SubscribeDomainsReq struct {
	Domains      []string          // the domain list that require to subscribe
	TTL          timeutil.Duration // time to live, <= 0 means forever
	Asynchronous bool              // asynchronous call
}

// ListSubscribeDomainsResp is a response.
type ListSubscribeDomainsResp struct {
	Domains []string // the domain list that the agent is subscribing
}
