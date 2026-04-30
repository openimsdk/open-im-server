package chain

import (
	"context"
	"fmt"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/log"
)

type TronIndexer struct {
	client          *TronClient
	db              controller.RedPacketDatabase
	pollInterval    time.Duration
	lastBlockNum    int64
	contractAddress string
}

func NewTronIndexer(client *TronClient, db controller.RedPacketDatabase, pollInterval int, startBlock int64) *TronIndexer {
	if pollInterval <= 0 {
		pollInterval = 3
	}
	return &TronIndexer{
		client:          client,
		db:              db,
		pollInterval:    time.Duration(pollInterval) * time.Second,
		lastBlockNum:    startBlock,
		contractAddress: client.contractBase58,
	}
}

func (t *TronIndexer) Start(ctx context.Context) {
	log.ZInfo(ctx, "starting RedPacket TRON event indexer")

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.ZError(ctx, "redpacket tron indexer panic recovered", fmt.Errorf("%v", r))
			}
		}()
		ticker := time.NewTicker(t.pollInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.ZInfo(ctx, "redpacket tron indexer stopped")
				return
			case <-ticker.C:
				if err := t.poll(ctx); err != nil {
					log.ZWarn(ctx, "redpacket tron indexer poll error", err)
					time.Sleep(2 * time.Second)
				}
			}
		}
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.ZError(ctx, "redpacket tron compensation panic recovered", fmt.Errorf("%v", r))
			}
		}()
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := t.compensate(ctx); err != nil {
					log.ZWarn(ctx, "redpacket tron compensation error", err)
				}
			}
		}
	}()
}

func (t *TronIndexer) compensate(ctx context.Context) error {
	now := time.Now().Unix()
	packets, err := t.db.GetExpiredPendingPackets(ctx, now)
	if err != nil {
		return fmt.Errorf("get expired packets failed: %w", err)
	}
	for _, rp := range packets {
		if err := t.db.UpdateRedPacketStatus(ctx, rp.PacketID, "EXPIRED"); err != nil {
			log.ZWarn(ctx, "redpacket tron compensation mark expired failed", err, "packetID", rp.PacketID)
			continue
		}
		log.ZInfo(ctx, "redpacket tron compensation: marked packet EXPIRED", "packetID", rp.PacketID)
	}
	return nil
}

func (t *TronIndexer) poll(ctx context.Context) error {
	currentBlock, err := t.getNowBlock(ctx)
	if err != nil {
		return fmt.Errorf("get now block failed: %w", err)
	}

	if currentBlock <= t.lastBlockNum {
		return nil
	}

	log.ZDebug(ctx, "redpacket tron scanning blocks", "from", t.lastBlockNum+1, "to", currentBlock)

	// Advance the cursor only up to the last successfully processed block so
	// that a transient RPC failure does not cause blocks to be silently skipped.
	lastOK := t.lastBlockNum
	for blockNum := t.lastBlockNum + 1; blockNum <= currentBlock; blockNum++ {
		if err := t.scanBlock(ctx, blockNum); err != nil {
			log.ZWarn(ctx, "redpacket tron scan block failed", err, "block", blockNum)
			break
		}
		lastOK = blockNum
	}

	t.lastBlockNum = lastOK
	return nil
}

func (t *TronIndexer) getNowBlock(ctx context.Context) (int64, error) {
	var resp map[string]interface{}
	err := postJSON(ctx, t.client.fullNodeURL+"/wallet/getnowblock", map[string]interface{}{}, &resp)
	if err != nil {
		return 0, err
	}

	if blockHeader, ok := resp["block_header"].(map[string]interface{}); ok {
		if rawData, ok := blockHeader["raw_data"].(map[string]interface{}); ok {
			if number, ok := rawData["number"].(float64); ok {
				return int64(number), nil
			}
		}
	}

	return 0, fmt.Errorf("could not parse block number")
}

func (t *TronIndexer) scanBlock(ctx context.Context, blockNum int64) error {
	var blockResp map[string]interface{}
	err := postJSON(ctx, t.client.fullNodeURL+"/wallet/getblockbynum", map[string]interface{}{
		"num": blockNum,
	}, &blockResp)
	if err != nil {
		return err
	}

	transactions, ok := blockResp["transactions"].([]interface{})
	if !ok {
		return nil
	}

	for _, txInterface := range transactions {
		tx, ok := txInterface.(map[string]interface{})
		if !ok {
			continue
		}

		txID, _ := tx["txID"].(string)
		if txID == "" {
			continue
		}

		if err := t.processTransaction(ctx, txID); err != nil {
			log.ZWarn(ctx, "redpacket tron process tx failed", err, "txID", txID)
		}
	}

	return nil
}

