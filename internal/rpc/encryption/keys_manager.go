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

package encryption

import (
	"context"
	"fmt"
	"time"

	"github.com/openimsdk/open-im-server/v3/internal/rpc/encryption/stores"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model/signal"
	"github.com/openimsdk/tools/log"
)

const (
	// Maximum number of one-time prekeys that can be uploaded at once
	MaxOneTimePreKeys = 100
)

type KeysManager struct {
	identityStore     stores.IdentityStoreInterface
	preKeyStore       stores.PreKeyStoreInterface
	signedPreKeyStore stores.SignedPreKeyStoreInterface
}

func NewKeysManager(
	identityStore stores.IdentityStoreInterface,
	preKeyStore stores.PreKeyStoreInterface,
	signedPreKeyStore stores.SignedPreKeyStoreInterface,
) *KeysManager {
	return &KeysManager{
		identityStore:     identityStore,
		preKeyStore:       preKeyStore,
		signedPreKeyStore: signedPreKeyStore,
	}
}

// GetIdentityKey retrieves the identity key for a user/device
func (km *KeysManager) GetIdentityKey(ctx context.Context, userID string, deviceID int32) (*signal.SignalIdentityKey, error) {
	return km.identityStore.Get(ctx, userID, deviceID)
}

// SetIdentityKey sets the identity key for a user/device
func (km *KeysManager) SetIdentityKey(ctx context.Context, userID string, deviceID int32, identityKey []byte, registrationID int32) error {
	now := time.Now()

	// Check if identity key already exists
	existing, err := km.identityStore.Get(ctx, userID, deviceID)
	if err == nil && existing != nil {
		// Update existing identity key
		existing.IdentityKey = identityKey
		existing.RegistrationID = registrationID
		existing.UpdatedTime = now
		return km.identityStore.Update(ctx, userID, deviceID, existing)
	}

	// Create new identity key
	identityKeyRecord := &signal.SignalIdentityKey{
		UserID:         userID,
		DeviceID:       deviceID,
		IdentityKey:    identityKey,
		RegistrationID: registrationID,
		CreatedTime:    now,
		UpdatedTime:    now,
	}

	return km.identityStore.Create(ctx, identityKeyRecord)
}

// GetActiveSignedPreKey retrieves the active signed prekey for a user/device
func (km *KeysManager) GetActiveSignedPreKey(ctx context.Context, userID string, deviceID int32) (*signal.SignalSignedPreKey, error) {
	return km.signedPreKeyStore.GetActive(ctx, userID, deviceID)
}

// SetSignedPreKey sets a signed prekey for a user/device
func (km *KeysManager) SetSignedPreKey(ctx context.Context, userID string, deviceID int32, signedPreKey *SignedPreKeyResponse) error {
	now := time.Now()

	signedPreKeyRecord := &signal.SignalSignedPreKey{
		UserID:      userID,
		DeviceID:    deviceID,
		KeyID:       signedPreKey.KeyId,
		PublicKey:   signedPreKey.PublicKey,
		Signature:   signedPreKey.Signature,
		CreatedTime: now,
		Active:      true,
	}

	// Deactivate existing signed prekeys and set this one as active
	err := km.signedPreKeyStore.SetActive(ctx, userID, deviceID, signedPreKey.KeyId)
	if err != nil {
		log.ZWarn(ctx, "failed to deactivate existing signed prekeys", err)
	}

	// Create or update the signed prekey
	existing, err := km.signedPreKeyStore.GetByKeyID(ctx, userID, deviceID, signedPreKey.KeyId)
	if err == nil && existing != nil {
		return km.signedPreKeyStore.Update(ctx, userID, deviceID, signedPreKey.KeyId, signedPreKeyRecord)
	}

	return km.signedPreKeyStore.Create(ctx, signedPreKeyRecord)
}

// GetOneTimePreKey retrieves an available one-time prekey and marks it as used
func (km *KeysManager) GetOneTimePreKey(ctx context.Context, userID string, deviceID int32) (*signal.SignalPreKey, error) {
	preKey, err := km.preKeyStore.GetAvailable(ctx, userID, deviceID)
	if err != nil {
		return nil, fmt.Errorf("no available one-time prekey: %w", err)
	}

	// Mark the prekey as used
	err = km.preKeyStore.MarkUsed(ctx, userID, deviceID, preKey.KeyID)
	if err != nil {
		log.ZError(ctx, "failed to mark prekey as used", err, "userID", userID, "deviceID", deviceID, "keyID", preKey.KeyID)
		// Don't fail the request, but log the error
	}

	return preKey, nil
}

// SetOneTimePreKeys sets multiple one-time prekeys for a user/device
func (km *KeysManager) SetOneTimePreKeys(ctx context.Context, userID string, deviceID int32, preKeys []*PreKeyResponse) (int, error) {
	if len(preKeys) > MaxOneTimePreKeys {
		return 0, fmt.Errorf("too many one-time prekeys: %d (max: %d)", len(preKeys), MaxOneTimePreKeys)
	}

	now := time.Now()
	var preKeyRecords []*signal.SignalPreKey

	for _, preKey := range preKeys {
		preKeyRecord := &signal.SignalPreKey{
			UserID:      userID,
			DeviceID:    deviceID,
			KeyID:       preKey.KeyId,
			PublicKey:   preKey.PublicKey,
			Used:        false,
			CreatedTime: now,
		}
		preKeyRecords = append(preKeyRecords, preKeyRecord)
	}

	err := km.preKeyStore.CreateBatch(ctx, preKeyRecords)
	if err != nil {
		return 0, fmt.Errorf("failed to create one-time prekeys: %w", err)
	}

	return len(preKeyRecords), nil
}

