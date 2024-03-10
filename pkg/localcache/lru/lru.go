// Copyright © 2024 OpenIM. All rights reserved.
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

package lru

import "github.com/hashicorp/golang-lru/v2/simplelru"

type EvictCallback[K comparable, V any] simplelru.EvictCallback[K, V]

type LRU[K comparable, V any] interface {
	Get(key K, fetch func() (V, error)) (V, error)
	Del(key K) bool
	Stop()
}

type Target interface {
	IncrGetHit()
	IncrGetSuccess()
	IncrGetFailed()

	IncrDelHit()
	IncrDelNotFound()
}
