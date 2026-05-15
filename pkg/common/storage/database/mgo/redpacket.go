package mgo

import (
	"context"
	"math/big"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ---- RedPacket ----

type RedPacketMgo struct {
	coll *mongo.Collection
}

func NewRedPacketMongo(db *mongo.Database) (database.RedPacket, error) {
	coll := db.Collection("red_packet")
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "biz_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "packet_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "group_id", Value: 1}},
		},
	})
	if err != nil {
		return nil, err
	}
	return &RedPacketMgo{coll: coll}, nil
}

func (m *RedPacketMgo) Create(ctx context.Context, rp *model.RedPacket) error {
	_, err := m.coll.InsertOne(ctx, rp)
	return err
}

func (m *RedPacketMgo) GetByBizID(ctx context.Context, bizID string) (*model.RedPacket, error) {
	var rp model.RedPacket
	err := m.coll.FindOne(ctx, bson.M{"biz_id": bizID}).Decode(&rp)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errs.ErrRecordNotFound.WrapMsg("red packet not found", "bizID", bizID)
		}
		return nil, err
	}
	return &rp, nil
}

func (m *RedPacketMgo) GetByPacketID(ctx context.Context, packetID string) (*model.RedPacket, error) {
	var rp model.RedPacket
	err := m.coll.FindOne(ctx, bson.M{"packet_id": packetID}).Decode(&rp)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errs.ErrRecordNotFound.WrapMsg("red packet not found", "packetID", packetID)
		}
		return nil, err
	}
	return &rp, nil
}

func (m *RedPacketMgo) UpdateCreated(ctx context.Context, rp *model.RedPacket) error {
	updates := bson.M{
		"chain_type":        rp.ChainType,
		"packet_id":         rp.PacketID,
		"tx_hash":           rp.TxHash,
		"chain_id":          rp.ChainID,
		"contract_address":  rp.ContractAddress,
		"creator_wallet":    rp.CreatorWallet,
		"packet_type":       rp.PacketType,
		"token":             rp.Token,
		"total_amount":      rp.TotalAmount,
		"total_shares":      rp.TotalShares,
		"expiry_at":         rp.ExpiryAt,
		"group_id":          rp.GroupID,
		"scope_type":        rp.ScopeType,
		"receiver_user_id":  rp.ReceiverUserID,
		"receiver_user_ids": rp.ReceiverUserIDs,
		"status":            rp.Status,
		"updated_at":        time.Now(),
	}
	res, err := m.coll.UpdateOne(ctx, bson.M{"biz_id": rp.BizID}, bson.M{"$set": updates})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errs.ErrRecordNotFound.WrapMsg("red packet not found", "bizID", rp.BizID)
	}
	return nil
}

func (m *RedPacketMgo) UpdateStatus(ctx context.Context, packetID, status string) error {
	res, err := m.coll.UpdateOne(ctx, bson.M{"packet_id": packetID},
		bson.M{"$set": bson.M{"status": status, "updated_at": time.Now()}})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errs.ErrRecordNotFound.WrapMsg("red packet not found", "packetID", packetID)
	}
	return nil
}

