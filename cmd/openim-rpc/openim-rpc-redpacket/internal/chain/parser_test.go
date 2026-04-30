package chain

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestParseEventsFromLogs_ParsesRefundEvent(t *testing.T) {
	abiJSON, err := ExtractABIFromEmbeddedArtifact()
	if err != nil {
		t.Fatalf("ExtractABIFromEmbeddedArtifact() error = %v", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(abiJSON)))
	if err != nil {
		t.Fatalf("abi.JSON() error = %v", err)
	}

	eventDef := parsedABI.Events["PacketRefunded"]
	packetID := big.NewInt(101)
	operator := common.HexToAddress("0x1111111111111111111111111111111111111111")
	refundTo := common.HexToAddress("0x2222222222222222222222222222222222222222")
	amount := big.NewInt(8888)

	data, err := eventDef.Inputs.NonIndexed().Pack(amount)
	if err != nil {
		t.Fatalf("Pack() error = %v", err)
	}

	log := &types.Log{
		Address: common.HexToAddress("0x3333333333333333333333333333333333333333"),
		Topics: []common.Hash{
			eventDef.ID,
			common.BigToHash(packetID),
			common.BytesToHash(common.LeftPadBytes(operator.Bytes(), 32)),
			common.BytesToHash(common.LeftPadBytes(refundTo.Bytes(), 32)),
		},
		Data:        data,
		BlockNumber: 77,
		TxHash:      common.HexToHash("0xabc"),
	}

	events, err := ParseEventsFromLogs([]*types.Log{log}, parsedABI)
	if err != nil {
		t.Fatalf("ParseEventsFromLogs() error = %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	event := events[0]
	if event.Name != "PacketRefunded" {
		t.Fatalf("unexpected event name: %s", event.Name)
	}
	if got := GetPacketIDFromEvent(event).String(); got != "101" {
		t.Fatalf("packet id mismatch: got %s", got)
	}
	if got := GetAddressFromEvent(event, "operator").Hex(); got != operator.Hex() {
		t.Fatalf("operator mismatch: got %s want %s", got, operator.Hex())
	}
	if got := GetAddressFromEvent(event, "refundTo").Hex(); got != refundTo.Hex() {
		t.Fatalf("refundTo mismatch: got %s want %s", got, refundTo.Hex())
	}
	if got := GetAmountFromEvent(event).String(); got != "8888" {
		t.Fatalf("amount mismatch: got %s", got)
	}
	if event.BlockNumber != 77 {
		t.Fatalf("block number mismatch: got %d", event.BlockNumber)
	}
	if event.TxHash != common.HexToHash("0xabc") {
		t.Fatalf("tx hash mismatch: got %s", event.TxHash.Hex())
	}
}
