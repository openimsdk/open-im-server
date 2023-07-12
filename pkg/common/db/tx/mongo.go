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

package tx

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
)

func NewMongo(client *mongo.Client) CtxTx {
	return &_Mongo{
		client: client,
	}
}

type _Mongo struct {
	client *mongo.Client
}

func (m *_Mongo) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	sess, err := m.client.StartSession()
	if err != nil {
		return err
	}
	sCtx := mongo.NewSessionContext(ctx, sess)
	defer sess.EndSession(sCtx)
	if err := fn(sCtx); err != nil {
		_ = sess.AbortTransaction(sCtx)
		return err
	}
	return utils.Wrap(sess.CommitTransaction(sCtx), "")
}
