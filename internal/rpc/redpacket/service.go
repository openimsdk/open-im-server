package redpacket

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/openimsdk/open-im-server/v3/internal/rpc/redpacket/chain"
	"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	pbredpacket "github.com/openimsdk/protocol/redpacket"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
)

func (s *redPacketServer) CreateOrder(ctx context.Context, req *pbredpacket.CreateOrderReq) (*pbredpacket.CreateOrderResp, error) {
	currentUserID := mcontext.GetOpUserID(ctx)
	if currentUserID == "" {
		return nil, servererrs.ErrNoPermission.WrapMsg("op user id is empty")
	}

	bizID := uuid.NewString()
	chainType, err := normalizeChainType(req.ChainType)
	if err != nil {
		return nil, err
	}
	scopeType := normalizeScopeType(req.ScopeType)
	if err := validateCreateScope(scopeType, req.GroupID, req.ReceiverUserID, req.ReceiverUserIDs); err != nil {
		return nil, err
	}
	if err := s.validateCreateHook(ctx, req); err != nil {
		return nil, err
	}

	chainID := req.ChainID
	contractAddress := strings.TrimSpace(req.ContractAddress)
	if chainType == "EVM" && s.chainClient != nil {
		if chainID == 0 {
			if chainValue := s.chainClient.ChainID(); chainValue != nil {
				chainID = chainValue.Int64()
			}
		}
		if contractAddress == "" {
			contractAddress = s.chainClient.ContractAddress().Hex()
		}
	}
	if chainType == "TRON" && s.tronClient != nil && contractAddress == "" {
		contractAddress = s.tronClient.ContractAddress()
	}

	rp := &model.RedPacket{
		BizID:           bizID,
		ChainType:       chainType,
		ChainID:         chainID,
		ContractAddress: contractAddress,
		CreatorUserID:   currentUserID,
		CreatorWallet:   req.CreatorWallet,
		GroupID:         req.GroupID,
		ScopeType:       scopeType,
		ReceiverUserID:  req.ReceiverUserID,
		ReceiverUserIDs: append([]string(nil), req.ReceiverUserIDs...),
		PacketType:      req.PacketType,
		Token:           req.Token,
		TotalAmount:     req.TotalAmount,
		TotalShares:     req.TotalShares,
		ExpiryAt:        req.ExpiryAt,
		Status:          "PENDING",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.db.CreateRedPacket(ctx, rp); err != nil {
		log.ZError(ctx, "create redpacket failed", err, "bizID", bizID)
		return nil, servererrs.ErrDatabase.WrapMsg("failed to create red packet")
	}

	return &pbredpacket.CreateOrderResp{BizID: bizID}, nil
}

func (s *redPacketServer) CreatedCallback(ctx context.Context, req *pbredpacket.CreatedCallbackReq) (*pbredpacket.CreatedCallbackResp, error) {
	if strings.TrimSpace(req.BizID) == "" || strings.TrimSpace(req.TxHash) == "" {
		return nil, errs.ErrArgs.WrapMsg("biz_id and tx_hash are required")
	}

	rp, err := s.db.GetRedPacketByBizID(ctx, req.BizID)
	if err != nil {
		return nil, err
	}

	groupID := firstNonEmpty(req.GroupID, rp.GroupID)
	scopeType := normalizeScopeType(firstNonEmpty(req.ScopeType, rp.ScopeType))
	receiverUserID := firstNonEmpty(req.ReceiverUserID, rp.ReceiverUserID)
	receiverUserIDs := rp.ReceiverUserIDs
	if len(req.ReceiverUserIDs) > 0 {
		receiverUserIDs = append([]string(nil), req.ReceiverUserIDs...)
	}

	if err := validateCreateScope(scopeType, groupID, receiverUserID, receiverUserIDs); err != nil {
		return nil, err
	}

	createdPacket, err := s.resolveCreatedPacket(ctx, rp, req.TxHash, req.PacketID)
	if err != nil {
		return nil, err
	}

	if err := s.db.UpdateRedPacketCreated(ctx, &model.RedPacket{
		BizID:           req.BizID,
		ChainType:       rp.ChainType,
		PacketID:        createdPacket.PacketID,
		ChainID:         createdPacket.ChainID,
		ContractAddress: createdPacket.ContractAddress,
		TxHash:          req.TxHash,
		GroupID:         groupID,
		ScopeType:       scopeType,
		ReceiverUserID:  receiverUserID,
		ReceiverUserIDs: receiverUserIDs,
		Status:          "ACTIVE",
	}); err != nil {
		return nil, err
	}
	return &pbredpacket.CreatedCallbackResp{}, nil
}

func (s *redPacketServer) GetDetail(ctx context.Context, req *pbredpacket.GetDetailReq) (*pbredpacket.GetDetailResp, error) {
	if strings.TrimSpace(req.PacketID) == "" {
		return nil, errs.ErrArgs.WrapMsg("packet_id is required")
	}

	rp, err := s.db.GetRedPacketByPacketID(ctx, req.PacketID)
	if err != nil {
		return nil, err
	}
	claims, err := s.db.GetClaimsByPacketID(ctx, req.PacketID)
	if err != nil {
		claims = nil
	}

	return &pbredpacket.GetDetailResp{
		Record: redPacketModelToProto(rp),
		Claims: claimsModelToProto(claims),
	}, nil
}

func (s *redPacketServer) IssueClaimSign(ctx context.Context, req *pbredpacket.IssueClaimSignReq) (*pbredpacket.IssueClaimSignResp, error) {
	currentUserID := mcontext.GetOpUserID(ctx)
	if currentUserID == "" {
		return nil, servererrs.ErrNoPermission.WrapMsg("op user id is empty")
	}
	if strings.TrimSpace(req.PacketID) == "" || strings.TrimSpace(req.Claimer) == "" {
		return nil, errs.ErrArgs.WrapMsg("packet_id and claimer are required")
	}
	if err := s.canClaim(ctx, req.PacketID, req.Claimer, currentUserID); err != nil {
		return nil, err
	}

	packetIDBig := new(big.Int)
	if _, ok := packetIDBig.SetString(req.PacketID, 10); !ok {
		return nil, errs.ErrArgs.WrapMsg("invalid packet_id", "packetID", req.PacketID)
	}

	claimerAddr := common.HexToAddress(req.Claimer)
	nonce := fmt.Sprintf("%d", time.Now().UnixNano())
	authNonceBig := new(big.Int)
	authNonceBig.SetString(nonce, 10)
	deadline := time.Now().Add(5 * time.Minute).Unix()
	randomSeedBig := new(big.Int)
	if req.RandomSeed != "" && req.RandomSeed != "0" {
		if _, ok := randomSeedBig.SetString(req.RandomSeed, 10); !ok {
			return nil, errs.ErrArgs.WrapMsg("invalid random_seed", "randomSeed", req.RandomSeed)
		}
	} else {
		randomSeedBig.SetInt64(time.Now().UnixNano())
	}
	deadlineBig := big.NewInt(deadline)

	var digest [32]byte
	var err error
	if s.chainClient != nil {
		digest, err = s.chainClient.GetSignMessage(ctx, packetIDBig, claimerAddr, authNonceBig, randomSeedBig, deadlineBig)
		if err != nil {
			return nil, errs.ErrInternalServer.WrapMsg("getSignMessage failed: " + err.Error())
		}
	} else {
		digest = crypto.Keccak256Hash([]byte(fmt.Sprintf("%s:%s:%s:%s:%d", req.PacketID, req.Claimer, nonce, randomSeedBig.String(), deadline)))
	}

	var signature []byte
	if s.signerKey != nil {
		signature, err = crypto.Sign(digest[:], s.signerKey)
		if err != nil {
			return nil, errs.ErrInternalServer.WrapMsg("sign failed: " + err.Error())
		}
		if len(signature) == 65 && signature[64] < 27 {
			signature[64] += 27
		}
	} else {
		signature = []byte("0xplaceholder-signature-for-testing")
	}

	sigHex := "0x" + hex.EncodeToString(signature)

	auth := &model.RedPacketClaimAuth{
		PacketID:   req.PacketID,
		Claimer:    req.Claimer,
		AuthNonce:  nonce,
		RandomSeed: randomSeedBig.String(),
		Deadline:   deadline,
		Signature:  sigHex,
		CreatedAt:  time.Now(),
	}

	if err := s.db.CreateClaimAuth(ctx, auth); err != nil {
		return nil, servererrs.ErrDatabase.WrapMsg("save claim auth failed: " + err.Error())
	}

	return &pbredpacket.IssueClaimSignResp{
		AuthNonce:  nonce,
		Deadline:   deadline,
		Signature:  sigHex,
		RandomSeed: randomSeedBig.String(),
	}, nil
}

func (s *redPacketServer) ClaimResult(ctx context.Context, req *pbredpacket.ClaimResultReq) (*pbredpacket.ClaimResultResp, error) {
	currentUserID := mcontext.GetOpUserID(ctx)
	if currentUserID == "" {
		return nil, servererrs.ErrNoPermission.WrapMsg("op user id is empty")
	}
	if strings.TrimSpace(req.PacketID) == "" || strings.TrimSpace(req.Claimer) == "" || strings.TrimSpace(req.TxHash) == "" {
		return nil, errs.ErrArgs.WrapMsg("packet_id, claimer and tx_hash are required")
	}

	rp, err := s.db.GetRedPacketByPacketID(ctx, req.PacketID)
	if err != nil {
		return nil, err
	}

	if err := validateClaimBase(rp, currentUserID, req.Claimer); err != nil {
		return nil, err
	}

	claim := &model.RedPacketClaim{
		PacketID:      req.PacketID,
		UserID:        currentUserID,
		ClaimerWallet: req.Claimer,
		ClaimTxHash:   req.TxHash,
		Status:        "PENDING",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.db.SaveClaim(ctx, claim); err != nil {
		return nil, err
	}

	claimedEvent, err := s.resolveClaimedEvent(ctx, rp, req.TxHash)
	if err != nil {
		log.ZWarn(ctx, "resolve claim event failed", err, "txHash", req.TxHash)
		return &pbredpacket.ClaimResultResp{}, nil
	}
	if claimedEvent == nil {
		return &pbredpacket.ClaimResultResp{}, nil
	}
	if !strings.EqualFold(claimedEvent.ClaimerWallet, req.Claimer) {
		return nil, errs.ErrArgs.WrapMsg(fmt.Sprintf("claim event claimer mismatch: got %s want %s", claimedEvent.ClaimerWallet, req.Claimer))
	}

	confirmed := &model.RedPacketClaim{
		PacketID:      req.PacketID,
		UserID:        currentUserID,
		ClaimerWallet: claimedEvent.ClaimerWallet,
		AuthNonce:     claimedEvent.AuthNonce,
		ClaimTxHash:   req.TxHash,
		ClaimedAmount: claimedEvent.Amount,
		BlockNumber:   claimedEvent.BlockNumber,
		Status:        "CONFIRMED",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := s.db.SaveClaim(ctx, confirmed); err != nil {
		return nil, err
	}

	if claimedEvent.AuthNonce != "" {
		if err := s.db.MarkClaimAuthUsed(ctx, claimedEvent.AuthNonce); err != nil {
			log.ZWarn(ctx, "mark claim auth used failed", err, "authNonce", claimedEvent.AuthNonce)
		}
	}

	nextStatus := derivePacketStatusAfterClaim(rp, claimedEvent.Amount)
	if err := s.db.UpdateRedPacketClaimProgress(ctx, req.PacketID, claimedEvent.Amount, nextStatus); err != nil {
		return nil, err
	}
	return &pbredpacket.ClaimResultResp{}, nil
}

// canClaim runs the claim-eligibility check (formerly RedPacketService.CanClaim).
func (s *redPacketServer) canClaim(ctx context.Context, packetID, claimer, userID string) error {
	rp, err := s.db.GetRedPacketByPacketID(ctx, packetID)
	if err != nil {
		return err
	}

	if err := validateClaimBase(rp, userID, claimer); err != nil {
		return err
	}
	if err := s.ensureWalletBinding(ctx, userID, claimer, rp.ChainType); err != nil {
		return err
	}

	switch rp.PacketType {
	case 0:
		return s.validateFixedPacketClaim(ctx, rp, userID, claimer)
	case 1:
		return s.validateRandomPacketClaim(ctx, rp, userID, claimer)
	case 2:
		return s.validateTransferPacketClaim(ctx, rp, userID, claimer)
	default:
		return errs.ErrArgs.WrapMsg(fmt.Sprintf("unsupported packet_type: %d", rp.PacketType))
	}
}

type claimedEventSnapshot struct {
	ClaimerWallet string
	AuthNonce     string
	Amount        string
	BlockNumber   uint64
}

type createdPacketSnapshot struct {
	PacketID        string
	ChainID         int64
	ContractAddress string
	CreatorWallet   string
	PacketType      int32
	Token           string
	TotalAmount     string
	TotalShares     int32
	ExpiryAt        int64
}

func (s *redPacketServer) resolveCreatedPacket(ctx context.Context, rp *model.RedPacket, txHashHex, fallbackPacketID string) (*createdPacketSnapshot, error) {
	switch rp.ChainType {
	case "EVM":
		if s.chainClient == nil {
			if fallbackPacketID == "" {
				return nil, errs.ErrArgs.WrapMsg("packet_id is required when EVM client is unavailable")
			}
			return buildFallbackCreatedPacket(rp, fallbackPacketID), nil
		}

		events, err := s.chainClient.ParseTransactionReceipt(ctx, common.HexToHash(txHashHex))
		if err != nil {
			if fallbackPacketID == "" {
				return nil, errs.ErrInternalServer.WrapMsg("parse created tx failed: " + err.Error())
			}
			return buildFallbackCreatedPacket(rp, fallbackPacketID), nil
		}

		for _, event := range events {
			if event.Name != "PacketCreated" {
				continue
			}
			createdPacket := buildCreatedPacketSnapshot(rp, event)
			if chainValue := s.chainClient.ChainID(); chainValue != nil {
				createdPacket.ChainID = chainValue.Int64()
			}
			createdPacket.ContractAddress = s.chainClient.ContractAddress().Hex()
			if err := validateCreatedPacket(rp, createdPacket); err != nil {
				return nil, err
			}
			return createdPacket, nil
		}

		if fallbackPacketID == "" {
			return nil, errs.ErrInternalServer.WrapMsg("PacketCreated event not found in tx: " + txHashHex)
		}
		return buildFallbackCreatedPacket(rp, fallbackPacketID), nil
	case "TRON":
		if s.tronClient == nil {
			if fallbackPacketID == "" {
				return nil, errs.ErrArgs.WrapMsg("packet_id is required when TRON client is unavailable")
			}
			return buildFallbackCreatedPacket(rp, fallbackPacketID), nil
		}

		events, err := s.tronClient.ParseTransactionReceipt(ctx, txHashHex)
		if err != nil {
			if fallbackPacketID == "" {
				return nil, errs.ErrInternalServer.WrapMsg("parse tron created tx failed: " + err.Error())
			}
			return buildFallbackCreatedPacket(rp, fallbackPacketID), nil
		}

		for _, event := range events {
			if event.Name != "PacketCreated" {
				continue
			}
			createdPacket := buildCreatedPacketSnapshot(rp, event)
			createdPacket.ContractAddress = firstNonEmpty(s.tronClient.ContractAddress(), rp.ContractAddress)
			if err := validateCreatedPacket(rp, createdPacket); err != nil {
				return nil, err
			}
			return createdPacket, nil
		}

		if fallbackPacketID == "" {
			return nil, errs.ErrInternalServer.WrapMsg("PacketCreated event not found in TRON tx: " + txHashHex)
		}
		return buildFallbackCreatedPacket(rp, fallbackPacketID), nil
	default:
		return nil, errs.ErrArgs.WrapMsg("unsupported chain_type: " + rp.ChainType)
	}
}

// validateCreateHook reserves a centralized validation extension point split by packet type.
func (s *redPacketServer) validateCreateHook(ctx context.Context, req *pbredpacket.CreateOrderReq) error {
	switch req.PacketType {
	case 0:
		return s.validateFixedPacketCreate(ctx, req)
	case 1:
		return s.validateRandomPacketCreate(ctx, req)
	case 2:
		return s.validateTransferPacketCreate(ctx, req)
	default:
		return errs.ErrArgs.WrapMsg(fmt.Sprintf("unsupported packet_type: %d", req.PacketType))
	}
}

func (s *redPacketServer) validateFixedPacketCreate(ctx context.Context, req *pbredpacket.CreateOrderReq) error {
	return nil
}

func (s *redPacketServer) validateRandomPacketCreate(ctx context.Context, req *pbredpacket.CreateOrderReq) error {
	return nil
}

func (s *redPacketServer) validateTransferPacketCreate(ctx context.Context, req *pbredpacket.CreateOrderReq) error {
	return nil
}

func buildFallbackCreatedPacket(rp *model.RedPacket, packetID string) *createdPacketSnapshot {
	return &createdPacketSnapshot{
		PacketID:        packetID,
		ChainID:         rp.ChainID,
		ContractAddress: rp.ContractAddress,
		CreatorWallet:   strings.ToLower(rp.CreatorWallet),
		PacketType:      rp.PacketType,
		Token:           normalizeTokenAddress(rp.Token),
		TotalAmount:     rp.TotalAmount,
		TotalShares:     rp.TotalShares,
		ExpiryAt:        rp.ExpiryAt,
	}
}

func buildCreatedPacketSnapshot(rp *model.RedPacket, event *chain.ParsedEvent) *createdPacketSnapshot {
	return &createdPacketSnapshot{
		PacketID:        chain.GetPacketIDFromEvent(event).String(),
		ChainID:         rp.ChainID,
		ContractAddress: rp.ContractAddress,
		CreatorWallet:   strings.ToLower(chain.GetAddressFromEvent(event, "creator").Hex()),
		PacketType:      int32(chain.GetUintFromEvent(event, "packetType").Int64()),
		Token:           strings.ToLower(chain.GetAddressFromEvent(event, "token").Hex()),
		TotalAmount:     chain.GetUintFromEvent(event, "totalAmount").String(),
		TotalShares:     int32(chain.GetUintFromEvent(event, "totalShares").Int64()),
		ExpiryAt:        chain.GetUintFromEvent(event, "expiryAt").Int64(),
	}
}

func validateCreatedPacket(rp *model.RedPacket, createdPacket *createdPacketSnapshot) error {
	if createdPacket == nil {
		return errs.ErrInternalServer.WrapMsg("created packet is nil")
	}
	if createdPacket.CreatorWallet != "" && strings.ToLower(rp.CreatorWallet) != createdPacket.CreatorWallet {
		return errs.ErrArgs.WrapMsg(fmt.Sprintf("creator mismatch: got %s want %s", createdPacket.CreatorWallet, rp.CreatorWallet))
	}
	if createdPacket.PacketType != rp.PacketType {
		return errs.ErrArgs.WrapMsg(fmt.Sprintf("packet type mismatch: got %d want %d", createdPacket.PacketType, rp.PacketType))
	}
	if createdPacket.TotalAmount != rp.TotalAmount {
		return errs.ErrArgs.WrapMsg(fmt.Sprintf("total amount mismatch: got %s want %s", createdPacket.TotalAmount, rp.TotalAmount))
	}
	if createdPacket.TotalShares != rp.TotalShares {
		return errs.ErrArgs.WrapMsg(fmt.Sprintf("total shares mismatch: got %d want %d", createdPacket.TotalShares, rp.TotalShares))
	}
	expectedToken := normalizeTokenAddress(rp.Token)
	if createdPacket.Token != expectedToken {
		return errs.ErrArgs.WrapMsg(fmt.Sprintf("token mismatch: got %s want %s", createdPacket.Token, expectedToken))
	}
	if rp.ExpiryAt > 0 && createdPacket.ExpiryAt != rp.ExpiryAt {
		return errs.ErrArgs.WrapMsg(fmt.Sprintf("expiry mismatch: got %d want %d", createdPacket.ExpiryAt, rp.ExpiryAt))
	}
	return nil
}

func validateClaimBase(rp *model.RedPacket, userID, claimer string) error {
	if rp == nil {
		return servererrs.ErrRecordNotFound.WrapMsg("packet not found")
	}
	if strings.TrimSpace(userID) == "" {
		return errs.ErrArgs.WrapMsg("user_id is required")
	}
	if strings.TrimSpace(claimer) == "" {
		return errs.ErrArgs.WrapMsg("claimer is required")
	}
	if rp.Status != "ACTIVE" {
		return errs.ErrArgs.WrapMsg("packet is not active, current status: " + rp.Status)
	}
	if rp.ExpiryAt > 0 && rp.ExpiryAt <= time.Now().Unix() {
		return errs.ErrArgs.WrapMsg("packet is expired")
	}
	if rp.Status == "REFUNDED" {
		return errs.ErrArgs.WrapMsg("packet is refunded")
	}
	return nil
}

func (s *redPacketServer) validateFixedPacketClaim(ctx context.Context, rp *model.RedPacket, userID, claimer string) error {
	if strings.TrimSpace(rp.GroupID) == "" {
		return errs.ErrArgs.WrapMsg("group_id is required for fixed packet claim")
	}
	if err := s.ensureNotClaimed(ctx, rp.PacketID, userID, claimer); err != nil {
		return err
	}
	return s.ensureGroupEligibility(ctx, rp.GroupID, userID)
}

func (s *redPacketServer) validateRandomPacketClaim(ctx context.Context, rp *model.RedPacket, userID, claimer string) error {
	if strings.TrimSpace(rp.GroupID) == "" {
		return errs.ErrArgs.WrapMsg("group_id is required for random packet claim")
	}
	if err := s.ensureNotClaimed(ctx, rp.PacketID, userID, claimer); err != nil {
		return err
	}
	return s.ensureGroupEligibility(ctx, rp.GroupID, userID)
}

func (s *redPacketServer) validateTransferPacketClaim(ctx context.Context, rp *model.RedPacket, userID, claimer string) error {
	if err := s.ensureNotClaimed(ctx, rp.PacketID, userID, claimer); err != nil {
		return err
	}
	if strings.TrimSpace(rp.ReceiverUserID) == "" {
		return errs.ErrArgs.WrapMsg("receiver_user_id is required for transfer claim")
	}
	if rp.ReceiverUserID != userID {
		return errs.ErrNoPermission.WrapMsg("user is not the designated receiver")
	}
	return s.ensureFriendRelationship(ctx, rp.CreatorUserID, userID)
}

func (s *redPacketServer) ensureNotClaimed(ctx context.Context, packetID, userID, claimer string) error {
	if strings.TrimSpace(userID) != "" {
		claim, err := s.db.GetClaimByPacketIDAndUserID(ctx, packetID, userID)
		if err == nil && claim != nil && claim.Status != "FAILED" {
			return errs.ErrArgs.WrapMsg("user already claimed")
		}
		if err != nil && !errs.ErrRecordNotFound.Is(err) {
			return err
		}
	}

	claim, err := s.db.GetClaimByPacketIDAndClaimer(ctx, packetID, claimer)
	if err == nil && claim != nil && claim.Status != "FAILED" {
		return errs.ErrArgs.WrapMsg("already claimed")
	}
	if err != nil && !errs.ErrRecordNotFound.Is(err) {
		return err
	}
	return nil
}

func (s *redPacketServer) ensureWalletBinding(ctx context.Context, userID, claimer, chainType string) error {
	if _, err := s.db.GetActiveWalletBinding(ctx, userID, chainType, claimer); err != nil {
		if errs.ErrRecordNotFound.Is(err) {
			return errs.ErrNoPermission.WrapMsg("wallet is not bound to user")
		}
		return err
	}
	return nil
}

// ensureGroupEligibility reserves centralized group membership checks.
func (s *redPacketServer) ensureGroupEligibility(ctx context.Context, groupID, userID string) error {
	return nil
}

// ensureFriendRelationship reserves centralized relation validation for transfer packets.
func (s *redPacketServer) ensureFriendRelationship(ctx context.Context, creatorUserID, receiverUserID string) error {
	return nil
}

func (s *redPacketServer) resolveClaimedEvent(ctx context.Context, rp *model.RedPacket, txHash string) (*claimedEventSnapshot, error) {
	var (
		events []*chain.ParsedEvent
		err    error
	)

	switch rp.ChainType {
	case "EVM":
		if s.chainClient == nil {
			return nil, nil
		}
		events, err = s.chainClient.ParseTransactionReceipt(ctx, common.HexToHash(txHash))
	case "TRON":
		if s.tronClient == nil {
			return nil, nil
		}
		events, err = s.tronClient.ParseTransactionReceipt(ctx, txHash)
	default:
		return nil, errs.ErrArgs.WrapMsg("unsupported chain_type: " + rp.ChainType)
	}
	if err != nil {
		return nil, err
	}

	for _, event := range events {
		if event.Name != "PacketClaimed" {
			continue
		}
		packetID := chain.GetPacketIDFromEvent(event).String()
		claimerWallet := strings.ToLower(chain.GetAddressFromEvent(event, "claimer").Hex())
		if packetID != rp.PacketID {
			return nil, errs.ErrArgs.WrapMsg(fmt.Sprintf("claim event packet mismatch: got %s want %s", packetID, rp.PacketID))
		}
		return &claimedEventSnapshot{
			ClaimerWallet: claimerWallet,
			AuthNonce:     chain.GetUintFromEvent(event, "authNonce").String(),
			Amount:        chain.GetAmountFromEvent(event).String(),
			BlockNumber:   event.BlockNumber,
		}, nil
	}

	return nil, nil
}

func derivePacketStatusAfterClaim(rp *model.RedPacket, claimedAmount string) string {
	if rp == nil {
		return ""
	}
	if rp.PacketType == 2 {
		return "COMPLETED"
	}

	nextShares := rp.ClaimedShares + 1
	if rp.TotalShares > 0 && nextShares >= rp.TotalShares {
		return "COMPLETED"
	}

	totalClaimed := addNumericStrings(rp.ClaimedAmount, claimedAmount)
	if rp.TotalAmount != "" && totalClaimed == rp.TotalAmount {
		return "COMPLETED"
	}

	return "ACTIVE"
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

func normalizeScopeType(scopeType string) string {
	switch strings.ToUpper(strings.TrimSpace(scopeType)) {
	case "GROUP", "DIRECT", "PUBLIC":
		return strings.ToUpper(strings.TrimSpace(scopeType))
	default:
		return "PUBLIC"
	}
}

func normalizeChainType(chainType string) (string, error) {
	switch strings.ToUpper(strings.TrimSpace(chainType)) {
	case "EVM":
		return "EVM", nil
	case "TRON":
		return "TRON", nil
	default:
		return "", errs.ErrArgs.WrapMsg("unsupported chain_type: " + chainType)
	}
}

func validateCreateScope(scopeType, groupID, receiverUserID string, receiverUserIDs []string) error {
	switch scopeType {
	case "GROUP":
		if strings.TrimSpace(groupID) == "" {
			return errs.ErrArgs.WrapMsg("group_id is required when scope_type=GROUP")
		}
	case "DIRECT":
		if strings.TrimSpace(receiverUserID) == "" && len(receiverUserIDs) == 0 {
			return errs.ErrArgs.WrapMsg("receiver_user_id or receiver_user_ids is required when scope_type=DIRECT")
		}
	}
	return nil
}

func normalizeTokenAddress(token string) string {
	if strings.TrimSpace(token) == "" {
		return strings.ToLower(common.Address{}.Hex())
	}
	return strings.ToLower(common.HexToAddress(token).Hex())
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func redPacketModelToProto(rp *model.RedPacket) *pbredpacket.RedPacketRecord {
	if rp == nil {
		return nil
	}
	return &pbredpacket.RedPacketRecord{
		BizID:           rp.BizID,
		ChainType:       rp.ChainType,
		PacketID:        rp.PacketID,
		ChainID:         rp.ChainID,
		ContractAddress: rp.ContractAddress,
		CreatorUserID:   rp.CreatorUserID,
		CreatorWallet:   rp.CreatorWallet,
		GroupID:         rp.GroupID,
		ScopeType:       rp.ScopeType,
		ReceiverUserID:  rp.ReceiverUserID,
		ReceiverUserIDs: append([]string(nil), rp.ReceiverUserIDs...),
		PacketType:      rp.PacketType,
		Token:           rp.Token,
		TotalAmount:     rp.TotalAmount,
		TotalShares:     rp.TotalShares,
		ClaimedAmount:   rp.ClaimedAmount,
		ClaimedShares:   rp.ClaimedShares,
		ExpiryAt:        rp.ExpiryAt,
		TxHash:          rp.TxHash,
		Status:          rp.Status,
		CreatedAt:       rp.CreatedAt.Unix(),
		UpdatedAt:       rp.UpdatedAt.Unix(),
	}
}

func claimsModelToProto(claims []*model.RedPacketClaim) []*pbredpacket.RedPacketClaimRecord {
	out := make([]*pbredpacket.RedPacketClaimRecord, 0, len(claims))
	for _, c := range claims {
		if c == nil {
			continue
		}
		out = append(out, &pbredpacket.RedPacketClaimRecord{
			PacketID:      c.PacketID,
			UserID:        c.UserID,
			ClaimerWallet: c.ClaimerWallet,
			AuthNonce:     c.AuthNonce,
			ClaimTxHash:   c.ClaimTxHash,
			ClaimedAmount: c.ClaimedAmount,
			BlockNumber:   c.BlockNumber,
			Status:        c.Status,
			CreatedAt:     c.CreatedAt.Unix(),
			UpdatedAt:     c.UpdatedAt.Unix(),
		})
	}
	return out
}
