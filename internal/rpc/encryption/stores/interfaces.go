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

package stores

import (
	"context"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model/signal"
)

// IdentityStoreInterface defines the interface for identity key storage operations
type IdentityStoreInterface interface {
	Create(ctx context.Context, identityKey *signal.SignalIdentityKey) error
	Update(ctx context.Context, userID string, deviceID int32, identityKey *signal.SignalIdentityKey) error
	Get(ctx context.Context, userID string, deviceID int32) (*signal.SignalIdentityKey, error)
	Delete(ctx context.Context, userID string, deviceID int32) error
	GetByUserID(ctx context.Context, userID string) ([]*signal.SignalIdentityKey, error)
	Exists(ctx context.Context, userID string, deviceID int32) (bool, error)
}

// PreKeyStoreInterface defines the interface for one-time prekey storage operations
type PreKeyStoreInterface interface {
	Create(ctx context.Context, prekey *signal.SignalPreKey) error
	CreateBatch(ctx context.Context, prekeys []*signal.SignalPreKey) error
	GetAvailable(ctx context.Context, userID string, deviceID int32) (*signal.SignalPreKey, error)
	MarkUsed(ctx context.Context, userID string, deviceID int32, keyID uint32) error
	Delete(ctx context.Context, userID string, deviceID int32, keyID uint32) error
	DeleteAllByUserDevice(ctx context.Context, userID string, deviceID int32) error
	CountAvailable(ctx context.Context, userID string, deviceID int32) (int64, error)
	GetByKeyID(ctx context.Context, userID string, deviceID int32, keyID uint32) (*signal.SignalPreKey, error)
	CleanupUsed(ctx context.Context, olderThan time.Duration) (int64, error)
}

// SignedPreKeyStoreInterface defines the interface for signed prekey storage operations
type SignedPreKeyStoreInterface interface {
	Create(ctx context.Context, signedPrekey *signal.SignalSignedPreKey) error
	Update(ctx context.Context, userID string, deviceID int32, keyID uint32, signedPrekey *signal.SignalSignedPreKey) error
	GetActive(ctx context.Context, userID string, deviceID int32) (*signal.SignalSignedPreKey, error)
	GetByKeyID(ctx context.Context, userID string, deviceID int32, keyID uint32) (*signal.SignalSignedPreKey, error)
	SetActive(ctx context.Context, userID string, deviceID int32, keyID uint32) error
	Delete(ctx context.Context, userID string, deviceID int32, keyID uint32) error
	GetAll(ctx context.Context, userID string, deviceID int32) ([]*signal.SignalSignedPreKey, error)
	CleanupInactive(ctx context.Context, olderThan time.Duration) (int64, error)
	Exists(ctx context.Context, userID string, deviceID int32) (bool, error)
}

