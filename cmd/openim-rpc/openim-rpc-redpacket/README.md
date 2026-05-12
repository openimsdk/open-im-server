# RedPacket RPC Service

A Web3 Red Packet RPC service that has been migrated to the standard OpenIM
service layout: gRPC over `protocol/redpacket`, MongoDB via the `mgo` +
`controller` pattern, and command/discovery wiring through `pkg/common/cmd`
and `pkg/common/startrpc`.

For HTTP access, the service is exposed by the API gateway under `/redpacket/*`
(see `internal/api/redpacket.go`).

## Layout

```
.
├── main.go                                              # cmd.NewRedPacketRpcCmd().Exec()
├── README.md
├── backend-api.md                                       # Legacy API docs, kept for reference
├── client-integration-guide.md                          # Legacy integration docs, kept for reference
├── red-packet-go-backend-eth-tron.md                    # Architecture / chain integration design
└── redpacket-web3-integration-design.md                 # Web3 integration design
```

The actual implementation lives in:

- `protocol/redpacket/redpacket.proto`                   – gRPC contract
- `pkg/common/storage/model/redpacket.go`                – Mongo BSON models
- `pkg/common/storage/database/redpacket.go`             – DAO interfaces
- `pkg/common/storage/database/mgo/redpacket.go`         – Mongo DAO impl
- `pkg/common/storage/controller/redpacket.go`           – Aggregated database façade
- `pkg/common/cmd/rpc_redpacket.go`                      – Cobra entry, startrpc bootstrap
- `internal/rpc/redpacket/`                              – gRPC service, chain client, indexers
- `internal/api/redpacket.go`                            – Gin gateway handlers
- `config/openim-rpc-redpacket.yml`                      – Service configuration

## Features

- ✅ Create red packet orders + on-chain `Created` callback reconciliation
- ✅ Red packet detail query (with full claim history)
- ✅ Claim signature issuance using the contract's `getSignMessage(...)`
- ✅ Claim result reporting + idempotent persistence by tx hash
- ✅ EVM event indexer (claim / refund)
- ✅ TRON full-node JSON-RPC integration scaffold
- ✅ EVM SIWE-style wallet binding (challenge / sign / confirm)
- ✅ Admin endpoints (signer / allowed token / expiry / allow-all-tokens / native-token)

## Configuration

See `config/openim-rpc-redpacket.yml` (alongside other OpenIM RPC configs).

```yaml
rpc:
  registerIP: ""
  listenIP: 0.0.0.0
  autoSetPorts: false
  ports: [10560]

prometheus:
  enable: false
  ports: [12560]

chain:                # Optional — leave rpcURL empty to disable EVM
  rpcURL: ""
  contractAddress: ""
  chainID: 0
  signerPrivateKey: ""
  configAdminPrivateKey: ""

tron:                 # Optional — leave fullNodeURL empty to disable TRON
  fullNodeURL: ""
  contractBase58: ""
  ownerBase58: ""
  privateKeyHex: ""
  feeLimit: 100000000

indexer:
  pollInterval: 5
```

`config/share.yml` registers the service name as `redPacket`.

## Limitations / TODO

- TRON `ConfirmWalletBind` signature verification is not yet implemented and
  returns `not implemented`.
- TRON event decoding in `chain/tron_indexer.go` is still a scaffold and only
  identifies events by topic-0; payload decoding will be added once the
  contract event signatures are finalized.
- Admin endpoints (`/redpacket/admin/*`) currently mirror the legacy mock
  behaviour for EVM and only forward live calls on TRON.

See `backend-api.md`, `client-integration-guide.md`, and the design docs for
detailed specifications.
