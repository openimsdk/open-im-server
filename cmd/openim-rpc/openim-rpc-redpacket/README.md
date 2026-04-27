# RedPacket Backend Service

A Web3 Red Packet service supporting Ethereum and TRON, following the design documents:

- `backend-api.md` - API specifications
- `redpacket-web3-integration-design.md` - Architecture and flows
- `red-packet-go-backend-eth-tron.md` - Blockchain integration details

## Features

- ✅ Create red packet orders (`/api/redpacket/create-order`)
- ✅ Created callback for on-chain transaction results
- ✅ Red packet detail query with claim history
- ✅ Claim signature issuance (`/api/redpacket/claim-sign`)
- ✅ Claim result reporting
- ✅ SQLite/MySQL support
- ✅ Blockchain signature logic ready for ETH/TRON
- ✅ Admin configuration endpoints

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
│   └── chain/           # Blockchain integration (to be expanded)
├── pkg/resp/            # Response helpers
├── router/              # Route definitions
├── main.go
├── go.mod
└── README.md
```

## Next Steps (from design docs)

1. **Full Blockchain Integration**
   - Implement `ChainClient` for ETH and TRON
   - Add event indexer for `PacketCreated`, `PacketClaimed`, `PacketRefunded`
   - Implement proper signature generation using `getSignMessage`

2. **Advanced Features**
   - Admin configuration APIs (`setSigner`, `setToken`, etc.)
   - Refund logic
   - Rate limiting and authentication
   - Monitoring and metrics

3. **Production**
   - Add proper authentication middleware
   - Configure production database
   - Set up monitoring and logging
   - Deploy with Docker/K8s

See the three design documents for detailed specifications.

## API Documentation

See `backend-api.md` for complete API reference with examples.
