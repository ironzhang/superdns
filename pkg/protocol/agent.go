package protocol

// WatchDomainsReq is a request when SDK need to watch domains, it will send this request.
type WatchDomainsReq struct {
	Domains []string // the domain list that require to watch
}
