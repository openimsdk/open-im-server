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
	opUserID := mcontext.GetOpUserID(ctx)
	if opUserID == "" {
		return nil, servererrs.ErrNoPermission.WrapMsg("op user id is empty")
	}
	if strings.TrimSpace(req.BizID) == "" || strings.TrimSpace(req.TxHash) == "" {
		return nil, errs.ErrArgs.WrapMsg("biz_id and tx_hash are required")
	}

	rp, err := s.db.GetRedPacketByBizID(ctx, req.BizID)
	if err != nil {
		return nil, err
	}
	if rp.CreatorUserID != opUserID {
		return nil, servererrs.ErrNoPermission.WrapMsg("only the creator can submit the creation callback")
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
		CreatorWallet:   createdPacket.CreatorWallet,
		PacketType:      createdPacket.PacketType,
		Token:           createdPacket.Token,
		TotalAmount:     createdPacket.TotalAmount,
		TotalShares:     createdPacket.TotalShares,
		ExpiryAt:        createdPacket.ExpiryAt,
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
		return nil, errs.ErrInternalServer.WrapMsg("signer key not configured; cannot issue claim signature")
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

	txSuccess, events, err := s.parseChainReceiptWithStatus(ctx, rp, req.TxHash)
	if err != nil {
		log.ZWarn(ctx, "parse claim receipt failed", err, "txHash", req.TxHash)
		return &pbredpacket.ClaimResultResp{}, nil
	}
	if !txSuccess {
		if markErr := s.markClaimFailed(ctx, req.PacketID, currentUserID, req.Claimer, req.TxHash); markErr != nil {
			log.ZWarn(ctx, "mark claim failed status failed", markErr, "txHash", req.TxHash)
		}
		return &pbredpacket.ClaimResultResp{}, nil
	}

	claimedEvent, err := resolveClaimedEventFromParsedEvents(rp, events)
	if err != nil {
		log.ZWarn(ctx, "resolve claim event failed", err, "txHash", req.TxHash)
		if markErr := s.markClaimFailed(ctx, req.PacketID, currentUserID, req.Claimer, req.TxHash); markErr != nil {
			log.ZWarn(ctx, "mark claim failed status failed", markErr, "txHash", req.TxHash)
		}
		return &pbredpacket.ClaimResultResp{}, nil
	}
	if claimedEvent == nil {
		if markErr := s.markClaimFailed(ctx, req.PacketID, currentUserID, req.Claimer, req.TxHash); markErr != nil {
			log.ZWarn(ctx, "mark claim failed status failed", markErr, "txHash", req.TxHash)
		}
		return &pbredpacket.ClaimResultResp{}, nil
	}
	if !strings.EqualFold(claimedEvent.ClaimerWallet, req.Claimer) {
		if markErr := s.markClaimFailed(ctx, req.PacketID, currentUserID, req.Claimer, req.TxHash); markErr != nil {
			log.ZWarn(ctx, "mark claim failed status failed", markErr, "txHash", req.TxHash)
		}
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

	// Pass "" for status so the DB layer auto-derives COMPLETED/ACTIVE.
	// Pass req.TxHash as the idempotency key so concurrent indexer processing
	// of the same transaction cannot double-count the claim.
	if err := s.db.UpdateRedPacketClaimProgress(ctx, req.PacketID, claimedEvent.Amount, "", req.TxHash); err != nil {
		return nil, err
	}
	return &pbredpacket.ClaimResultResp{}, nil
}

func (s *redPacketServer) parseChainReceiptWithStatus(ctx context.Context, rp *model.RedPacket, txHash string) (bool, []*chain.ParsedEvent, error) {
	switch rp.ChainType {
	case "EVM":
		if s.chainClient == nil {
			return false, nil, errs.ErrInternalServer.WrapMsg("evm client is unavailable")
		}
		return s.chainClient.ParseTransactionReceiptWithStatus(ctx, common.HexToHash(txHash))
	case "TRON":
		if s.tronClient == nil {
			return false, nil, errs.ErrInternalServer.WrapMsg("tron client is unavailable")
		}
		return s.tronClient.ParseTransactionReceiptWithStatus(ctx, txHash)
	default:
		return false, nil, errs.ErrArgs.WrapMsg("unsupported chain_type: " + rp.ChainType)
	}
}

func (s *redPacketServer) markClaimFailed(ctx context.Context, packetID, userID, claimer, txHash string) error {
	return s.db.SaveClaim(ctx, &model.RedPacketClaim{
		PacketID:      packetID,
		UserID:        userID,
		ClaimerWallet: claimer,
		ClaimTxHash:   txHash,
		Status:        "FAILED",
		UpdatedAt:     time.Now(),
	})
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

type refundedEventSnapshot struct {
	RefundTo    string
	Amount      string
	BlockNumber uint64
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
		// Offline mode: no chain client configured; caller must supply packet_id directly.
		if s.chainClient == nil {
			if fallbackPacketID == "" {
				return nil, errs.ErrArgs.WrapMsg("packet_id is required when EVM client is unavailable")
			}
			return buildFallbackCreatedPacket(rp, fallbackPacketID), nil
		}

		success, events, err := s.chainClient.ParseTransactionReceiptWithStatus(ctx, common.HexToHash(txHashHex))
		if err != nil {
			return nil, errs.ErrInternalServer.WrapMsg("parse created tx failed: " + err.Error())
		}
		if !success {
			return nil, errs.ErrArgs.WrapMsg("created tx execution failed on chain")
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
		return nil, errs.ErrInternalServer.WrapMsg("PacketCreated event not found in tx: " + txHashHex)
	case "TRON":
		// Offline mode: no chain client configured; caller must supply packet_id directly.
		if s.tronClient == nil {
			if fallbackPacketID == "" {
				return nil, errs.ErrArgs.WrapMsg("packet_id is required when TRON client is unavailable")
			}
			return buildFallbackCreatedPacket(rp, fallbackPacketID), nil
		}

		success, events, err := s.tronClient.ParseTransactionReceiptWithStatus(ctx, txHashHex)
		if err != nil {
			return nil, errs.ErrInternalServer.WrapMsg("parse tron created tx failed: " + err.Error())
		}
		if !success {
			return nil, errs.ErrArgs.WrapMsg("created tx execution failed on chain")
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
		return nil, errs.ErrInternalServer.WrapMsg("PacketCreated event not found in TRON tx: " + txHashHex)
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

// validateCreateBaseFields validates the fields shared by every red packet type.
// It does not look up creator identity or scope; those are handled by the per-type hooks.
func validateCreateBaseFields(req *pbredpacket.CreateOrderReq) (*big.Int, error) {
	if strings.TrimSpace(req.CreatorWallet) == "" {
		return nil, errs.ErrArgs.WrapMsg("creator_wallet is required")
	}
	if strings.TrimSpace(req.TotalAmount) == "" {
		return nil, errs.ErrArgs.WrapMsg("total_amount is required")
	}
	total, ok := new(big.Int).SetString(req.TotalAmount, 10)
	if !ok || total.Sign() <= 0 {
		return nil, errs.ErrArgs.WrapMsg("total_amount must be a positive integer string", "totalAmount", req.TotalAmount)
	}
	if req.ExpiryAt != 0 && req.ExpiryAt <= time.Now().Unix() {
		return nil, errs.ErrArgs.WrapMsg("expiry_at must be 0 or a future unix timestamp", "expiryAt", req.ExpiryAt)
	}
	return total, nil
}

// validateCreatorScope verifies group membership / friend relationship for the creator
// based on the requested scope. PUBLIC scope skips relationship checks.
func (s *redPacketServer) validateCreatorScope(ctx context.Context, req *pbredpacket.CreateOrderReq) error {
	creatorUserID := mcontext.GetOpUserID(ctx)
	if creatorUserID == "" {
		return servererrs.ErrNoPermission.WrapMsg("op user id is empty")
	}
	switch normalizeScopeType(req.ScopeType) {
	case "GROUP":
		return s.ensureGroupEligibility(ctx, req.GroupID, creatorUserID)
	case "DIRECT":
		if strings.TrimSpace(req.ReceiverUserID) != "" {
			if err := s.ensureFriendRelationship(ctx, creatorUserID, req.ReceiverUserID); err != nil {
				return err
			}
		}
		for _, receiverID := range req.ReceiverUserIDs {
			if strings.TrimSpace(receiverID) == "" {
				continue
			}
			if err := s.ensureFriendRelationship(ctx, creatorUserID, receiverID); err != nil {
				return err
			}
		}
		return nil
	default:
		return nil
	}
}

// validateFixedPacketCreate validates fixed red packets:
//   - shared base fields
//   - scope_type must be GROUP (fixed packets are group-only; claim validators require group_id)
//   - 0 < total_shares <= maxTotalShares
//   - total_amount must be divisible by total_shares (each share is an integer in min units)
//   - creator must be an active member of the group
func (s *redPacketServer) validateFixedPacketCreate(ctx context.Context, req *pbredpacket.CreateOrderReq) error {
	total, err := validateCreateBaseFields(req)
	if err != nil {
		return err
	}
	if normalizeScopeType(req.ScopeType) != "GROUP" {
		return errs.ErrArgs.WrapMsg("fixed packet must use scope_type=GROUP")
	}
	if req.TotalShares <= 0 {
		return errs.ErrArgs.WrapMsg("total_shares must be positive for fixed packet", "totalShares", req.TotalShares)
	}
	if req.TotalShares > maxTotalShares {
		return errs.ErrArgs.WrapMsg(fmt.Sprintf("total_shares must not exceed %d for fixed packet", maxTotalShares), "totalShares", req.TotalShares)
	}
	shares := big.NewInt(int64(req.TotalShares))
	if new(big.Int).Mod(total, shares).Sign() != 0 {
		return errs.ErrArgs.WrapMsg("total_amount must be divisible by total_shares for fixed packet",
			"totalAmount", req.TotalAmount, "totalShares", req.TotalShares)
	}
	return s.validateCreatorScope(ctx, req)
}

// validateRandomPacketCreate validates random (lucky) red packets:
//   - shared base fields
//   - scope_type must be GROUP (random packets are group-only; claim validators require group_id)
//   - 0 < total_shares <= maxTotalShares
//   - total_amount >= total_shares (at least 1 min unit per share)
//   - creator must be an active member of the group
func (s *redPacketServer) validateRandomPacketCreate(ctx context.Context, req *pbredpacket.CreateOrderReq) error {
	total, err := validateCreateBaseFields(req)
	if err != nil {
		return err
	}
	if normalizeScopeType(req.ScopeType) != "GROUP" {
		return errs.ErrArgs.WrapMsg("random packet must use scope_type=GROUP")
	}
	if req.TotalShares <= 0 {
		return errs.ErrArgs.WrapMsg("total_shares must be positive for random packet", "totalShares", req.TotalShares)
	}
	if req.TotalShares > maxTotalShares {
		return errs.ErrArgs.WrapMsg(fmt.Sprintf("total_shares must not exceed %d for random packet", maxTotalShares), "totalShares", req.TotalShares)
	}
	shares := big.NewInt(int64(req.TotalShares))
	if total.Cmp(shares) < 0 {
		return errs.ErrArgs.WrapMsg("total_amount must be >= total_shares for random packet",
			"totalAmount", req.TotalAmount, "totalShares", req.TotalShares)
	}
	return s.validateCreatorScope(ctx, req)
}

// validateTransferPacketCreate validates transfer red packets:
//   - shared base fields
//   - scope_type must be DIRECT (transfer is a 1-to-1 direct send)
//   - total_shares == 1
//   - exactly one receiver_user_id (receiver_user_ids must be empty)
//   - receiver must not be the creator (no self-transfer)
//   - creator and receiver must be friends
func (s *redPacketServer) validateTransferPacketCreate(ctx context.Context, req *pbredpacket.CreateOrderReq) error {
	if _, err := validateCreateBaseFields(req); err != nil {
		return err
	}
	if normalizeScopeType(req.ScopeType) != "DIRECT" {
		return errs.ErrArgs.WrapMsg("transfer packet must use scope_type=DIRECT")
	}
	if req.TotalShares != 1 {
		return errs.ErrArgs.WrapMsg("transfer packet must have total_shares == 1", "totalShares", req.TotalShares)
	}
	// Reject ambiguous input: receiver_user_ids is not applicable for transfer.
	if len(req.ReceiverUserIDs) > 0 {
		return errs.ErrArgs.WrapMsg("transfer packet uses receiver_user_id (singular), not receiver_user_ids")
	}
	receiverUserID := strings.TrimSpace(req.ReceiverUserID)
	if receiverUserID == "" {
		return errs.ErrArgs.WrapMsg("receiver_user_id is required for transfer packet")
	}
	creatorUserID := mcontext.GetOpUserID(ctx)
	if creatorUserID == "" {
		return servererrs.ErrNoPermission.WrapMsg("op user id is empty")
	}
	if creatorUserID == receiverUserID {
		return errs.ErrArgs.WrapMsg("transfer packet cannot be sent to yourself")
	}
	return s.ensureFriendRelationship(ctx, creatorUserID, receiverUserID)
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
	// Check status first to give precise error messages for each terminal state.
	switch rp.Status {
	case "ACTIVE":
		// ok, continue to expiry check
	case "REFUNDED":
		return errs.ErrArgs.WrapMsg("packet has been refunded")
	case "EXPIRED":
		return errs.ErrArgs.WrapMsg("packet has expired")
	default:
		return errs.ErrArgs.WrapMsg("packet is not claimable, current status: " + rp.Status)
	}
	// Guard against the race where status is still ACTIVE but expiry has passed.
	if rp.ExpiryAt > 0 && rp.ExpiryAt <= time.Now().Unix() {
		return errs.ErrArgs.WrapMsg("packet has expired")
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

// ensureGroupEligibility verifies that userID is an active member of groupID.
func (s *redPacketServer) ensureGroupEligibility(ctx context.Context, groupID, userID string) error {
	groupID = strings.TrimSpace(groupID)
	userID = strings.TrimSpace(userID)
	if groupID == "" {
		return errs.ErrArgs.WrapMsg("group_id is required for group claim")
	}
	if userID == "" {
		return errs.ErrArgs.WrapMsg("user_id is required for group claim")
	}
	if s.groupClient == nil {
		return servererrs.ErrInternalServer.WrapMsg("group client is not initialized")
	}
	if _, err := s.groupClient.GetGroupMemberInfo(ctx, groupID, userID); err != nil {
		if errs.ErrRecordNotFound.Is(err) {
			return errs.ErrNoPermission.WrapMsg("user is not a member of the group", "groupID", groupID, "userID", userID)
		}
		return err
	}
	return nil
}

// ensureFriendRelationship verifies that userA and userB are mutual friends.
// It is used in two contexts:
//   - validateCreatorScope (DIRECT scope): checking that each listed receiver is
//     a friend of the creator. In that path userA == userB is theoretically possible
//     (creator adding themselves to a list), which is allowed here; the transfer
//     validator has its own explicit self-transfer prohibition.
//   - validateTransferPacketClaim: re-confirming the friendship at claim time.
//
// Self-transfer is intentionally allowed at this level; call sites that need to
// prohibit it (e.g. validateTransferPacketCreate) must do so before calling here.
func (s *redPacketServer) ensureFriendRelationship(ctx context.Context, userA, userB string) error {
	userA = strings.TrimSpace(userA)
	userB = strings.TrimSpace(userB)
	if userA == "" || userB == "" {
		return errs.ErrArgs.WrapMsg("both user IDs are required for friend relationship check")
	}
	if userA == userB {
		return nil
	}
	if s.relationClient == nil {
		return servererrs.ErrInternalServer.WrapMsg("relation client is not initialized")
	}
	ok, err := s.relationClient.IsFriend(ctx, userA, userB)
	if err != nil {
		return err
	}
	if !ok {
		return errs.ErrNoPermission.WrapMsg("users are not friends", "userA", userA, "userB", userB)
	}
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
	return resolveClaimedEventFromParsedEvents(rp, events)
}

func resolveClaimedEventFromParsedEvents(rp *model.RedPacket, events []*chain.ParsedEvent) (*claimedEventSnapshot, error) {
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

func resolveRefundedEventFromParsedEvents(rp *model.RedPacket, events []*chain.ParsedEvent) (*refundedEventSnapshot, error) {
	for _, event := range events {
		if event.Name != "PacketRefunded" {
			continue
		}
		packetID := chain.GetPacketIDFromEvent(event).String()
		if packetID != rp.PacketID {
			return nil, errs.ErrArgs.WrapMsg(fmt.Sprintf("refund event packet mismatch: got %s want %s", packetID, rp.PacketID))
		}
		return &refundedEventSnapshot{
			RefundTo:    strings.ToLower(chain.GetAddressFromEvent(event, "refundTo").Hex()),
			Amount:      chain.GetAmountFromEvent(event).String(),
			BlockNumber: event.BlockNumber,
		}, nil
	}
	return nil, nil
}

// maxTotalShares caps the number of shares to prevent abuse.
const maxTotalShares = 10_000

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

// RequestRefund allows the red-packet creator to submit an on-chain refund
// transaction for an expired packet. The indexer will asynchronously pick up
// the on-chain RefundPacket event and mark the packet as REFUNDED in the DB.
func (s *redPacketServer) RequestRefund(ctx context.Context, req *pbredpacket.RequestRefundReq) (*pbredpacket.RequestRefundResp, error) {
	currentUserID := mcontext.GetOpUserID(ctx)
	if currentUserID == "" {
		return nil, servererrs.ErrNoPermission.WrapMsg("op user id is empty")
	}
	if req.GetPacketID() == "" {
		return nil, errs.ErrArgs.WrapMsg("packet_id is required")
	}

	rp, err := s.db.GetRedPacketByPacketID(ctx, req.GetPacketID())
	if err != nil {
		return nil, err
	}
	if rp.CreatorUserID != currentUserID {
		return nil, errs.ErrNoPermission.WrapMsg("only the creator can request a refund")
	}
	if rp.Status == "REFUNDED" {
		return &pbredpacket.RequestRefundResp{TxHash: "", Status: "REFUNDED"}, nil
	}
	if rp.ExpiryAt > 0 && time.Now().Unix() < rp.ExpiryAt {
		return nil, errs.ErrArgs.WrapMsg("red packet has not expired yet")
	}

	// Submit the on-chain refund transaction.
	var txHash string
	if s.chainClient != nil {
		txHash, err = s.chainClient.RefundPacket(ctx, rp.PacketID)
		if err != nil {
			return nil, errs.ErrInternalServer.WrapMsg("submit refund tx failed: " + err.Error())
		}
	} else if s.tronClient != nil {
		packetIDBig, ok := new(big.Int).SetString(rp.PacketID, 10)
		if !ok {
			return nil, errs.ErrInternalServer.WrapMsg("invalid packet id format")
		}
		txHash, err = s.tronClient.SendAdminTransaction(ctx, "refundPacket", packetIDBig)
		if err != nil {
			return nil, errs.ErrInternalServer.WrapMsg("submit tron refund tx failed: " + err.Error())
		}
	} else {
		return nil, errs.ErrInternalServer.WrapMsg("no blockchain client configured")
	}

	log.ZInfo(ctx, "redpacket refund submitted", "packetID", rp.PacketID, "txHash", txHash)
	txSuccess, events, parseErr := s.parseChainReceiptWithStatus(ctx, rp, txHash)
	if parseErr != nil {
		log.ZWarn(ctx, "parse refund receipt failed, fallback to async indexer", parseErr, "packetID", rp.PacketID, "txHash", txHash)
		return &pbredpacket.RequestRefundResp{TxHash: txHash, Status: "PENDING"}, nil
	}
	if !txSuccess {
		return &pbredpacket.RequestRefundResp{TxHash: txHash, Status: "FAILED"}, nil
	}

	refundedEvent, err := resolveRefundedEventFromParsedEvents(rp, events)
	if err != nil {
		log.ZWarn(ctx, "resolve refunded event failed, fallback to async indexer", err, "packetID", rp.PacketID, "txHash", txHash)
		return &pbredpacket.RequestRefundResp{TxHash: txHash, Status: "PENDING"}, nil
	}
	if refundedEvent == nil {
		return &pbredpacket.RequestRefundResp{TxHash: txHash, Status: "PENDING"}, nil
	}

	if err := s.db.SaveRefund(ctx, &model.RedPacketRefund{
		PacketID:  rp.PacketID,
		RefundTo:  refundedEvent.RefundTo,
		TxHash:    txHash,
		Amount:    refundedEvent.Amount,
		CreatedAt: time.Now(),
	}); err != nil {
		return nil, err
	}
	if err := s.db.UpdateRedPacketStatus(ctx, rp.PacketID, "REFUNDED"); err != nil {
		return nil, err
	}
	return &pbredpacket.RequestRefundResp{TxHash: txHash, Status: "REFUNDED"}, nil
}

func (s *redPacketServer) GetRefund(ctx context.Context, req *pbredpacket.GetRefundReq) (*pbredpacket.GetRefundResp, error) {
	if req.GetPacketID() == "" {
		return nil, errs.ErrArgs.WrapMsg("packet_id is required")
	}
	refund, err := s.db.GetRefundByPacketID(ctx, req.GetPacketID())
	if err != nil {
		return nil, err
	}
	return &pbredpacket.GetRefundResp{
		PacketID:  refund.PacketID,
		RefundTo:  refund.RefundTo,
		TxHash:    refund.TxHash,
		Amount:    refund.Amount,
		CreatedAt: refund.CreatedAt.Unix(),
	}, nil
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
