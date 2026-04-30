# RedPacket Backend Service

A Web3 Red Packet service supporting Ethereum and TRON, following the design documents:

- `backend-api.md` - API specifications
- `client-integration-guide.md` - Frontend / gateway integration guide
- `redpacket-web3-integration-design.md` - Architecture and flows
- `red-packet-go-backend-eth-tron.md` - Blockchain integration details

## Features

- ✅ Create red packet orders (`/api/redpacket/create-order`)
- ✅ Created callback for on-chain transaction results
- ✅ Red packet detail query with claim history
- ✅ Claim signature issuance (`/api/redpacket/claim-sign`)
- ✅ Claim result reporting
- ✅ SQLite/MySQL support
- ✅ EVM signature generation via `getSignMessage(...)`
- ✅ Basic EVM event indexing for claim/refund synchronization
- ✅ Idempotent claim/refund persistence by transaction hash
- ✅ Admin configuration endpoints

## Current Status

This service is runnable and suitable for continued iteration, but it is not yet fully production-complete.

Working well now:

- EVM-side claim signing uses the real `authNonce` in the digest
- Claim pre-checks cover packet existence, active status, expiry, and already-claimed cases
- EVM ABI and event parsing are aligned with the current contract events
- Claim and refund events can be persisted idempotently

Still incomplete:

- ETH admin endpoints are still mostly mock behavior
- `PacketCreated` indexing is not yet fully wired for automatic order reconciliation
- TRON `getSignMessage` flow is not complete
- TRON event decoding is still a scaffold
- Admin APIs still need authentication and audit controls

## Quick Start

```bash
cd cmd/openim-rpc/openim-rpc-redpacket

# 1. Configure (optional)
cp config/config.yaml config/config.yaml.bak
# Edit config/config.yaml with your blockchain settings

# 2. Build and run
go run .

# Or build binary
go build -o redpacket .
./redpacket
```

Service will start on `http://localhost:8080`

## Test the API

```bash
# Health check
curl http://localhost:8080/health

# Create red packet
curl -X POST http://localhost:8080/api/redpacket/create-order \
  -H "Content-Type: application/json" \
  -d '{
    "creator_user_id": "u1001",
    "creator_wallet": "0x1111111111111111111111111111111111111111",
    "packet_type": 1,
    "total_amount": "1000000000000000000",
    "total_shares": 10
  }'
```

## Project Structure

```
.
├── config/              # Configuration
├── internal/
│   ├── handler/         # HTTP handlers (Gin)
│   ├── model/           # Database models (GORM)
│   ├── repository/      # Data access layer
│   ├── service/         # Business logic
│   └── chain/           # Blockchain integration and event indexing
├── pkg/resp/            # Response helpers
├── router/              # Route definitions
├── main.go
├── go.mod
└── README.md
```

## Recommended Next Steps

1. Implement real ETH admin transactions for signer/token/expiry configuration
2. Finish `PacketCreated` indexing and automatic order reconciliation
3. Complete TRON `getSignMessage` and reliable event decoding
4. Add authentication, audit, and rate limiting for sensitive endpoints
5. Extend end-to-end test coverage

See the three design documents for detailed specifications.

## API Documentation

See:

- `backend-api.md` for complete API reference with request / response examples
- `client-integration-guide.md` for frontend, wallet-binding, and claim-sign integration steps
