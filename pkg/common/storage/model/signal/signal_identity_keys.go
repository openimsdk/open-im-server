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
	SignalIdentityKeyCollection = "signal_identity_keys"
)

// SignalIdentityKey represents the identity key for Signal protocol
type SignalIdentityKey struct {
	UserID         string    `bson:"user_id" json:"userID"`
	DeviceID       int32     `bson:"device_id" json:"deviceID"`
	IdentityKey    []byte    `bson:"identity_key" json:"identityKey"`
	RegistrationID int32     `bson:"registration_id" json:"registrationID"`
	CreatedTime    time.Time `bson:"created_time" json:"createdTime"`
	UpdatedTime    time.Time `bson:"updated_time" json:"updatedTime"`
}

type SignalIdentityKeyModelInterface interface {
	// Create creates a new identity key record
	Create(ctx context.Context, identityKey *SignalIdentityKey) error
	
	// Update updates an existing identity key
	Update(ctx context.Context, userID string, deviceID int32, identityKey *SignalIdentityKey) error
	
	// Get retrieves an identity key by user ID and device ID
	Get(ctx context.Context, userID string, deviceID int32) (*SignalIdentityKey, error)
	
	// Delete removes an identity key
	Delete(ctx context.Context, userID string, deviceID int32) error
	
	// GetByUserID retrieves all identity keys for a user
	GetByUserID(ctx context.Context, userID string) ([]*SignalIdentityKey, error)
	
	// Exists checks if an identity key exists
	Exists(ctx context.Context, userID string, deviceID int32) (bool, error)
}

func (SignalIdentityKey) TableName() string {
	return SignalIdentityKeyCollection
}

// Indexes returns the indexes for the collection
func (SignalIdentityKey) Indexes() []string {
	return []string{
		"user_id",
		"user_id_device_id", // compound index for (user_id, device_id)
		"created_time",
	}
}