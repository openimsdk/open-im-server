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

package prommetrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net"
)

var (
	MsgInsertRedisSuccessCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "msg_insert_redis_success_total",
		Help: "The number of successful insert msg to redis",
	})
	MsgInsertRedisFailedCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "msg_insert_redis_failed_total",
		Help: "The number of failed insert msg to redis",
	})
	MsgInsertMongoSuccessCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "msg_insert_mongo_success_total",
		Help: "The number of successful insert msg to mongo",
	})
	MsgInsertMongoFailedCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "msg_insert_mongo_failed_total",
		Help: "The number of failed insert msg to mongo",
	})
	SeqSetFailedCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "seq_set_failed_total",
		Help: "The number of failed set seq",
	})
)

func TransferInit(listener net.Listener) error {
	reg := prometheus.NewRegistry()
	cs := append(
		baseCollector,
		MsgInsertRedisSuccessCounter,
		MsgInsertRedisFailedCounter,
		MsgInsertMongoSuccessCounter,
		MsgInsertMongoFailedCounter,
		SeqSetFailedCounter,
	)
	return Init(reg, listener, commonPath, promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}), cs...)
}
