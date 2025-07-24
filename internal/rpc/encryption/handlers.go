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
	"encoding/base64"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/openimsdk/tools/log"
)

// Request/Response structures for HTTP API
type GetPreKeysResponse struct {
	IdentityKey   *IdentityKeyInfo    `json:"identityKey"`
	SignedPreKey  *SignedPreKeyInfo   `json:"signedPreKey"`
	OneTimePreKey *PreKeyInfo         `json:"oneTimePreKey,omitempty"`
	RegistrationID int32              `json:"registrationId"`
}

type IdentityKeyInfo struct {
	IdentityKey    string `json:"identityKey"`
	RegistrationID int32  `json:"registrationId"`
	CreatedTime    int64  `json:"createdTime"`
}

type PreKeyInfo struct {
	KeyID     uint32 `json:"keyId"`
	PublicKey string `json:"publicKey"`
}

type SignedPreKeyInfo struct {
	KeyID       uint32 `json:"keyId"`
	PublicKey   string `json:"publicKey"`
	Signature   string `json:"signature"`
	CreatedTime int64  `json:"createdTime"`
}

type SetPreKeysRequest struct {
	IdentityKey     string             `json:"identityKey,omitempty"`
	SignedPreKey    *SignedPreKeyInfo  `json:"signedPreKey,omitempty"`
	OneTimePreKeys  []PreKeyInfo       `json:"oneTimePreKeys,omitempty"`
	RegistrationID  int32              `json:"registrationId,omitempty"`
}

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// GetPreKeys handles GET /api/v1/encryption/prekeys/:user_id/:device_id
func (s *Server) GetPreKeys(c *gin.Context) {
	userID := c.Param("user_id")
	deviceIDStr := c.Param("device_id")
	
	deviceID, err := strconv.ParseInt(deviceIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "Invalid device_id",
		})
		return
	}
	
	log.ZInfo(c.Request.Context(), "GetPreKeys", "userID", userID, "deviceID", deviceID)
	
	// Get identity key
	identityKey, err := s.keysManager.GetIdentityKey(c.Request.Context(), userID, int32(deviceID))
	if err != nil {
		log.ZError(c.Request.Context(), "failed to get identity key", err)
		c.JSON(http.StatusNotFound, APIResponse{
			Code:    404,
			Message: "Identity key not found",
		})
		return
	}
	
	// Get signed prekey
	signedPreKey, err := s.keysManager.GetActiveSignedPreKey(c.Request.Context(), userID, int32(deviceID))
	if err != nil {
		log.ZError(c.Request.Context(), "failed to get signed prekey", err)
		c.JSON(http.StatusNotFound, APIResponse{
			Code:    404,
			Message: "Signed prekey not found",
		})
		return
	}
	
	// Get one-time prekey (optional)
	oneTimePreKey, err := s.keysManager.GetOneTimePreKey(c.Request.Context(), userID, int32(deviceID))
	if err != nil {
		log.ZWarn(c.Request.Context(), "no one-time prekey available", err)
		oneTimePreKey = nil
	}
	
	response := &GetPreKeysResponse{
		IdentityKey: &IdentityKeyInfo{
			IdentityKey:    base64.StdEncoding.EncodeToString(identityKey.IdentityKey),
			RegistrationID: identityKey.RegistrationID,
			CreatedTime:    identityKey.CreatedTime.Unix(),
		},
		SignedPreKey: &SignedPreKeyInfo{
			KeyID:       signedPreKey.KeyID,
			PublicKey:   base64.StdEncoding.EncodeToString(signedPreKey.PublicKey),
			Signature:   base64.StdEncoding.EncodeToString(signedPreKey.Signature),
			CreatedTime: signedPreKey.CreatedTime.Unix(),
		},
		RegistrationID: identityKey.RegistrationID,
	}
	
	if oneTimePreKey != nil {
		response.OneTimePreKey = &PreKeyInfo{
			KeyID:     oneTimePreKey.KeyID,
			PublicKey: base64.StdEncoding.EncodeToString(oneTimePreKey.PublicKey),
		}
	}
	
	c.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "success",
		Data:    response,
	})
}

