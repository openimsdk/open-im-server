package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"redpacket/internal/authctx"
	"redpacket/internal/model"
	"redpacket/internal/repository"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestService(t *testing.T) (*RedPacketService, repository.Repository) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm.Open() error = %v", err)
	}

	if err := db.AutoMigrate(
		&model.RedPacket{},
		&model.RedPacketClaim{},
		&model.RedPacketClaimAuth{},
		&model.RedPacketRefund{},
		&model.WalletBindingChallenge{},
		&model.WalletBinding{},
	); err != nil {
		t.Fatalf("AutoMigrate() error = %v", err)
	}

	repo := repository.New(db)
	svc := NewRedPacketService(repo, nil, nil, "")
	return svc, repo
}

func seedWalletBinding(t *testing.T, repo repository.Repository, userID, chainType, wallet string) {
	t.Helper()

	err := repo.UpsertWalletBinding(context.Background(), &model.WalletBinding{
		UserID:        userID,
		ChainType:     chainType,
		WalletAddress: wallet,
		Status:        "ACTIVE",
		ChallengeID:   "test-challenge",
		VerifiedAt:    time.Now(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	})
	if err != nil {
		t.Fatalf("UpsertWalletBinding() error = %v", err)
	}
}

func TestCanClaimRejectsExpiredAndAlreadyClaimed(t *testing.T) {
	svc, repo := newTestService(t)
	ctx := authctx.WithCurrentUserID(context.Background(), "u2")

	activePacket := &model.RedPacket{
		BizID:         "biz-active",
		ChainType:     "EVM",
		PacketID:      "1001",
		CreatorUserID: "u1",
		CreatorWallet: "0xabc",
		GroupID:       "g-active",
		Status:        "ACTIVE",
		ExpiryAt:      time.Now().Add(10 * time.Minute).Unix(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := repo.CreateRedPacket(ctx, activePacket); err != nil {
		t.Fatalf("CreateRedPacket(active) error = %v", err)
	}
	seedWalletBinding(t, repo, "u2", "EVM", "0xclaimer")

	claim := &model.RedPacketClaim{
		PacketID:      "1001",
		ClaimerWallet: "0xclaimer",
		ClaimTxHash:   "0xtx1",
		Status:        "CONFIRMED",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := repo.SaveClaim(ctx, claim); err != nil {
		t.Fatalf("SaveClaim() error = %v", err)
	}

	if err := svc.CanClaim(ctx, "1001", "0xclaimer", "u2"); err == nil || err.Error() != "already claimed" {
		t.Fatalf("expected already claimed error, got %v", err)
	}

	expiredPacket := &model.RedPacket{
		BizID:         "biz-expired",
		ChainType:     "EVM",
		PacketID:      "1002",
		CreatorUserID: "u1",
		CreatorWallet: "0xabc",
		GroupID:       "g-expired",
		Status:        "ACTIVE",
		ExpiryAt:      time.Now().Add(-1 * time.Minute).Unix(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := repo.CreateRedPacket(ctx, expiredPacket); err != nil {
		t.Fatalf("CreateRedPacket(expired) error = %v", err)
	}
	seedWalletBinding(t, repo, "u3", "EVM", "0xfresh")

	if err := svc.CanClaim(authctx.WithCurrentUserID(context.Background(), "u3"), "1002", "0xfresh", "u3"); err == nil || err.Error() != "packet is expired" {
		t.Fatalf("expected expired error, got %v", err)
	}
}

func TestCanClaimRejectsAlreadyClaimedByUserID(t *testing.T) {
	svc, repo := newTestService(t)
	ctx := authctx.WithCurrentUserID(context.Background(), "u2")

	packet := &model.RedPacket{
		BizID:         "biz-user-claimed",
		ChainType:     "EVM",
		PacketID:      "1003",
		CreatorUserID: "u1",
		CreatorWallet: "0xabc",
		GroupID:       "g-user-claimed",
		Status:        "ACTIVE",
		ExpiryAt:      time.Now().Add(10 * time.Minute).Unix(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := repo.CreateRedPacket(ctx, packet); err != nil {
		t.Fatalf("CreateRedPacket() error = %v", err)
	}
	seedWalletBinding(t, repo, "u2", "EVM", "0xanother-wallet")

	claim := &model.RedPacketClaim{
		PacketID:      "1003",
		UserID:        "u2",
		ClaimerWallet: "0xclaimer",
		ClaimTxHash:   "0xtx-user-claimed",
		Status:        "CONFIRMED",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := repo.SaveClaim(ctx, claim); err != nil {
		t.Fatalf("SaveClaim() error = %v", err)
	}

	if err := svc.CanClaim(ctx, "1003", "0xanother-wallet", "u2"); err == nil || err.Error() != "user already claimed" {
		t.Fatalf("expected user already claimed error, got %v", err)
	}
}

func TestCanClaimUsesPacketTypeRules(t *testing.T) {
	svc, repo := newTestService(t)
	ctx := authctx.WithCurrentUserID(context.Background(), "u2")

	groupPacket := &model.RedPacket{
		BizID:         "biz-group",
		ChainType:     "EVM",
		PacketID:      "1101",
		CreatorUserID: "u1",
		CreatorWallet: "0xabc",
		PacketType:    0,
		Status:        "ACTIVE",
		ExpiryAt:      time.Now().Add(10 * time.Minute).Unix(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := repo.CreateRedPacket(ctx, groupPacket); err != nil {
		t.Fatalf("CreateRedPacket(group) error = %v", err)
	}
	seedWalletBinding(t, repo, "u2", "EVM", "0xclaimer")
	if err := svc.CanClaim(ctx, "1101", "0xclaimer", "u2"); err == nil || err.Error() != "group_id is required for fixed packet claim" {
		t.Fatalf("expected missing group_id error, got %v", err)
	}

	transferPacket := &model.RedPacket{
		BizID:          "biz-transfer",
		ChainType:      "EVM",
		PacketID:       "1102",
		CreatorUserID:  "u1",
		CreatorWallet:  "0xabc",
		PacketType:     2,
		ReceiverUserID: "u-receiver",
		Status:         "ACTIVE",
		ExpiryAt:       time.Now().Add(10 * time.Minute).Unix(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := repo.CreateRedPacket(ctx, transferPacket); err != nil {
		t.Fatalf("CreateRedPacket(transfer) error = %v", err)
	}
	seedWalletBinding(t, repo, "u-other", "EVM", "0xclaimer")
	if err := svc.CanClaim(ctx, "1102", "0xclaimer", "u-other"); err == nil || err.Error() != "user is not the designated receiver" {
		t.Fatalf("expected designated receiver error, got %v", err)
	}
}

func TestCreateOrderPersistsScopeFields(t *testing.T) {
	svc, repo := newTestService(t)
	ctx := authctx.WithCurrentUserID(context.Background(), "u-create")

	result, err := svc.CreateOrder(ctx, &CreateOrderRequest{
		ChainType:       "EVM",
		CreatorWallet:   "0x1111111111111111111111111111111111111111",
		GroupID:         "g-100",
		ScopeType:       "group",
		PacketType:      1,
		Token:           "0x2222222222222222222222222222222222222222",
		TotalAmount:     "1000",
		TotalShares:     10,
		ReceiverUserIDs: []string{"u2", "u3"},
	})
	if err != nil {
		t.Fatalf("CreateOrder() error = %v", err)
	}

	bizID, _ := result["biz_id"].(string)
	record, err := repo.GetRedPacketByBizID(ctx, bizID)
	if err != nil {
		t.Fatalf("GetRedPacketByBizID() error = %v", err)
	}

	if record.ScopeType != "GROUP" {
		t.Fatalf("scope type mismatch: got %s", record.ScopeType)
	}
	if record.ChainType != "EVM" {
		t.Fatalf("chain type mismatch: got %s", record.ChainType)
	}
	if record.GroupID != "g-100" {
		t.Fatalf("group id mismatch: got %s", record.GroupID)
	}

	var got []string
	if err := json.Unmarshal([]byte(record.ReceiverUserIDs), &got); err != nil {
		t.Fatalf("Unmarshal(receiver_user_ids) error = %v", err)
	}
	if len(got) != 2 || got[0] != "u2" || got[1] != "u3" {
		t.Fatalf("receiver_user_ids mismatch: got %+v", got)
	}
}

func TestCreatedCallbackUpdatesBindingAndScope(t *testing.T) {
	svc, repo := newTestService(t)
	ctx := authctx.WithCurrentUserID(context.Background(), "u-create")

	result, err := svc.CreateOrder(ctx, &CreateOrderRequest{
		ChainType:     "TRON",
		CreatorWallet: "0x1111111111111111111111111111111111111111",
		PacketType:    2,
		Token:         "0x0000000000000000000000000000000000000000",
		TotalAmount:   "1000",
		TotalShares:   1,
	})
	if err != nil {
		t.Fatalf("CreateOrder() error = %v", err)
	}

	bizID, _ := result["biz_id"].(string)
	err = svc.CreatedCallback(ctx, &CreatedCallbackRequest{
		BizID:          bizID,
		TxHash:         "0xabc123",
		PacketID:       "3001",
		ScopeType:      "DIRECT",
		ReceiverUserID: "u-receiver",
	})
	if err != nil {
		t.Fatalf("CreatedCallback() error = %v", err)
	}

	record, err := repo.GetRedPacketByBizID(ctx, bizID)
	if err != nil {
		t.Fatalf("GetRedPacketByBizID() error = %v", err)
	}

	if record.PacketID != "3001" {
		t.Fatalf("packet id mismatch: got %s", record.PacketID)
	}
	if record.ChainType != "TRON" {
		t.Fatalf("chain type mismatch: got %s", record.ChainType)
	}
	if record.TxHash != "0xabc123" {
		t.Fatalf("tx hash mismatch: got %s", record.TxHash)
	}
	if record.Status != "ACTIVE" {
		t.Fatalf("status mismatch: got %s", record.Status)
	}
	if record.ScopeType != "DIRECT" {
		t.Fatalf("scope type mismatch: got %s", record.ScopeType)
	}
	if record.ReceiverUserID != "u-receiver" {
		t.Fatalf("receiver user mismatch: got %s", record.ReceiverUserID)
	}
}

func TestIssueClaimSignValidatesInputsAndPersistsAuth(t *testing.T) {
	svc, repo := newTestService(t)
	ctx := authctx.WithCurrentUserID(context.Background(), "u2")

	packet := &model.RedPacket{
		BizID:         "biz-sign",
		ChainType:     "EVM",
		PacketID:      "2001",
		CreatorUserID: "u1",
		CreatorWallet: "0xabc",
		GroupID:       "g-sign",
		Status:        "ACTIVE",
		ExpiryAt:      time.Now().Add(10 * time.Minute).Unix(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := repo.CreateRedPacket(ctx, packet); err != nil {
		t.Fatalf("CreateRedPacket() error = %v", err)
	}
	seedWalletBinding(t, repo, "u2", "EVM", "0x1111111111111111111111111111111111111111")

	if _, err := svc.IssueClaimSign(ctx, "bad-packet-id", "0x1111111111111111111111111111111111111111", "0"); err == nil {
		t.Fatalf("expected invalid packet id error")
	}

	result, err := svc.IssueClaimSign(ctx, "2001", "0x1111111111111111111111111111111111111111", "123")
	if err != nil {
		t.Fatalf("IssueClaimSign() error = %v", err)
	}

	auth, err := repo.GetClaimAuth(ctx, "2001", "0x1111111111111111111111111111111111111111")
	if err != nil {
		t.Fatalf("GetClaimAuth() error = %v", err)
	}
	if auth.AuthNonce == "" {
		t.Fatalf("expected auth nonce to be persisted")
	}
	if auth.RandomSeed != "123" {
		t.Fatalf("random seed mismatch: got %s", auth.RandomSeed)
	}
	if result["auth_nonce"] == "" {
		t.Fatalf("expected auth_nonce in response")
	}
}

func TestClaimResultPersistsPendingWithoutChainParser(t *testing.T) {
	svc, repo := newTestService(t)
	ctx := authctx.WithCurrentUserID(context.Background(), "u2")

	packet := &model.RedPacket{
		BizID:         "biz-claim-result",
		ChainType:     "EVM",
		PacketID:      "2101",
		CreatorUserID: "u1",
		CreatorWallet: "0xabc",
		GroupID:       "g-1",
		PacketType:    0,
		Status:        "ACTIVE",
		ExpiryAt:      time.Now().Add(10 * time.Minute).Unix(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := repo.CreateRedPacket(ctx, packet); err != nil {
		t.Fatalf("CreateRedPacket() error = %v", err)
	}
	seedWalletBinding(t, repo, "u2", "EVM", "0x1111111111111111111111111111111111111111")

	err := svc.ClaimResult(ctx, &ClaimResultRequest{
		PacketID: "2101",
		Claimer:  "0x1111111111111111111111111111111111111111",
		TxHash:   "0xtx-claim",
	})
	if err != nil {
		t.Fatalf("ClaimResult() error = %v", err)
	}

	claim, err := repo.GetClaimByPacketIDAndClaimer(ctx, "2101", "0x1111111111111111111111111111111111111111")
	if err != nil {
		t.Fatalf("GetClaimByPacketIDAndClaimer() error = %v", err)
	}
	if claim.Status != "PENDING" {
		t.Fatalf("claim status mismatch: got %s", claim.Status)
	}
	if claim.UserID != "u2" {
		t.Fatalf("user id mismatch: got %s", claim.UserID)
	}
}
