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

package signal

import (
	"context"
	"time"
)

const (
	SignalPreKeyCollection = "signal_prekeys"
)

// SignalPreKey represents one-time prekeys for Signal protocol
type SignalPreKey struct {
	UserID      string    `bson:"user_id" json:"userID"`
	DeviceID    int32     `bson:"device_id" json:"deviceID"`
	KeyID       uint32    `bson:"key_id" json:"keyID"`
	PublicKey   []byte    `bson:"public_key" json:"publicKey"`
	Used        bool      `bson:"used" json:"used"`
	CreatedTime time.Time `bson:"created_time" json:"createdTime"`
	UsedTime    *time.Time `bson:"used_time,omitempty" json:"usedTime,omitempty"`
}

type SignalPreKeyModelInterface interface {
	// Create creates a new prekey record
	Create(ctx context.Context, prekey *SignalPreKey) error
	
	// CreateBatch creates multiple prekey records in batch
	CreateBatch(ctx context.Context, prekeys []*SignalPreKey) error
	
	// GetAvailable retrieves an available (unused) prekey for a user/device
	GetAvailable(ctx context.Context, userID string, deviceID int32) (*SignalPreKey, error)
	
	// MarkUsed marks a prekey as used
	MarkUsed(ctx context.Context, userID string, deviceID int32, keyID uint32) error
	
	// Delete removes a prekey
	Delete(ctx context.Context, userID string, deviceID int32, keyID uint32) error
	
	// DeleteAllByUserDevice removes all prekeys for a user/device
	DeleteAllByUserDevice(ctx context.Context, userID string, deviceID int32) error
	
	// CountAvailable returns the count of available prekeys for a user/device
	CountAvailable(ctx context.Context, userID string, deviceID int32) (int64, error)
	
	// GetByKeyID retrieves a specific prekey by key ID
	GetByKeyID(ctx context.Context, userID string, deviceID int32, keyID uint32) (*SignalPreKey, error)
	
	// CleanupUsed removes used prekeys older than the specified duration
	CleanupUsed(ctx context.Context, olderThan time.Duration) (int64, error)
}

func (SignalPreKey) TableName() string {
	return SignalPreKeyCollection
}

// Indexes returns the indexes for the collection
func (SignalPreKey) Indexes() []string {
	return []string{
		"user_id",
		"user_id_device_id", // compound index for (user_id, device_id)
		"user_id_device_id_key_id", // compound unique index for (user_id, device_id, key_id)
		"user_id_device_id_used", // compound index for (user_id, device_id, used)
		"used_time",
		"created_time",
	}
}