func (m *RedPacketMgo) UpdateClaimProgress(ctx context.Context, packetID, claimedAmount, status, claimTxHash string) error {
	var rp model.RedPacket
	err := m.coll.FindOne(ctx, bson.M{"packet_id": packetID}).Decode(&rp)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errs.ErrRecordNotFound.WrapMsg("red packet not found", "packetID", packetID)
		}
		return err
	}

	totalClaimed := addNumericStrings(rp.ClaimedAmount, claimedAmount)
	nextShares := rp.ClaimedShares + 1

	// Auto-derive status when the caller does not force one.
	nextStatus := status
	if nextStatus == "" {
		if rp.PacketType == 2 {
			nextStatus = "COMPLETED"
		} else if rp.TotalShares > 0 && nextShares >= rp.TotalShares {
			nextStatus = "COMPLETED"
		} else {
			tcBig, tok := new(big.Int).SetString(totalClaimed, 10)
			taBig, taok := new(big.Int).SetString(rp.TotalAmount, 10)
			if tok && taok && tcBig.Cmp(taBig) >= 0 {
				nextStatus = "COMPLETED"
			}
		}
	}

	setFields := bson.M{
		"claimed_amount": totalClaimed,
		"claimed_shares": nextShares,
		"updated_at":     time.Now(),
	}
	if nextStatus != "" {
		setFields["status"] = nextStatus
	}

	// The $addToSet + $ne filter makes the whole update idempotent per claimTxHash:
	// if two code paths (RPC handler and indexer) both attempt to process the same
	// transaction, only the first UpdateOne will match and the second is a no-op.
	filter := bson.M{"packet_id": packetID}
	if claimTxHash != "" {
		filter["processed_claim_hashes"] = bson.M{"$ne": claimTxHash}
	}
	update := bson.M{"$set": setFields}
	if claimTxHash != "" {
		update["$addToSet"] = bson.M{"processed_claim_hashes": claimTxHash}
	}

	_, err = m.coll.UpdateOne(ctx, filter, update)
	return err
}

func addNumericStrings(current, delta string) string {
	left := new(big.Int)
	if current != "" {
		left.SetString(current, 10)
	}
	right := new(big.Int)
	if delta != "" {
		right.SetString(delta, 10)
	}
	return new(big.Int).Add(left, right).String()
}

// ---- RedPacketClaim ----

type RedPacketClaimMgo struct {
	coll *mongo.Collection
}

func NewRedPacketClaimMongo(db *mongo.Database) (database.RedPacketClaim, error) {
	coll := db.Collection("red_packet_claim")
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "claim_tx_hash", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "packet_id", Value: 1}, {Key: "user_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "packet_id", Value: 1}, {Key: "claimer_wallet", Value: 1}},
		},
	})
	if err != nil {
		return nil, err
	}
	return &RedPacketClaimMgo{coll: coll}, nil
}

func (m *RedPacketClaimMgo) Save(ctx context.Context, claim *model.RedPacketClaim) error {
	if claim.UserID != "" {
		var existing model.RedPacketClaim
		err := m.coll.FindOne(ctx, bson.M{
			"packet_id": claim.PacketID,
			"user_id":   claim.UserID,
		}).Decode(&existing)
		if err == nil {
			updates := bson.M{
				"claimer_wallet": claim.ClaimerWallet,
				"auth_nonce":     claim.AuthNonce,
				"claim_tx_hash":  claim.ClaimTxHash,
				"claimed_amount": claim.ClaimedAmount,
				"block_number":   claim.BlockNumber,
				"status":         claim.Status,
				"updated_at":     claim.UpdatedAt,
			}
			_, err := m.coll.UpdateOne(ctx,
				bson.M{"packet_id": claim.PacketID, "user_id": claim.UserID},
				bson.M{"$set": updates})
			return err
		}
		if err != mongo.ErrNoDocuments {
			return err
		}
	}

	_, err := m.coll.UpdateOne(ctx,
		bson.M{"claim_tx_hash": claim.ClaimTxHash},
		bson.M{"$set": claim},
		options.Update().SetUpsert(true),
	)
	return err
}

func (m *RedPacketClaimMgo) GetByPacketIDAndClaimer(ctx context.Context, packetID, claimer string) (*model.RedPacketClaim, error) {
	var claim model.RedPacketClaim
	err := m.coll.FindOne(ctx,
		bson.M{"packet_id": packetID, "claimer_wallet": claimer},
		options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}}),
	).Decode(&claim)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errs.ErrRecordNotFound.WrapMsg("claim not found", "packetID", packetID, "claimer", claimer)
		}
		return nil, err
	}
	return &claim, nil
}

func (m *RedPacketClaimMgo) GetByPacketIDAndUserID(ctx context.Context, packetID, userID string) (*model.RedPacketClaim, error) {
	var claim model.RedPacketClaim
	err := m.coll.FindOne(ctx,
		bson.M{"packet_id": packetID, "user_id": userID},
		options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}}),
	).Decode(&claim)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errs.ErrRecordNotFound.WrapMsg("claim not found", "packetID", packetID, "userID", userID)
		}
		return nil, err
	}
	return &claim, nil
}

