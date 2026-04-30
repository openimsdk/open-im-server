package chain

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"redpacket/internal/model"
	"redpacket/internal/repository"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Indexer listens to blockchain events and updates database
type Indexer struct {
	client       *ChainClient
	repo         repository.Repository
	pollInterval time.Duration
	lastBlock    uint64
	contractAddr common.Address
}

// NewIndexer creates a new event indexer
func NewIndexer(client *ChainClient, repo repository.Repository, pollInterval int, startBlock uint64) *Indexer {
	if pollInterval <= 0 {
		pollInterval = 5
	}

	return &Indexer{
		client:       client,
		repo:         repo,
		pollInterval: time.Duration(pollInterval) * time.Second,
		lastBlock:    startBlock,
		contractAddr: client.contractAddr,
	}
}

// Start begins polling for new events
func (i *Indexer) Start(ctx context.Context) {
	log.Println("🚀 Starting RedPacket event indexer...")

	go func() {
		ticker := time.NewTicker(i.pollInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("Indexer stopped")
				return
			case <-ticker.C:
				if err := i.poll(ctx); err != nil {
					log.Printf("Indexer poll error: %v", err)
				}
			}
		}
	}()
}

func (i *Indexer) poll(ctx context.Context) error {
	// Get latest block
	header, err := i.client.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return fmt.Errorf("get header failed: %w", err)
	}

	currentBlock := header.Number.Uint64()
	if currentBlock <= i.lastBlock {
		return nil
	}

	// Query logs from lastBlock+1 to currentBlock
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(i.lastBlock + 1)),
		ToBlock:   big.NewInt(int64(currentBlock)),
		Addresses: []common.Address{i.contractAddr},
	}

	logs, err := i.client.client.FilterLogs(ctx, query)
	if err != nil {
		return fmt.Errorf("filter logs failed: %w", err)
	}

	// Convert to pointer slice for parser
	logPtrs := make([]*types.Log, len(logs))
	for i, log := range logs {
		logPtrs[i] = &log
	}

	// Parse and process events
	events, err := ParseEventsFromLogs(logPtrs, i.client.contractABI)
	if err != nil {
		return err
	}

	for _, event := range events {
		if err := i.processEvent(ctx, event, logPtrs); err != nil {
			log.Printf("Process event %s failed: %v", event.Name, err)
		}
	}

	i.lastBlock = currentBlock
	log.Printf("✅ Indexed up to block %d, processed %d events", currentBlock, len(events))
	return nil
}

func (i *Indexer) processEvent(ctx context.Context, event *ParsedEvent, logs []*types.Log) error {
	switch event.Name {
	case "PacketCreated":
		return i.handlePacketCreated(ctx, event)
	case "PacketClaimed":
		return i.handlePacketClaimed(ctx, event)
	case "PacketRefunded":
		return i.handlePacketRefunded(ctx, event)
	default:
		log.Printf("Unknown event: %s", event.Name)
		return nil
	}
}

func (i *Indexer) handlePacketCreated(ctx context.Context, event *ParsedEvent) error {
	packetID := GetPacketIDFromEvent(event)
	creator := GetAddressFromEvent(event, "creator")

	log.Printf("📦 PacketCreated: packetId=%s, creator=%s", packetID.String(), creator.Hex())

	// Update database - in real implementation, link with biz_id via offchain record
	// This would typically be triggered by the created-callback first
	return nil
}

func (i *Indexer) handlePacketClaimed(ctx context.Context, event *ParsedEvent) error {
	packetID := GetPacketIDFromEvent(event)
	claimer := GetAddressFromEvent(event, "claimer")
	amount := GetAmountFromEvent(event)
	authNonce := GetUintFromEvent(event, "authNonce")

	log.Printf("🎁 PacketClaimed: packetId=%s, claimer=%s, amount=%s",
		packetID.String(), claimer.Hex(), amount.String())

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

	if err := i.repo.SaveClaim(ctx, claim); err != nil {
		return err
	}
	if err := i.repo.MarkClaimAuthUsed(ctx, authNonce.String()); err != nil {
		return err
	}

	return i.repo.UpdateRedPacketClaimProgress(ctx, packetID.String(), amount.String(), "")
}

func (i *Indexer) handlePacketRefunded(ctx context.Context, event *ParsedEvent) error {
	packetID := GetPacketIDFromEvent(event)
	operator := GetAddressFromEvent(event, "operator")
	refundTo := GetAddressFromEvent(event, "refundTo")
	amount := GetAmountFromEvent(event)

	log.Printf("♻️ PacketRefunded: packetId=%s, operator=%s, refundTo=%s, amount=%s",
		packetID.String(), operator.Hex(), refundTo.Hex(), amount.String())

	if err := i.repo.SaveRefund(ctx, &model.RedPacketRefund{
		PacketID:  packetID.String(),
		RefundTo:  refundTo.Hex(),
		TxHash:    event.TxHash.Hex(),
		Amount:    amount.String(),
		CreatedAt: time.Now(),
	}); err != nil {
		return err
	}

	return i.repo.UpdateRedPacketStatus(ctx, packetID.String(), "REFUNDED")
}
