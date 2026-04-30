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
	processedTxs    map[string]bool
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
		processedTxs:    make(map[string]bool),
	}
}

func (t *TronIndexer) Start(ctx context.Context) {
	log.ZInfo(ctx, "starting RedPacket TRON event indexer")

	go func() {
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

	for blockNum := t.lastBlockNum + 1; blockNum <= currentBlock; blockNum++ {
		if err := t.scanBlock(ctx, blockNum); err != nil {
			log.ZWarn(ctx, "redpacket tron scan block failed", err, "block", blockNum)
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
		if txID == "" || t.processedTxs[txID] {
			continue
		}

		if err := t.processTransaction(ctx, txID); err != nil {
			log.ZWarn(ctx, "redpacket tron process tx failed", err, "txID", txID)
		} else {
			t.processedTxs[txID] = true
		}
	}

	return nil
}

func (t *TronIndexer) processTransaction(ctx context.Context, txID string) error {
	var txInfo map[string]interface{}
	err := postJSON(ctx, t.client.fullNodeURL+"/wallet/gettransactioninfobyid", map[string]interface{}{
		"value": txID,
	}, &txInfo)
	if err != nil {
		return err
	}

	contractAddress := t.client.contractBase58
	if logs, ok := txInfo["log"].([]interface{}); ok && len(logs) > 0 {
		for _, logEntry := range logs {
			if logMap, ok := logEntry.(map[string]interface{}); ok {
				if address, ok := logMap["address"].(string); ok && address == contractAddress {
					eventType := t.parseTronEvent(logMap)
					log.ZDebug(ctx, "redpacket tron event detected", "event", eventType, "txID", txID)

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
	if topics, ok := logEntry["topics"].([]interface{}); ok && len(topics) > 0 {
		if topic0, ok := topics[0].(string); ok {
			switch topic0 {
			case "0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0":
				return "Transfer"
			default:
				return "UnknownEvent"
			}
		}
	}
	return "UnknownEvent"
}

func (t *TronIndexer) handleTronPacketCreated(ctx context.Context, logData map[string]interface{}, txID string) {
	log.ZInfo(ctx, "tron PacketCreated event", "txID", txID)
}

func (t *TronIndexer) handleTronPacketClaimed(ctx context.Context, logData map[string]interface{}, txID string) {
	log.ZInfo(ctx, "tron PacketClaimed event", "txID", txID)

	claimer := "unknown"
	amount := "0"

	if topics, ok := logData["topics"].([]interface{}); ok && len(topics) > 1 {
		if claimerTopic, ok := topics[1].(string); ok {
			claimer = claimerTopic
		}
	}

	claim := &model.RedPacketClaim{
		PacketID:      "tron-packet-" + txID[:8],
		ClaimerWallet: claimer,
		ClaimTxHash:   txID,
		ClaimedAmount: amount,
		Status:        "CONFIRMED",
	}

	if err := t.db.SaveClaim(ctx, claim); err != nil {
		log.ZWarn(ctx, "redpacket tron save claim failed", err)
	}
}

func (t *TronIndexer) handleTronPacketRefunded(ctx context.Context, logData map[string]interface{}, txID string) {
	log.ZInfo(ctx, "tron PacketRefunded event", "txID", txID)
}

func (t *TronIndexer) GetLastProcessedBlock() int64 {
	return t.lastBlockNum
}