func (m *RedPacketClaimMgo) ListByPacketID(ctx context.Context, packetID string) ([]*model.RedPacketClaim, error) {
	cursor, err := m.coll.Find(ctx,
		bson.M{"packet_id": packetID},
		options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}),
	)
	if err != nil {
		return nil, err
	}
	var claims []*model.RedPacketClaim
	if err := cursor.All(ctx, &claims); err != nil {
		return nil, err
	}
	return claims, nil
}

// ---- RedPacketClaimAuth ----

type RedPacketClaimAuthMgo struct {
	coll *mongo.Collection
}

func NewRedPacketClaimAuthMongo(db *mongo.Database) (database.RedPacketClaimAuth, error) {
	coll := db.Collection("red_packet_claim_auth")
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "auth_nonce", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "packet_id", Value: 1}, {Key: "claimer", Value: 1}},
		},
	})
	if err != nil {
		return nil, err
	}
	return &RedPacketClaimAuthMgo{coll: coll}, nil
}

func (m *RedPacketClaimAuthMgo) Create(ctx context.Context, auth *model.RedPacketClaimAuth) error {
	_, err := m.coll.InsertOne(ctx, auth)
	return err
}

func (m *RedPacketClaimAuthMgo) Get(ctx context.Context, packetID, claimer string) (*model.RedPacketClaimAuth, error) {
	var auth model.RedPacketClaimAuth
	err := m.coll.FindOne(ctx, bson.M{
		"packet_id": packetID,
		"claimer":   claimer,
		"used":      false,
	}).Decode(&auth)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errs.ErrRecordNotFound.WrapMsg("claim auth not found", "packetID", packetID, "claimer", claimer)
		}
		return nil, err
	}
	return &auth, nil
}

func (m *RedPacketClaimAuthMgo) MarkUsed(ctx context.Context, authNonce string) error {
	res, err := m.coll.UpdateOne(ctx,
		bson.M{"auth_nonce": authNonce},
		bson.M{"$set": bson.M{"used": true}},
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errs.ErrRecordNotFound.WrapMsg("claim auth not found", "authNonce", authNonce)
	}
	return nil
}

// ---- RedPacketRefund ----

type RedPacketRefundMgo struct {
	coll *mongo.Collection
}

func NewRedPacketRefundMongo(db *mongo.Database) (database.RedPacketRefund, error) {
	coll := db.Collection("red_packet_refund")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "tx_hash", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, err
	}
	return &RedPacketRefundMgo{coll: coll}, nil
}

func (m *RedPacketRefundMgo) Save(ctx context.Context, refund *model.RedPacketRefund) error {
	_, err := m.coll.UpdateOne(ctx,
		bson.M{"tx_hash": refund.TxHash},
		bson.M{"$setOnInsert": refund},
		options.Update().SetUpsert(true),
	)
	return err
}

func (m *RedPacketRefundMgo) GetByPacketID(ctx context.Context, packetID string) (*model.RedPacketRefund, error) {
	var r model.RedPacketRefund
	err := m.coll.FindOne(ctx, bson.M{"packet_id": packetID}).Decode(&r)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errs.ErrRecordNotFound.WrapMsg("refund not found", "packetID", packetID)
		}
		return nil, err
	}
	return &r, nil
}

// ---- WalletBindingChallenge ----

type WalletBindingChallengeMgo struct {
	coll *mongo.Collection
}

func NewWalletBindingChallengeMongo(db *mongo.Database) (database.WalletBindingChallenge, error) {
	coll := db.Collection("wallet_binding_challenge")
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "challenge_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "wallet_address", Value: 1}},
		},
	})
	if err != nil {
		return nil, err
	}
	return &WalletBindingChallengeMgo{coll: coll}, nil
}

func (m *WalletBindingChallengeMgo) Create(ctx context.Context, challenge *model.WalletBindingChallenge) error {
	_, err := m.coll.InsertOne(ctx, challenge)
	return err
}

