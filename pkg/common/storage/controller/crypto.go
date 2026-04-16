package controller

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/tx"
	"github.com/openimsdk/tools/log"
)

type CryptoDatabase interface {
	RegisterDevice(ctx context.Context, userID, deviceID, platform, deviceModel, appVersion string) (*model.CryptoDevice, error)
	GetDevices(ctx context.Context, userID string) ([]*model.CryptoDevice, error)
	GetDevice(ctx context.Context, userID, deviceID string) (*model.CryptoDevice, error)
	RevokeDevice(ctx context.Context, userID, deviceID string) error
	TouchDevice(ctx context.Context, userID, deviceID string) error

	GetGroupKeyVersion(ctx context.Context, groupID string) (int64, error)
	BumpGroupKeyVersion(ctx context.Context, groupID, operatorUserID, eventType string) (int64, error)
	GetGroupKeyEvents(ctx context.Context, groupID string, sinceVersion int64) ([]*model.GroupKeyEvent, error)
}

type cryptoDatabase struct {
	deviceDB     database.CryptoDevice
	keyVersionDB database.GroupKeyVersion
	keyEventDB   database.GroupKeyEvent
	tx           tx.Tx
}

func NewCryptoDatabase(
	deviceDB database.CryptoDevice,
	keyVersionDB database.GroupKeyVersion,
	keyEventDB database.GroupKeyEvent,
	tx tx.Tx,
) CryptoDatabase {
	return &cryptoDatabase{
		deviceDB:     deviceDB,
		keyVersionDB: keyVersionDB,
		keyEventDB:   keyEventDB,
		tx:           tx,
	}
}

func (c *cryptoDatabase) RegisterDevice(ctx context.Context, userID, deviceID, platform, deviceModel, appVersion string) (*model.CryptoDevice, error) {
	virgilIdentity := userID + ":" + deviceID
	now := time.Now()
	device := &model.CryptoDevice{
		DeviceID:       deviceID,
		UserID:         userID,
		Platform:       platform,
		DeviceModel:    deviceModel,
		AppVersion:     appVersion,
		VirgilIdentity: virgilIdentity,
		Status:         "active",
		LastSeenAt:     now,
		CreateTime:     now,
	}
	if err := c.deviceDB.Create(ctx, device); err != nil {
		return nil, err
	}
	return device, nil
}

func (c *cryptoDatabase) GetDevices(ctx context.Context, userID string) ([]*model.CryptoDevice, error) {
	return c.deviceDB.FindByUserID(ctx, userID)
}

func (c *cryptoDatabase) GetDevice(ctx context.Context, userID, deviceID string) (*model.CryptoDevice, error) {
	return c.deviceDB.FindByUserIDAndDeviceID(ctx, userID, deviceID)
}

func (c *cryptoDatabase) RevokeDevice(ctx context.Context, userID, deviceID string) error {
	return c.deviceDB.UpdateStatus(ctx, userID, deviceID, "revoked")
}

func (c *cryptoDatabase) TouchDevice(ctx context.Context, userID, deviceID string) error {
	return c.deviceDB.UpdateLastSeen(ctx, userID, deviceID)
}

func (c *cryptoDatabase) GetGroupKeyVersion(ctx context.Context, groupID string) (int64, error) {
	v, err := c.keyVersionDB.Find(ctx, groupID)
	if err != nil {
		return 0, err
	}
	return v.GroupKeyVersion, nil
}

func (c *cryptoDatabase) BumpGroupKeyVersion(ctx context.Context, groupID, operatorUserID, eventType string) (int64, error) {
	log.ZDebug(ctx, "cryptoDatabase BumpGroupKeyVersion begin",
		"groupID", groupID,
		"operatorUserID", operatorUserID,
		"eventType", eventType,
	)
	var newVersion int64
	err := c.tx.Transaction(ctx, func(ctx context.Context) error {
		var err error
		newVersion, err = c.keyVersionDB.IncrVersion(ctx, groupID)
		if err != nil {
			return err
		}
		event := &model.GroupKeyEvent{
			EventID:         uuid.New().String(),
			GroupID:         groupID,
			GroupKeyVersion: newVersion,
			EventType:       eventType,
			OperatorUserID:  operatorUserID,
			CreateTime:      time.Now(),
		}
		return c.keyEventDB.Create(ctx, event)
	})
	if err != nil {
		log.ZError(ctx, "cryptoDatabase BumpGroupKeyVersion failed", err,
			"groupID", groupID,
			"operatorUserID", operatorUserID,
			"eventType", eventType,
		)
		return 0, err
	}
	log.ZDebug(ctx, "cryptoDatabase BumpGroupKeyVersion success",
		"groupID", groupID,
		"operatorUserID", operatorUserID,
		"eventType", eventType,
		"newGroupKeyVersion", newVersion,
	)
	return newVersion, nil
}

func (c *cryptoDatabase) GetGroupKeyEvents(ctx context.Context, groupID string, sinceVersion int64) ([]*model.GroupKeyEvent, error) {
	return c.keyEventDB.FindSinceVersion(ctx, groupID, sinceVersion)
}