// processTransaction parses the on-chain receipt through the ABI (same path as
// the ETH indexer) and dispatches each decoded event to the appropriate handler.
func (t *TronIndexer) processTransaction(ctx context.Context, txID string) error {
	events, err := t.client.ParseTransactionReceipt(ctx, txID)
	if err != nil {
		return fmt.Errorf("parse tron tx receipt failed: %w", err)
	}

	for _, event := range events {
		log.ZDebug(ctx, "redpacket tron event detected", "event", event.Name, "txID", txID)
		switch event.Name {
		case "PacketCreated":
			if err := t.handleTronPacketCreated(ctx, event, txID); err != nil {
				log.ZWarn(ctx, "redpacket tron handlePacketCreated failed", err, "txID", txID)
			}
		case "PacketClaimed":
			if err := t.handleTronPacketClaimed(ctx, event, txID); err != nil {
				log.ZWarn(ctx, "redpacket tron handlePacketClaimed failed", err, "txID", txID)
			}
		case "PacketRefunded":
			if err := t.handleTronPacketRefunded(ctx, event, txID); err != nil {
				log.ZWarn(ctx, "redpacket tron handlePacketRefunded failed", err, "txID", txID)
			}
		}
	}
	return nil
}

func (t *TronIndexer) handleTronPacketCreated(ctx context.Context, event *ParsedEvent, txID string) error {
	packetID := GetPacketIDFromEvent(event)
	creator := GetAddressFromEvent(event, "creator")
	log.ZInfo(ctx, "tron PacketCreated event", "packetID", packetID.String(), "creator", creator.Hex(), "txID", txID)
	return nil
}

func (t *TronIndexer) handleTronPacketClaimed(ctx context.Context, event *ParsedEvent, txID string) error {
	packetID := GetPacketIDFromEvent(event)
	claimer := GetAddressFromEvent(event, "claimer")
	amount := GetAmountFromEvent(event)
	authNonce := GetUintFromEvent(event, "authNonce")

	log.ZInfo(ctx, "tron PacketClaimed event", "packetID", packetID.String(), "claimer", claimer.Hex(), "amount", amount.String(), "txID", txID)

	claim := &model.RedPacketClaim{
		PacketID:      packetID.String(),
		ClaimerWallet: claimer.Hex(),
		AuthNonce:     authNonce.String(),
		ClaimTxHash:   txID,
		ClaimedAmount: amount.String(),
		BlockNumber:   event.BlockNumber,
		Status:        "CONFIRMED",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := t.db.SaveClaim(ctx, claim); err != nil {
		return err
	}
	if err := t.db.MarkClaimAuthUsed(ctx, authNonce.String()); err != nil {
		return err
	}
	// Pass "" for forced status; DB layer auto-derives COMPLETED/ACTIVE.
	// txID is the idempotency key: prevents double-counting if ClaimResult RPC
	// already processed this same transaction.
	return t.db.UpdateRedPacketClaimProgress(ctx, packetID.String(), amount.String(), "", txID)
}

func (t *TronIndexer) handleTronPacketRefunded(ctx context.Context, event *ParsedEvent, txID string) error {
	packetID := GetPacketIDFromEvent(event)
	refundTo := GetAddressFromEvent(event, "refundTo")
	amount := GetAmountFromEvent(event)

	log.ZInfo(ctx, "tron PacketRefunded event", "packetID", packetID.String(), "refundTo", refundTo.Hex(), "amount", amount.String(), "txID", txID)

	if err := t.db.SaveRefund(ctx, &model.RedPacketRefund{
		PacketID:  packetID.String(),
		RefundTo:  refundTo.Hex(),
		TxHash:    txID,
		Amount:    amount.String(),
		CreatedAt: time.Now(),
	}); err != nil {
		return err
	}
	return t.db.UpdateRedPacketStatus(ctx, packetID.String(), "REFUNDED")
}

func (t *TronIndexer) GetLastProcessedBlock() int64 {
	return t.lastBlockNum
}
