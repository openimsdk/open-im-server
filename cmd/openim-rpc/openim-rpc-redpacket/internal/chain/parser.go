package chain

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// ParsedEvent represents a parsed blockchain event
type ParsedEvent struct {
	Name string
	Data map[string]interface{}
}

// ParseEventsFromLogs parses logs using the contract ABI
func ParseEventsFromLogs(logs []*types.Log, contractABI abi.ABI) ([]*ParsedEvent, error) {
	var events []*ParsedEvent

	for _, log := range logs {
		if len(log.Topics) == 0 {
			continue
		}

		event, err := parseEvent(log, contractABI)
		if err == nil && event != nil {
			events = append(events, event)
		}
	}

	return events, nil
}

func parseEvent(log *types.Log, contractABI abi.ABI) (*ParsedEvent, error) {
	for name, event := range contractABI.Events {
		if event.ID != log.Topics[0] {
			continue
		}

		data := make(map[string]interface{})

		// Parse indexed parameters from topics
		indexedIdx := 1
		for _, arg := range event.Inputs {
			if arg.Indexed {
				if indexedIdx < len(log.Topics) {
					if arg.Type.T == abi.AddressTy {
						data[arg.Name] = common.BytesToAddress(log.Topics[indexedIdx].Bytes())
					} else if arg.Type.T == abi.UintTy || arg.Type.T == abi.IntTy {
						data[arg.Name] = new(big.Int).SetBytes(log.Topics[indexedIdx].Bytes())
					} else {
						data[arg.Name] = log.Topics[indexedIdx].Hex()
					}
					indexedIdx++
				}
			}
		}

		// Parse non-indexed parameters from data
		if len(log.Data) > 0 {
			unpacked, err := event.Inputs.Unpack(log.Data)
			if err == nil {
				nonIndexedIdx := 0
				for _, arg := range event.Inputs {
					if !arg.Indexed {
						if nonIndexedIdx < len(unpacked) {
							data[arg.Name] = unpacked[nonIndexedIdx]
							nonIndexedIdx++
						}
					}
				}
			}
		}

		return &ParsedEvent{
			Name: name,
			Data: data,
		}, nil
	}

	return nil, fmt.Errorf("unknown event: %s", log.Topics[0].Hex())
}

// GetPacketIDFromEvent extracts packetId from event data
func GetPacketIDFromEvent(event *ParsedEvent) *big.Int {
	if id, ok := event.Data["packetId"]; ok {
		if b, ok := id.(*big.Int); ok {
			return b
		}
	}
	return big.NewInt(0)
}

// GetClaimerFromEvent extracts claimer address from event
func GetClaimerFromEvent(event *ParsedEvent) common.Address {
	if claimer, ok := event.Data["claimer"]; ok {
		if addr, ok := claimer.(common.Address); ok {
			return addr
		}
	}
	return common.Address{}
}

// GetAmountFromEvent extracts amount from event
func GetAmountFromEvent(event *ParsedEvent) *big.Int {
	if amount, ok := event.Data["amount"]; ok {
		if b, ok := amount.(*big.Int); ok {
			return b
		}
	}
	return big.NewInt(0)
}