func (m *WalletBindingChallengeMgo) Get(ctx context.Context, challengeID string) (*model.WalletBindingChallenge, error) {
	var c model.WalletBindingChallenge
	err := m.coll.FindOne(ctx, bson.M{"challenge_id": challengeID}).Decode(&c)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errs.ErrRecordNotFound.WrapMsg("wallet binding challenge not found", "challengeID", challengeID)
		}
		return nil, err
	}
	return &c, nil
}

func (m *WalletBindingChallengeMgo) Update(ctx context.Context, c *model.WalletBindingChallenge) error {
	updates := bson.M{
		"status":      c.Status,
		"signature":   c.Signature,
		"verified_at": c.VerifiedAt,
		"updated_at":  c.UpdatedAt,
	}
	res, err := m.coll.UpdateOne(ctx, bson.M{"challenge_id": c.ChallengeID}, bson.M{"$set": updates})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errs.ErrRecordNotFound.WrapMsg("wallet binding challenge not found", "challengeID", c.ChallengeID)
	}
	return nil
}

// ---- WalletBinding ----

type WalletBindingMgo struct {
	coll *mongo.Collection
}

func NewWalletBindingMongo(db *mongo.Database) (database.WalletBinding, error) {
	coll := db.Collection("wallet_binding")
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "chain_type", Value: 1}, {Key: "wallet_address", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
	})
	if err != nil {
		return nil, err
	}
	return &WalletBindingMgo{coll: coll}, nil
}

// GetExpiredPending returns red packets that have expired but are still in
// "ACTIVE" status (i.e., on-chain creation confirmed, not yet fully claimed or refunded).
func (m *RedPacketMgo) GetExpiredPending(ctx context.Context, now int64) ([]*model.RedPacket, error) {
	cur, err := m.coll.Find(ctx, bson.M{
		"status":    "ACTIVE",
		"expiry_at": bson.M{"$lt": now, "$gt": 0},
	})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []*model.RedPacket
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (m *WalletBindingMgo) Upsert(ctx context.Context, b *model.WalletBinding) error {
	filter := bson.M{
		"user_id":        b.UserID,
		"chain_type":     b.ChainType,
		"wallet_address": b.WalletAddress,
	}
	updates := bson.M{
		"chain_id":     b.ChainID,
		"status":       b.Status,
		"challenge_id": b.ChallengeID,
		"verified_at":  b.VerifiedAt,
		"revoked_at":   b.RevokedAt,
		"updated_at":   b.UpdatedAt,
	}
	setOnInsert := bson.M{
		"created_at": b.CreatedAt,
	}
	_, err := m.coll.UpdateOne(ctx, filter,
		bson.M{"$set": updates, "$setOnInsert": setOnInsert},
		options.Update().SetUpsert(true),
	)
	return err
}

func (m *WalletBindingMgo) GetActive(ctx context.Context, userID, chainType, walletAddress string) (*model.WalletBinding, error) {
	var b model.WalletBinding
	err := m.coll.FindOne(ctx, bson.M{
		"user_id":        userID,
		"chain_type":     chainType,
		"wallet_address": walletAddress,
		"status":         "ACTIVE",
	}).Decode(&b)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errs.ErrRecordNotFound.WrapMsg("active wallet binding not found", "userID", userID, "chainType", chainType, "walletAddress", walletAddress)
		}
		return nil, err
	}
	return &b, nil
}

// ---- AdminAuditLog ----

type AdminAuditLogMgo struct {
	coll *mongo.Collection
}

func NewAdminAuditLogMongo(db *mongo.Database) (database.AdminAuditLog, error) {
	coll := db.Collection("admin_audit_log")
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "operator_id", Value: 1}}},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
	})
	if err != nil {
		return nil, err
	}
	return &AdminAuditLogMgo{coll: coll}, nil
}

func (m *AdminAuditLogMgo) Create(ctx context.Context, entry *model.AdminAuditLog) error {
	_, err := m.coll.InsertOne(ctx, entry)
	return err
}