// GetPreKeyCount returns the count of available one-time prekeys for a user/device
func (km *KeysManager) GetPreKeyCount(ctx context.Context, userID string, deviceID int32) (int64, error) {
	return km.preKeyStore.CountAvailable(ctx, userID, deviceID)
}

// GetSignedPreKeyInfo returns information about signed prekey existence and last rotation
func (km *KeysManager) GetSignedPreKeyInfo(ctx context.Context, userID string, deviceID int32) (exists bool, lastRotation time.Time, err error) {
	signedPreKey, err := km.signedPreKeyStore.GetActive(ctx, userID, deviceID)
	if err != nil {
		return false, time.Time{}, err
	}

	if signedPreKey != nil {
		return true, signedPreKey.CreatedTime, nil
	}

	return false, time.Time{}, nil
}

// CleanupExpiredKeys removes expired keys from storage
func (km *KeysManager) CleanupExpiredKeys(ctx context.Context) error {
	// Cleanup used one-time prekeys older than 7 days
	usedPreKeysCleanedCount, err := km.preKeyStore.CleanupUsed(ctx, 7*24*time.Hour)
	if err != nil {
		log.ZError(ctx, "failed to cleanup used prekeys", err)
	} else {
		log.ZInfo(ctx, "cleaned up used prekeys", "count", usedPreKeysCleanedCount)
	}

	// Cleanup inactive signed prekeys older than 30 days
	inactiveSignedPreKeysCount, err := km.signedPreKeyStore.CleanupInactive(ctx, 30*24*time.Hour)
	if err != nil {
		log.ZError(ctx, "failed to cleanup inactive signed prekeys", err)
	} else {
		log.ZInfo(ctx, "cleaned up inactive signed prekeys", "count", inactiveSignedPreKeysCount)
	}

	return nil
}

// ValidateSignedPreKey validates the signature of a signed prekey
func (km *KeysManager) ValidateSignedPreKey(ctx context.Context, userID string, deviceID int32, signedPreKey *SignedPreKeyResponse) error {
	// Get the identity key to verify the signature
	_, err := km.GetIdentityKey(ctx, userID, deviceID)
	if err != nil {
		return fmt.Errorf("failed to get identity key for signature validation: %w", err)
	}

	// TODO: Implement actual signature validation using Signal Protocol
	// This would require integrating with the Signal Protocol library
	// For now, we'll skip the validation
	log.ZInfo(ctx, "signed prekey signature validation skipped (not implemented)",
		"userID", userID, "deviceID", deviceID, "keyID", signedPreKey.KeyId)

	return nil
}

// RotateSignedPreKey creates a new signed prekey and deactivates the old one
func (km *KeysManager) RotateSignedPreKey(ctx context.Context, userID string, deviceID int32, newSignedPreKey *SignedPreKeyResponse) error {
	// Validate the new signed prekey
	err := km.ValidateSignedPreKey(ctx, userID, deviceID, newSignedPreKey)
	if err != nil {
		return fmt.Errorf("signed prekey validation failed: %w", err)
	}

	// Set the new signed prekey (this will automatically deactivate the old one)
	return km.SetSignedPreKey(ctx, userID, deviceID, newSignedPreKey)
}

// GetPreKeyBundleForUser retrieves all necessary keys for X3DH key agreement
func (km *KeysManager) GetPreKeyBundleForUser(ctx context.Context, userID string, deviceID int32) (map[string]interface{}, error) {
	// Get identity key
	identityKey, err := km.GetIdentityKey(ctx, userID, deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get identity key: %w", err)
	}

	// Get signed prekey
	signedPreKey, err := km.GetActiveSignedPreKey(ctx, userID, deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get signed prekey: %w", err)
	}

	// Get one-time prekey (optional)
	oneTimePreKey, err := km.GetOneTimePreKey(ctx, userID, deviceID)
	if err != nil {
		log.ZWarn(ctx, "no one-time prekey available", err)
		oneTimePreKey = nil
	}

	bundle := map[string]interface{}{
		"identityKey": map[string]interface{}{
			"identityKey":    identityKey.IdentityKey,
			"registrationId": identityKey.RegistrationID,
			"createdTime":    identityKey.CreatedTime.Unix(),
		},
		"signedPreKey": map[string]interface{}{
			"keyId":       signedPreKey.KeyID,
			"publicKey":   signedPreKey.PublicKey,
			"signature":   signedPreKey.Signature,
			"createdTime": signedPreKey.CreatedTime.Unix(),
		},
		"registrationId": identityKey.RegistrationID,
	}

	if oneTimePreKey != nil {
		bundle["oneTimePreKey"] = map[string]interface{}{
			"keyId":     oneTimePreKey.KeyID,
			"publicKey": oneTimePreKey.PublicKey,
		}
	}

	return bundle, nil
}

// Response structures used in key manager
type SignedPreKeyResponse struct {
	KeyId     uint32 `json:"keyId"`
	PublicKey []byte `json:"publicKey"`
	Signature []byte `json:"signature"`
}

type PreKeyResponse struct {
	KeyId     uint32 `json:"keyId"`
	PublicKey []byte `json:"publicKey"`
}
