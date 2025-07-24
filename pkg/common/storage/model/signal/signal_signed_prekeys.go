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
	SignalSignedPreKeyCollection = "signal_signed_prekeys"
)

// SignalSignedPreKey represents signed prekeys for Signal protocol
type SignalSignedPreKey struct {
	UserID      string    `bson:"user_id" json:"userID"`
	DeviceID    int32     `bson:"device_id" json:"deviceID"`
	KeyID       uint32    `bson:"key_id" json:"keyID"`
	PublicKey   []byte    `bson:"public_key" json:"publicKey"`
	Signature   []byte    `bson:"signature" json:"signature"`
	CreatedTime time.Time `bson:"created_time" json:"createdTime"`
	Active      bool      `bson:"active" json:"active"` // Whether this is the active signed prekey
}

type SignalSignedPreKeyModelInterface interface {
	// Create creates a new signed prekey record
	Create(ctx context.Context, signedPrekey *SignalSignedPreKey) error
	
	// Update updates an existing signed prekey (for rotation)
	Update(ctx context.Context, userID string, deviceID int32, keyID uint32, signedPrekey *SignalSignedPreKey) error
	
	// GetActive retrieves the active signed prekey for a user/device
	GetActive(ctx context.Context, userID string, deviceID int32) (*SignalSignedPreKey, error)
	
	// GetByKeyID retrieves a specific signed prekey by key ID
	GetByKeyID(ctx context.Context, userID string, deviceID int32, keyID uint32) (*SignalSignedPreKey, error)
	
	// SetActive marks a signed prekey as active and deactivates others
	SetActive(ctx context.Context, userID string, deviceID int32, keyID uint32) error
	
	// Delete removes a signed prekey
	Delete(ctx context.Context, userID string, deviceID int32, keyID uint32) error
	
	// GetAll retrieves all signed prekeys for a user/device
	GetAll(ctx context.Context, userID string, deviceID int32) ([]*SignalSignedPreKey, error)
	
	// CleanupInactive removes inactive signed prekeys older than the specified duration
	CleanupInactive(ctx context.Context, olderThan time.Duration) (int64, error)
	
	// Exists checks if a signed prekey exists
	Exists(ctx context.Context, userID string, deviceID int32) (bool, error)
}

func (SignalSignedPreKey) TableName() string {
	return SignalSignedPreKeyCollection
}

// Indexes returns the indexes for the collection
func (SignalSignedPreKey) Indexes() []string {
	return []string{
		"user_id",
		"user_id_device_id", // compound index for (user_id, device_id)
		"user_id_device_id_key_id", // compound unique index for (user_id, device_id, key_id)
		"user_id_device_id_active", // compound index for (user_id, device_id, active)
		"active_created_time", // compound index for (active, created_time)
		"created_time",
	}
}