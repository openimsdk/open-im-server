package chain

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/log"
)

const defaultIndexerMaxBlocksPerPoll uint64 = 2000

type Indexer struct {
	client             *ChainClient
	db                 controller.RedPacketDatabase
	pollInterval       time.Duration
	lastBlock          uint64
	contractAddr       common.Address
	maxBlocksPerPoll   uint64 // 0 => defaultIndexerMaxBlocksPerPoll
}

func NewIndexer(client *ChainClient, db controller.RedPacketDatabase, pollInterval int, startBlock uint64, maxBlocksPerPoll int) *Indexer {
	if pollInterval <= 0 {
		pollInterval = 5
	}
	var maxB uint64
	if maxBlocksPerPoll > 0 {
		maxB = uint64(maxBlocksPerPoll)
	}
	return &Indexer{
		client:           client,
		db:               db,
		pollInterval:     time.Duration(pollInterval) * time.Second,
		lastBlock:        startBlock,
		contractAddr:     client.contractAddr,
		maxBlocksPerPoll: maxB,
	}
}

func (i *Indexer) chunkEndBlock(chainTip uint64) uint64 {
	maxSpan := i.maxBlocksPerPoll
	if maxSpan == 0 {
		maxSpan = defaultIndexerMaxBlocksPerPoll
	}
	if chainTip <= i.lastBlock {
		return i.lastBlock
	}
	span := chainTip - i.lastBlock
	if span > maxSpan {
		return i.lastBlock + maxSpan
	}
	return chainTip
}

func (i *Indexer) Start(ctx context.Context) {
	log.ZInfo(ctx, "starting RedPacket ETH event indexer")

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.ZError(ctx, "redpacket eth indexer panic recovered", fmt.Errorf("%v", r))
			}
		}()
		ticker := time.NewTicker(i.pollInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.ZInfo(ctx, "redpacket eth indexer stopped")
				return
			case <-ticker.C:
				if err := i.poll(ctx); err != nil {
					log.ZWarn(ctx, "redpacket eth indexer poll error", err)
				}
			}
		}
	}()

	// Compensation loop: periodically scan DB for expired-but-unclosed packets
	// and mark them EXPIRED so the UI reflects the correct state even if the
	// on-chain refund event was missed.
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.ZError(ctx, "redpacket eth compensation panic recovered", fmt.Errorf("%v", r))
			}
		}()
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := i.compensate(ctx); err != nil {
					log.ZWarn(ctx, "redpacket eth compensation error", err)
				}
			}
		}
	}()
}

func (i *Indexer) compensate(ctx context.Context) error {
	now := time.Now().Unix()
	packets, err := i.db.GetExpiredPendingPackets(ctx, now)
	if err != nil {
		return fmt.Errorf("get expired packets failed: %w", err)
	}
	for _, rp := range packets {
		if err := i.db.UpdateRedPacketStatus(ctx, rp.PacketID, "EXPIRED"); err != nil {
			log.ZWarn(ctx, "redpacket eth compensation mark expired failed", err, "packetID", rp.PacketID)
			continue
		}
		log.ZInfo(ctx, "redpacket eth compensation: marked packet EXPIRED", "packetID", rp.PacketID)
	}
	return nil
}

func (i *Indexer) poll(ctx context.Context) error {
	header, err := i.client.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return fmt.Errorf("get header failed: %w", err)
	}

	chainTip := header.Number.Uint64()
	if chainTip <= i.lastBlock {
		return nil
	}

	toBlock := i.chunkEndBlock(chainTip)
	if toBlock <= i.lastBlock {
		return nil
	}

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(i.lastBlock + 1)),
		ToBlock:   big.NewInt(int64(toBlock)),
		Addresses: []common.Address{i.contractAddr},
	}

	logs, err := i.client.client.FilterLogs(ctx, query)
	if err != nil {
		return fmt.Errorf("filter logs failed (blocks %d-%d): %w", i.lastBlock+1, toBlock, err)
	}

	logPtrs := make([]*types.Log, len(logs))
	for idx := range logs {
		logPtrs[idx] = &logs[idx]
	}

	events, err := ParseEventsFromLogs(logPtrs, i.client.contractABI)
	if err != nil {
		return err
	}

	for _, event := range events {
		if err := i.processEvent(ctx, event); err != nil {
			log.ZWarn(ctx, "process redpacket eth event failed", err, "event", event.Name)
		}
	}

	i.lastBlock = toBlock
	if toBlock < chainTip {
		log.ZDebug(ctx, "redpacket eth indexer chunk done, catching up", "indexedTo", toBlock, "chainTip", chainTip, "events", len(events))
	} else {
		log.ZInfo(ctx, "redpacket eth indexed", "block", toBlock, "events", len(events))
	}
	return nil
}

func (i *Indexer) processEvent(ctx context.Context, event *ParsedEvent) error {
	switch event.Name {
	case "PacketCreated":
		return i.handlePacketCreated(ctx, event)
	case "PacketClaimed":
		return i.handlePacketClaimed(ctx, event)
	case "PacketRefunded":
		return i.handlePacketRefunded(ctx, event)
	default:
		return nil
	}
}

func (i *Indexer) handlePacketCreated(ctx context.Context, event *ParsedEvent) error {
	packetID := GetPacketIDFromEvent(event)
	creator := GetAddressFromEvent(event, "creator")
	log.ZInfo(ctx, "PacketCreated event", "packetID", packetID.String(), "creator", creator.Hex())
	return nil
}

func (i *Indexer) handlePacketClaimed(ctx context.Context, event *ParsedEvent) error {
	packetID := GetPacketIDFromEvent(event)
	claimer := GetAddressFromEvent(event, "claimer")
	amount := GetAmountFromEvent(event)
	authNonce := GetUintFromEvent(event, "authNonce")

	log.ZInfo(ctx, "PacketClaimed event", "packetID", packetID.String(), "claimer", claimer.Hex(), "amount", amount.String())

	claim := &model.RedPacketClaim{
		PacketID:      packetID.String(),
		ClaimerWallet: claimer.Hex(),
		AuthNonce:     authNonce.String(),
		ClaimTxHash:   event.TxHash.Hex(),
		ClaimedAmount: amount.String(),
		BlockNumber:   event.BlockNumber,
		Status:        "CONFIRMED",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := i.db.SaveClaim(ctx, claim); err != nil {
		return err
	}
	if err := i.db.MarkClaimAuthUsed(ctx, authNonce.String()); err != nil {
		return err
	}
	// Pass "" for forced status; DB layer auto-derives COMPLETED/ACTIVE.
	// TxHash is the idempotency key: prevents double-counting if ClaimResult RPC
	// already processed this same transaction.
	return i.db.UpdateRedPacketClaimProgress(ctx, packetID.String(), amount.String(), "", event.TxHash.Hex())
}

func (i *Indexer) handlePacketRefunded(ctx context.Context, event *ParsedEvent) error {
	packetID := GetPacketIDFromEvent(event)
	refundTo := GetAddressFromEvent(event, "refundTo")
	amount := GetAmountFromEvent(event)

	log.ZInfo(ctx, "PacketRefunded event", "packetID", packetID.String(), "refundTo", refundTo.Hex(), "amount", amount.String())

	if err := i.db.SaveRefund(ctx, &model.RedPacketRefund{
		PacketID:  packetID.String(),
		RefundTo:  refundTo.Hex(),
		TxHash:    event.TxHash.Hex(),
		Amount:    amount.String(),
		CreatedAt: time.Now(),
	}); err != nil {
		return err
	}

	return i.db.UpdateRedPacketStatus(ctx, packetID.String(), "REFUNDED")
}