// SetPreKeys handles POST /api/v1/encryption/prekeys/:user_id/:device_id
func (s *Server) SetPreKeys(c *gin.Context) {
	userID := c.Param("user_id")
	deviceIDStr := c.Param("device_id")
	
	deviceID, err := strconv.ParseInt(deviceIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "Invalid device_id",
		})
		return
	}
	
	var req SetPreKeysRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "Invalid request body",
		})
		return
	}
	
	log.ZInfo(c.Request.Context(), "SetPreKeys", "userID", userID, "deviceID", deviceID)
	
	// Set identity key if provided
	if req.IdentityKey != "" {
		identityKeyBytes, err := base64.StdEncoding.DecodeString(req.IdentityKey)
		if err != nil {
			c.JSON(http.StatusBadRequest, APIResponse{
				Code:    400,
				Message: "Invalid identity key encoding",
			})
			return
		}
		
		err = s.keysManager.SetIdentityKey(c.Request.Context(), userID, int32(deviceID), identityKeyBytes, req.RegistrationID)
		if err != nil {
			log.ZError(c.Request.Context(), "failed to set identity key", err)
			c.JSON(http.StatusInternalServerError, APIResponse{
				Code:    500,
				Message: "Failed to set identity key",
			})
			return
		}
	}
	
	// Set signed prekey if provided
	if req.SignedPreKey != nil {
		publicKeyBytes, err := base64.StdEncoding.DecodeString(req.SignedPreKey.PublicKey)
		if err != nil {
			c.JSON(http.StatusBadRequest, APIResponse{
				Code:    400,
				Message: "Invalid signed prekey public key encoding",
			})
			return
		}
		
		signatureBytes, err := base64.StdEncoding.DecodeString(req.SignedPreKey.Signature)
		if err != nil {
			c.JSON(http.StatusBadRequest, APIResponse{
				Code:    400,
				Message: "Invalid signed prekey signature encoding",
			})
			return
		}
		
		signedPreKeyData := &SignedPreKeyResponse{
			KeyId:     req.SignedPreKey.KeyID,
			PublicKey: publicKeyBytes,
			Signature: signatureBytes,
		}
		
		err = s.keysManager.SetSignedPreKey(c.Request.Context(), userID, int32(deviceID), signedPreKeyData)
		if err != nil {
			log.ZError(c.Request.Context(), "failed to set signed prekey", err)
			c.JSON(http.StatusInternalServerError, APIResponse{
				Code:    500,
				Message: "Failed to set signed prekey",
			})
			return
		}
	}
	
	// Set one-time prekeys
	if len(req.OneTimePreKeys) > 0 {
		var preKeyData []*PreKeyResponse
		for _, pk := range req.OneTimePreKeys {
			publicKeyBytes, err := base64.StdEncoding.DecodeString(pk.PublicKey)
			if err != nil {
				c.JSON(http.StatusBadRequest, APIResponse{
					Code:    400,
					Message: "Invalid one-time prekey public key encoding",
				})
				return
			}
			
			preKeyData = append(preKeyData, &PreKeyResponse{
				KeyId:     pk.KeyID,
				PublicKey: publicKeyBytes,
			})
		}
		
		acceptedCount, err := s.keysManager.SetOneTimePreKeys(c.Request.Context(), userID, int32(deviceID), preKeyData)
		if err != nil {
			log.ZError(c.Request.Context(), "failed to set one-time prekeys", err)
			c.JSON(http.StatusInternalServerError, APIResponse{
				Code:    500,
				Message: "Failed to set one-time prekeys",
			})
			return
		}
		
		c.JSON(http.StatusOK, APIResponse{
			Code:    0,
			Message: "success",
			Data: map[string]interface{}{
				"preKeysAccepted": acceptedCount,
			},
		})
		return
	}
	
	c.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "success",
	})
}

// GetPreKeyCount handles GET /api/v1/encryption/prekeys/:user_id/:device_id/count
func (s *Server) GetPreKeyCount(c *gin.Context) {
	userID := c.Param("user_id")
	deviceIDStr := c.Param("device_id")
	
	deviceID, err := strconv.ParseInt(deviceIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "Invalid device_id",
		})
		return
	}
	
	count, err := s.keysManager.GetPreKeyCount(c.Request.Context(), userID, int32(deviceID))
	if err != nil {
		log.ZError(c.Request.Context(), "failed to get prekey count", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "Failed to get prekey count",
		})
		return
	}
	
	signedPreKeyExists, lastRotation, err := s.keysManager.GetSignedPreKeyInfo(c.Request.Context(), userID, int32(deviceID))
	if err != nil {
		signedPreKeyExists = false
	}
	
	data := map[string]interface{}{
		"oneTimePreKeyCount": count,
		"signedPreKeyExists": signedPreKeyExists,
	}
	
	if !lastRotation.IsZero() {
		data["lastSignedPreKeyRotation"] = lastRotation.Unix()
	}
	
	c.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// GetIdentityKey handles GET /api/v1/encryption/identity/:user_id/:device_id
func (s *Server) GetIdentityKey(c *gin.Context) {
	userID := c.Param("user_id")
	deviceIDStr := c.Param("device_id")
	
	deviceID, err := strconv.ParseInt(deviceIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "Invalid device_id",
		})
		return
	}
	
	identityKey, err := s.keysManager.GetIdentityKey(c.Request.Context(), userID, int32(deviceID))
	if err != nil {
		log.ZError(c.Request.Context(), "failed to get identity key", err)
		c.JSON(http.StatusNotFound, APIResponse{
			Code:    404,
			Message: "Identity key not found",
		})
		return
	}
	
	response := &IdentityKeyInfo{
		IdentityKey:    base64.StdEncoding.EncodeToString(identityKey.IdentityKey),
		RegistrationID: identityKey.RegistrationID,
		CreatedTime:    identityKey.CreatedTime.Unix(),
	}
	
	c.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "success",
		Data:    response,
	})
}