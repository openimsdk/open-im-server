package redpacket

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	pbredpacket "github.com/openimsdk/protocol/redpacket"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
)

func (s *redPacketServer) SetSigner(ctx context.Context, req *pbredpacket.SetSignerReq) (*pbredpacket.SetSignerResp, error) {
	if req.SignerAddress == "" {
		return nil, errs.ErrArgs.WrapMsg("signer_address is required")
	}
	if s.chainClient != nil {
		log.ZInfo(ctx, "redpacket admin setSigner (eth mock)", "signerAddress", req.SignerAddress)
		return &pbredpacket.SetSignerResp{Message: "signer address updated successfully"}, nil
	}
	if s.tronClient != nil {
		if _, err := s.tronClient.SendAdminTransaction(ctx, "setSigner", req.SignerAddress); err != nil {
			return nil, errs.ErrInternalServer.WrapMsg("setSigner failed: " + err.Error())
		}
		return &pbredpacket.SetSignerResp{Message: "signer address updated successfully"}, nil
	}
	return nil, errs.ErrInternalServer.WrapMsg("no blockchain client configured")
}

func (s *redPacketServer) SetToken(ctx context.Context, req *pbredpacket.SetTokenReq) (*pbredpacket.SetTokenResp, error) {
	if req.TokenAddress == "" {
		return nil, errs.ErrArgs.WrapMsg("token_address is required")
	}

	minAmountBig := new(big.Int)
	if req.MinAmount != "" {
		minAmountBig.SetString(req.MinAmount, 10)
	}

	if s.chainClient != nil {
		log.ZInfo(ctx, "redpacket admin setToken (eth mock)",
			"tokenAddress", req.TokenAddress,
			"allowed", req.Allowed,
			"minAmount", req.MinAmount,
		)
		return &pbredpacket.SetTokenResp{Message: "token configuration updated"}, nil
	}
	if s.tronClient != nil {
		if _, err := s.tronClient.SendAdminTransaction(ctx, "setAllowedToken", req.TokenAddress, req.Allowed, minAmountBig); err != nil {
			return nil, errs.ErrInternalServer.WrapMsg("setAllowedToken failed: " + err.Error())
		}
		return &pbredpacket.SetTokenResp{Message: "token configuration updated"}, nil
	}
	return nil, errs.ErrInternalServer.WrapMsg("no blockchain client configured")
}

func (s *redPacketServer) SetExpiry(ctx context.Context, req *pbredpacket.SetExpiryReq) (*pbredpacket.SetExpiryResp, error) {
	if req.ExpirySeconds <= 0 {
		return nil, errs.ErrArgs.WrapMsg("expiry_seconds must be positive")
	}
	if s.chainClient != nil {
		log.ZInfo(ctx, "redpacket admin setExpiry (eth mock)", "expirySeconds", req.ExpirySeconds)
		return &pbredpacket.SetExpiryResp{Message: "expiry duration updated"}, nil
	}
	if s.tronClient != nil {
		if _, err := s.tronClient.SendAdminTransaction(ctx, "setDefaultExpiryDuration", req.ExpirySeconds); err != nil {
			return nil, errs.ErrInternalServer.WrapMsg("setDefaultExpiryDuration failed: " + err.Error())
		}
		return &pbredpacket.SetExpiryResp{Message: "expiry duration updated"}, nil
	}
	return nil, errs.ErrInternalServer.WrapMsg("no blockchain client configured")
}

func (s *redPacketServer) SetAllowAllTokens(ctx context.Context, req *pbredpacket.SetAllowAllTokensReq) (*pbredpacket.SetAllowAllTokensResp, error) {
	if s.chainClient != nil {
		log.ZInfo(ctx, "redpacket admin setAllowAllTokens (eth mock)", "allowAll", req.AllowAll)
		return &pbredpacket.SetAllowAllTokensResp{Message: "allow all tokens setting updated"}, nil
	}
	if s.tronClient != nil {
		if _, err := s.tronClient.SendAdminTransaction(ctx, "setAllowAllTokens", req.AllowAll); err != nil {
			return nil, errs.ErrInternalServer.WrapMsg("setAllowAllTokens failed: " + err.Error())
		}
		return &pbredpacket.SetAllowAllTokensResp{Message: "allow all tokens setting updated"}, nil
	}
	return nil, errs.ErrInternalServer.WrapMsg("no blockchain client configured")
}

func (s *redPacketServer) SetNativeTokenEnabled(ctx context.Context, req *pbredpacket.SetNativeTokenEnabledReq) (*pbredpacket.SetNativeTokenEnabledResp, error) {
	if s.chainClient != nil {
		log.ZInfo(ctx, "redpacket admin setNativeTokenEnabled (eth mock)", "enabled", req.Enabled)
		return &pbredpacket.SetNativeTokenEnabledResp{Message: "native token setting updated"}, nil
	}
	if s.tronClient != nil {
		if _, err := s.tronClient.SendAdminTransaction(ctx, "setNativeTokenEnabled", req.Enabled); err != nil {
			return nil, errs.ErrInternalServer.WrapMsg("setNativeTokenEnabled failed: " + err.Error())
		}
		return &pbredpacket.SetNativeTokenEnabledResp{Message: "native token setting updated"}, nil
	}
	return nil, errs.ErrInternalServer.WrapMsg("no blockchain client configured")
}

func (s *redPacketServer) ParseTxEvents(ctx context.Context, req *pbredpacket.ParseTxEventsReq) (*pbredpacket.ParseTxEventsResp, error) {
	if req.TxHash == "" {
		return nil, errs.ErrArgs.WrapMsg("tx_hash is required")
	}

	if req.Chain == "tron" && s.tronClient != nil {
		return &pbredpacket.ParseTxEventsResp{
			Chain:  "tron",
			TxHash: req.TxHash,
			Note:   "TRON event parsing not fully implemented in this version",
		}, nil
	}

	if s.chainClient != nil {
		txHashBytes := common.HexToHash(req.TxHash)
		events, err := s.chainClient.ParseTransactionReceipt(ctx, txHashBytes)
		if err != nil {
			return nil, errs.ErrInternalServer.WrapMsg("parse tx receipt failed: " + err.Error())
		}

		out := make([]*pbredpacket.ParsedEvent, 0, len(events))
		for _, e := range events {
			data := make(map[string]string, len(e.Data))
			for k, v := range e.Data {
				data[k] = fmt.Sprintf("%v", v)
			}
			out = append(out, &pbredpacket.ParsedEvent{
				Name: e.Name,
				Data: data,
			})
		}
		return &pbredpacket.ParseTxEventsResp{
			Chain:  "eth",
			TxHash: req.TxHash,
			Events: out,
		}, nil
	}

	return nil, errs.ErrInternalServer.WrapMsg("no client available for chain: " + req.Chain)
}
