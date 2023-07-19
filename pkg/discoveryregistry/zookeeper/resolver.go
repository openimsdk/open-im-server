// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package zookeeper

import (
	"context"
	"strings"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"

	"google.golang.org/grpc/resolver"
)

type Resolver struct {
	target resolver.Target
	cc     resolver.ClientConn
	addrs  []resolver.Address

	getConnsRemote func(serviceName string) (conns []resolver.Address, err error)
}

func (r *Resolver) ResolveNowZK(o resolver.ResolveNowOptions) {
	log.ZDebug(
		context.Background(),
		"start resolve now",
		"target",
		r.target,
		"cc",
		r.cc.UpdateState,
		"serviceName",
		strings.TrimLeft(r.target.URL.Path, "/"),
	)
	newConns, err := r.getConnsRemote(strings.TrimLeft(r.target.URL.Path, "/"))
	if err != nil {
		log.ZError(context.Background(), "resolve now error", err, "target", r.target)
		return
	}
	r.addrs = newConns
	if err := r.cc.UpdateState(resolver.State{Addresses: newConns}); err != nil {
		log.ZError(
			context.Background(),
			"UpdateState error, conns is nil from svr",
			err,
			"conns",
			newConns,
			"zk path",
			r.target.URL.Path,
		)
		return
	}
	log.ZDebug(context.Background(), "resolve now finished", "target", r.target, "conns", r.addrs)
}

func (r *Resolver) ResolveNow(o resolver.ResolveNowOptions) {}

func (s *Resolver) Close() {}

func (s *ZkClient) Build(
	target resolver.Target,
	cc resolver.ClientConn,
	opts resolver.BuildOptions,
) (resolver.Resolver, error) {
	s.logger.Printf("build resolver: %+v, cc: %+v", target, cc.UpdateState)
	serviceName := strings.TrimLeft(target.URL.Path, "/")
	if oldResolver, ok := s.resolvers[serviceName]; ok {
		s.logger.Printf("rpc resolver exist: %+v, cc: %+v, key: %s", target, cc.UpdateState, serviceName)
		return oldResolver, nil
	}
	r := &Resolver{}
	r.target = target
	r.cc = cc
	r.getConnsRemote = s.GetConnsRemote
	r.ResolveNowZK(resolver.ResolveNowOptions{})
	s.lock.Lock()
	defer s.lock.Unlock()
	s.resolvers[serviceName] = r
	s.logger.Printf("build resolver finished: %+v, cc: %+v, key: %s", target, cc.UpdateState, serviceName)
	return r, nil
}

func (s *ZkClient) Scheme() string { return s.scheme }
