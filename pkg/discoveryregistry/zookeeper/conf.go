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
	"github.com/go-zookeeper/zk"
)

func (s *ZkClient) RegisterConf2Registry(key string, conf []byte) error {
	exists, _, err := s.conn.Exists(s.getPath(key))
	if err != nil {
		return err
	}
	if exists {
		if err := s.conn.Delete(s.getPath(key), 0); err != nil {
			return err
		}
	}
	_, err = s.conn.Create(s.getPath(key), conf, 0, zk.WorldACL(zk.PermAll))
	if err != zk.ErrNodeExists {
		return err
	}
	return nil
}

func (s *ZkClient) GetConfFromRegistry(key string) ([]byte, error) {
	bytes, _, err := s.conn.Get(s.getPath(key))
	return bytes, err
}
