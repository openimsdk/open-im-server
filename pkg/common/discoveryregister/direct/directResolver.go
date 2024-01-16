package direct

import (
	"google.golang.org/grpc/resolver"
	"strings"
)

const (
	slashSeparator = "/"
	// EndpointSepChar is the separator char in endpoints.
	EndpointSepChar = ','

	subsetSize = 32
	scheme     = "Direct"
)

type ResolverDirect struct {
}

func NewResolverDirect() *ResolverDirect {
	return &ResolverDirect{}
}

func (rd *ResolverDirect) Build(target resolver.Target, cc resolver.ClientConn, _ resolver.BuildOptions) (
	resolver.Resolver, error) {
	endpoints := strings.FieldsFunc(GetEndpoints(target), func(r rune) bool {
		return r == EndpointSepChar
	})
	endpoints = subset(endpoints, subsetSize)
	addrs := make([]resolver.Address, 0, len(endpoints))

	for _, val := range endpoints {
		addrs = append(addrs, resolver.Address{
			Addr: val,
		})
	}
	if err := cc.UpdateState(resolver.State{
		Addresses: addrs,
	}); err != nil {
		return nil, err
	}

	return &nopResolver{cc: cc}, nil
}
func init() {
	resolver.Register(&ResolverDirect{})
}
func (rd *ResolverDirect) Scheme() string {
	return scheme // return your custom scheme name
}
