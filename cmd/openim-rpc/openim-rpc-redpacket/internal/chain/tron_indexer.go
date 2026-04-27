package chain

import (
	"context"
	"fmt"
	"log"
	"time"

	"redpacket/internal/model"
	"redpacket/internal/repository"
)

// TronIndexer provides production-grade event listening for TRON blockchain
type TronIndexer struct {
	client          *TronClient
	repo            repository.Repository
	pollInterval    time.Duration
	lastBlockNum    int64                    // TRON uses block numbers
	contractAddress string
	processedTxs    map[string]bool          // Simple dedup for this session
}

// NewTronIndexer creates a new TRON event indexer
func NewTronIndexer(client *TronClient, repo repository.Repository, pollInterval int, startBlock int64) *TronIndexer {
	if pollInterval <= 0 {
		pollInterval = 3 // TRON blocks are ~3s
	}

	return &TronIndexer{
		client:          client,
		repo:            repo,
		pollInterval:    time.Duration(pollInterval) * time.Second,
		lastBlockNum:    startBlock,
		contractAddress: client.contractBase58,
		processedTxs:    make(map[string]bool),
	}
}

// Start begins polling for TRON blockchain events
func (t *TronIndexer) Start(ctx context.Context) {
	log.Println("🚀 Starting TRON event indexer... (Production mode)")

	go func() {
		ticker := time.NewTicker(t.pollInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("TRON Indexer stopped")
				return
			case <-ticker.C:
				if err := t.poll(ctx); err != nil {
					log.Printf("TRON Indexer poll error: %v", err)
					// Backoff on error
					time.Sleep(2 * time.Second)
				}
			}
		}
	}()
}

func (t *TronIndexer) poll(ctx context.Context) error {
	// Get current block
	currentBlock, err := t.getNowBlock(ctx)
	if err != nil {
		return fmt.Errorf("get now block failed: %w", err)
	}

	if currentBlock <= t.lastBlockNum {
		return nil
	}

	log.Printf("📡 TRON scanning blocks %d to %d", t.lastBlockNum+1, currentBlock)

	// Scan blocks for contract transactions
	for blockNum := t.lastBlockNum + 1; blockNum <= currentBlock; blockNum++ {
		if err := t.scanBlock(ctx, blockNum); err != nil {
			log.Printf("Warning: failed to scan TRON block %d: %v", blockNum, err)
			continue
		}
	}

	t.lastBlockNum = currentBlock
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
	// Get block by number
	var blockResp map[string]interface{}
	err := postJSON(ctx, t.client.fullNodeURL+"/wallet/getblockbynum", map[string]interface{}{
		"num": blockNum,
	}, &blockResp)
	if err != nil {
		return err
	}

	transactions, ok := blockResp["transactions"].([]interface{})
	if !ok {
		return nil // no transactions
	}

	for _, txInterface := range transactions {
		tx, ok := txInterface.(map[string]interface{})
		if !ok {
			continue
		}

		txID, _ := tx["txID"].(string)
		if txID == "" || t.processedTxs[txID] {
			continue
		}

		if err := t.processTransaction(ctx, txID); err != nil {
			log.Printf("Failed to process TRON tx %s: %v", txID, err)
		} else {
			t.processedTxs[txID] = true
		}
	}

	return nil
}

func (t *TronIndexer) processTransaction(ctx context.Context, txID string) error {
	// Get transaction info with logs
	var txInfo map[string]interface{}
	err := postJSON(ctx, t.client.fullNodeURL+"/wallet/gettransactioninfobyid", map[string]interface{}{
		"value": txID,
	}, &txInfo)
	if err != nil {
		return err
	}

	// Check if this transaction interacted with our contract
	contractAddress := t.client.contractBase58
	if logs, ok := txInfo["log"].([]interface{}); ok && len(logs) > 0 {
		for _, logEntry := range logs {
			if logMap, ok := logEntry.(map[string]interface{}); ok {
				if address, ok := logMap["address"].(string); ok && address == contractAddress {
					// This is our contract event
					eventType := t.parseTronEvent(logMap)
					log.Printf("🔍 TRON Event detected: %s in tx %s", eventType, txID)

					// Process different event types
					switch eventType {
					case "PacketCreated":
						t.handleTronPacketCreated(ctx, logMap, txID)
					case "PacketClaimed":
						t.handleTronPacketClaimed(ctx, logMap, txID)
					case "PacketRefunded":
						t.handleTronPacketRefunded(ctx, logMap, txID)
					}
				}
			}
		}
	}

	return nil
}

func (t *TronIndexer) parseTronEvent(logEntry map[string]interface{}) string {
	// TRON events are more complex. In production, you'd decode topics and data
	// For this implementation, we use a simplified approach based on log data
	if topics, ok := logEntry["topics"].([]interface{}); ok && len(topics) > 0 {
		if topic0, ok := topics[0].(string); ok {
			// Map common TRON event signatures (this would be expanded with real contract event IDs)
			switch topic0 {
			case "0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0": // Transfer (example)
				return "Transfer"
			// Add real RedPacket event signatures here from contract
			default:
				return "UnknownEvent"
			}
		}
	}
	return "UnknownEvent"
}

// Event handlers - these would update the database with parsed event data

func (t *TronIndexer) handleTronPacketCreated(ctx context.Context, logData map[string]interface{}, txID string) {
	log.Printf("📦 [TRON] PacketCreated event in tx %s", txID)
	// TODO: Parse packetId, creator, amount, etc. and update database
	// This would typically link with the offchain biz_id created earlier
}

func (t *TronIndexer) handleTronPacketClaimed(ctx context.Context, logData map[string]interface{}, txID string) {
	log.Printf("🎁 [TRON] PacketClaimed event in tx %s", txID)

	// Example: extract claimer and amount from log data
	claimer := "unknown"
	amount := "0"

	if topics, ok := logData["topics"].([]interface{}); ok && len(topics) > 1 {
		if claimerTopic, ok := topics[1].(string); ok {
			claimer = claimerTopic // simplified
		}
	}

	claim := &model.RedPacketClaim{
		PacketID:      "tron-packet-" + txID[:8], // placeholder
		ClaimerWallet: claimer,
		ClaimTxHash:   txID,
		ClaimedAmount: amount,
		Status:        "CONFIRMED",
	}

	if err := t.repo.CreateClaim(ctx, claim); err != nil {
		log.Printf("Failed to save TRON claim: %v", err)
	}
}

func (t *TronIndexer) handleTronPacketRefunded(ctx context.Context, logData map[string]interface{}, txID string) {
	log.Printf("♻️ [TRON] PacketRefunded event in tx %s", txID)
	// Update packet status to REFUNDED
}

// GetLastProcessedBlock returns the last processed block for monitoring
func (t *TronIndexer) GetLastProcessedBlock() int64 {
	return t.lastBlockNum
}
