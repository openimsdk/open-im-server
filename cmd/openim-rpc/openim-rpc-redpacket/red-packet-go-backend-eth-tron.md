# 红包 Go 后台对接（ETH + TRON）

这份文档按你的需求给出三部分：
- 后端签名（`claim` 鉴权签名，ETH/TRON 通用）
- ETH 后台调用 + 通过 `txhash` 解析事件
- TRON 后台调用流程 + 通过 `txhash` 解析事件

说明：以下签名逻辑严格对应当前合约 `RedPacketBase` 的 `getSignMessage/claim`。

---

## 1. 依赖

```bash
go get github.com/ethereum/go-ethereum@v1.14.12
```

---

## 2. 关键合约事实（当前仓库）

- 签名结构体：
  `Claim(uint256 packetId,address claimer,uint256 authNonce,uint256 randomSeed,uint256 deadline)`
- 领取函数：
  `claim(packetId, authNonce, randomSeed, deadline, signature)`
- 重点事件：
  - `PacketCreated(uint256,address,uint8,address,uint256,uint256,uint256)`
  - `PacketClaimed(uint256,address,uint256,uint256,uint256,uint256)`
  - `PacketRefunded(uint256,address,address,uint256)`

---

## 3. Go：后端 claim 签名（ETH/TRON 通用）

合约里验签是 `ecrecover(getSignMessage(...), v, r, s)`，所以后端要对 `digest` 做裸签名，不要加 `personal_sign` 前缀。

```go
package redpacket

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// SignClaimDigest 对合约返回的 digest 做裸签，返回 65 字节签名（r||s||v）
func SignClaimDigest(priv *ecdsa.PrivateKey, digest [32]byte) ([]byte, error) {
	sig, err := crypto.Sign(digest[:], priv)
	if err != nil {
		return nil, err
	}
	// go-ethereum 返回 v 为 0/1；EVM 合约通常期望 27/28
	sig[64] += 27
	return sig, nil
}

// RecoverAndCheckSigner 本地自检（可选）
func RecoverAndCheckSigner(digest [32]byte, sig []byte, expected common.Address) error {
	if len(sig) != 65 {
		return fmt.Errorf("invalid sig length: %d", len(sig))
	}
	cpy := make([]byte, 65)
	copy(cpy, sig)
	if cpy[64] >= 27 {
		cpy[64] -= 27
	}
	pub, err := crypto.SigToPub(digest[:], cpy)
	if err != nil {
		return err
	}
	got := crypto.PubkeyToAddress(*pub)
	if got != expected {
		return fmt.Errorf("signer mismatch, got=%s want=%s", got.Hex(), expected.Hex())
	}
	return nil
}

// BuildClaimTypeHash 仅当你要本地复算 digest 时才需要。
func BuildClaimTypeHash() common.Hash {
	return crypto.Keccak256Hash([]byte("Claim(uint256 packetId,address claimer,uint256 authNonce,uint256 randomSeed,uint256 deadline)"))
}

// BuildClaimStructHash 本地复算 structHash（可选）。
func BuildClaimStructHash(packetId *big.Int, claimer common.Address, authNonce, randomSeed, deadline *big.Int) common.Hash {
	typeHash := BuildClaimTypeHash()
	encoded := make([]byte, 0, 32*6)
	encoded = append(encoded, typeHash.Bytes()...)
	encoded = append(encoded, common.LeftPadBytes(packetId.Bytes(), 32)...)
	encoded = append(encoded, common.LeftPadBytes(claimer.Bytes(), 32)...)
	encoded = append(encoded, common.LeftPadBytes(authNonce.Bytes(), 32)...)
	encoded = append(encoded, common.LeftPadBytes(randomSeed.Bytes(), 32)...)
	encoded = append(encoded, common.LeftPadBytes(deadline.Bytes(), 32)...)
	return crypto.Keccak256Hash(encoded)
}
```

生产建议：
- 最稳妥方式是先链上调用 `getSignMessage(...)` 拿 `digest`，再签名。
- `authNonce` 必须按 `claimer` 做幂等和防重。
- `deadline` 建议 5~30 分钟。

---

## 4. Go：ETH 后台调用 + txhash 解析事件

### 4.1 通过 txhash 解析 `PacketCreated/PacketClaimed/PacketRefunded`

```go
package redpacket

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ParsedEvent struct {
	Name string
	Data map[string]any
}

func ParseEthEventsByTxHash(ctx context.Context, rpcURL, txHashHex, contractABIJSON string) ([]ParsedEvent, error) {
	cli, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	txHash := common.HexToHash(txHashHex)
	rcpt, err := cli.TransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, err
	}

	parsedABI, err := abi.JSON(strings.NewReader(contractABIJSON))
	if err != nil {
		return nil, err
	}

	var out []ParsedEvent
	for _, lg := range rcpt.Logs {
		ev, ok := eventFromLog(parsedABI, lg)
		if ok {
			out = append(out, ev)
		}
	}
	return out, nil
}

func eventFromLog(parsedABI abi.ABI, lg *types.Log) (ParsedEvent, bool) {
	if len(lg.Topics) == 0 {
		return ParsedEvent{}, false
	}
	for name, e := range parsedABI.Events {
		if e.ID != lg.Topics[0] {
			continue
		}
		vals := map[string]any{}

		// 非 indexed 参数
		nonIndexed, err := e.Inputs.NonIndexed().Unpack(lg.Data)
		if err != nil {
			return ParsedEvent{}, false
		}
		n := 0
		idxTopic := 1
		for _, input := range e.Inputs {
			if input.Indexed {
				if idxTopic >= len(lg.Topics) {
					return ParsedEvent{}, false
				}
				vals[input.Name] = decodeIndexedTopic(input.Type, lg.Topics[idxTopic])
				idxTopic++
			} else {
				vals[input.Name] = nonIndexed[n]
				n++
			}
		}
		return ParsedEvent{Name: name, Data: vals}, true
	}
	return ParsedEvent{}, false
}

func decodeIndexedTopic(t abi.Type, topic common.Hash) any {
	switch t.T {
	case abi.AddressTy:
		return common.BytesToAddress(topic.Bytes()[12:])
	default:
		return topic
	}
}

func PrettyPrintEvents(events []ParsedEvent) string {
	b, _ := json.MarshalIndent(events, "", "  ")
	return string(b)
}

func MustReadABIFromArtifact(artifactJSON []byte) (string, error) {
	var raw map[string]any
	if err := json.Unmarshal(artifactJSON, &raw); err != nil {
		return "", err
	}
	abiObj, ok := raw["abi"]
	if !ok {
		return "", fmt.Errorf("abi field not found")
	}
	abiBytes, err := json.Marshal(abiObj)
	if err != nil {
		return "", err
	}
	return string(abiBytes), nil
}
```

### 4.2 ETH 创建/领取调用（示意）

建议用 `abigen` 生成 Go binding 后调用（最稳）。

`abigen` 示例：
```bash
abigen --abi abi/contracts/RedPacket.sol/RedPacket.json --pkg redpacket --type RedPacket --out redpacket_binding.go
```

调用流程：
1. `createFixedPacket/createRandomPacket/createTransfer` 发交易
2. 拿到 `txHash` 后轮询 receipt
3. 用上面的 `ParseEthEventsByTxHash` 解出 `PacketCreated`，拿到 `packetId`
4. 后端签名下发给前端后，前端/后端发 `claim`
5. 用 `PacketClaimed.amount` 作为最终到账金额

---

## 5. Go：TRON 后台调用 + txhash 解析事件

TRON 的 EVM 合约事件最终也是 topic/data 结构，因此事件解码可复用 EVM ABI。

### 5.1 通过 txhash 解析 TRON 事件（推荐走 `/wallet/gettransactioninfobyid`）

```go
package redpacket

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type tronTxInfoResp struct {
	ID  string `json:"id"`
	Log []struct {
		Address string   `json:"address"` // 合约地址hex(无0x)
		Topics  []string `json:"topics"`  // topic hex(无0x)
		Data    string   `json:"data"`    // data hex(无0x)
	} `json:"log"`
}

func ParseTronEventsByTxHash(ctx context.Context, tronFullNodeURL, txID, contractABIJSON string) ([]ParsedEvent, error) {
	body := map[string]string{"value": txID}
	buf, _ := json.Marshal(body)

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, tronFullNodeURL+"/wallet/gettransactioninfobyid", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("tron http %d: %s", resp.StatusCode, string(raw))
	}

	var info tronTxInfoResp
	if err := json.Unmarshal(raw, &info); err != nil {
		return nil, err
	}

	parsedABI, err := abi.JSON(strings.NewReader(contractABIJSON))
	if err != nil {
		return nil, err
	}

	out := make([]ParsedEvent, 0, len(info.Log))
	for _, lg := range info.Log {
		if len(lg.Topics) == 0 {
			continue
		}
		topic0 := common.HexToHash("0x" + lg.Topics[0])

		for name, e := range parsedABI.Events {
			if e.ID != topic0 {
				continue
			}
			vals := map[string]any{}

			dataBytes, err := hex.DecodeString(strings.TrimPrefix(lg.Data, "0x"))
			if err != nil {
				return nil, err
			}
			nonIndexed, err := e.Inputs.NonIndexed().Unpack(dataBytes)
			if err != nil {
				return nil, err
			}

			n := 0
			idxTopic := 1
			for _, input := range e.Inputs {
				if input.Indexed {
					if idxTopic >= len(lg.Topics) {
						return nil, fmt.Errorf("missing indexed topic for event %s", name)
					}
					t := common.HexToHash("0x" + lg.Topics[idxTopic])
					vals[input.Name] = decodeIndexedTopic(input.Type, t)
					idxTopic++
				} else {
					vals[input.Name] = nonIndexed[n]
					n++
				}
			}

			out = append(out, ParsedEvent{Name: name, Data: vals})
			break
		}
	}

	return out, nil
}
```

### 5.2 TRON 后台调用流程（实践）

1. 组装 ABI 参数（与 ETH 一样）
2. 调用 TRON FullNode 的 `trigger*contract` 生成未签名交易
3. 用托管私钥签名交易并广播
4. 根据返回 `txID` 调用上面的 `ParseTronEventsByTxHash` 解事件

说明：TRON 发交易接口在不同节点服务（TronGrid/自建 FullNode/SDK 封装）字段细节略有差异，建议你在项目里固定一种（推荐固定 TronGrid 或 gotron-sdk 版本），避免线上环境差异。

---

## 6. 合约参数设置（管理员）

需要 `CONFIG_ADMIN_ROLE` 的函数：
- `setSigner(address signer)`
- `setAllowAllTokens(bool allowAllTokens)`
- `setNativeTokenEnabled(bool enabled)`
- `setAllowedToken(address token, bool allowed, uint256 minShareAmount)`
- `setDefaultExpiryDuration(uint256 duration)`

对应配置事件（可按 `txhash` 解析校验）：
- `SignerUpdated(oldSigner, newSigner)`
- `AllowAllTokensUpdated(allowAllTokens)`
- `NativeTokenEnabledUpdated(enabled)`
- `AllowedTokenUpdated(token, allowed, minShareAmount)`
- `DefaultExpiryDurationUpdated(duration)`

### 6.1 ETH：Go 设置合约参数（通用写法）

```go
package redpacket

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// SendEthAdminTx 通用管理员写调用：
// method 例如 "setNativeTokenEnabled"
// args 对应函数参数
func SendEthAdminTx(
	ctx context.Context,
	rpcURL string,
	contractAddr common.Address,
	priv *ecdsa.PrivateKey,
	contractABIJSON string,
	method string,
	args ...any,
) (common.Hash, error) {
	cli, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return common.Hash{}, err
	}
	defer cli.Close()

	from := crypto.PubkeyToAddress(priv.PublicKey)
	nonce, err := cli.PendingNonceAt(ctx, from)
	if err != nil {
		return common.Hash{}, err
	}
	chainID, err := cli.NetworkID(ctx)
	if err != nil {
		return common.Hash{}, err
	}
	gasPrice, err := cli.SuggestGasPrice(ctx)
	if err != nil {
		return common.Hash{}, err
	}

	parsedABI, err := abi.JSON(strings.NewReader(contractABIJSON))
	if err != nil {
		return common.Hash{}, err
	}
	data, err := parsedABI.Pack(method, args...)
	if err != nil {
		return common.Hash{}, err
	}

	msg := ethereum.CallMsg{
		From: from, To: &contractAddr, Data: data, Value: big.NewInt(0),
	}
	gasLimit, err := cli.EstimateGas(ctx, msg)
	if err != nil {
		return common.Hash{}, err
	}

	tx := types.NewTransaction(nonce, contractAddr, big.NewInt(0), gasLimit, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), priv)
	if err != nil {
		return common.Hash{}, err
	}
	if err = cli.SendTransaction(ctx, signedTx); err != nil {
		return common.Hash{}, err
	}
	return signedTx.Hash(), nil
}

// 例子：开启原生币、放开所有 token、设置 token 白名单与最小份额
func ExampleSetConfigEth(ctx context.Context, rpcURL, abiJSON, contractHex string, priv *ecdsa.PrivateKey, usdt common.Address) error {
	contract := common.HexToAddress(contractHex)

	tx1, err := SendEthAdminTx(ctx, rpcURL, contract, priv, abiJSON, "setNativeTokenEnabled", true)
	if err != nil {
		return err
	}
	fmt.Println("setNativeTokenEnabled tx:", tx1.Hex())

	tx2, err := SendEthAdminTx(ctx, rpcURL, contract, priv, abiJSON, "setAllowAllTokens", false)
	if err != nil {
		return err
	}
	fmt.Println("setAllowAllTokens tx:", tx2.Hex())

	tx3, err := SendEthAdminTx(ctx, rpcURL, contract, priv, abiJSON, "setAllowedToken", usdt, true, big.NewInt(1_000_000))
	if err != nil {
		return err
	}
	fmt.Println("setAllowedToken tx:", tx3.Hex())

	return nil
}
```

注意：`setAllowedToken(..., minShareAmount)` 的单位是 token 最小单位（例如 6 位精度 token，`1_000_000` 代表 1 个 token）。

### 6.2 TRON：Go 设置合约参数（FullNode HTTP）

TRON 推荐流程：`triggersmartcontract -> gettransactionsign -> broadcasttransaction`。

```go
package redpacket

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func encodeTronParams(abiJSON, method string, args ...any) (string, error) {
	parsed, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return "", err
	}
	m, ok := parsed.Methods[method]
	if !ok {
		return "", fmt.Errorf("method not found: %s", method)
	}
	packed, err := m.Inputs.Pack(args...)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(packed), nil
}

func postJSON(ctx context.Context, url string, body any, out any) error {
	b, _ := json.Marshal(body)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return fmt.Errorf("http %d: %s", resp.StatusCode, string(raw))
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return err
	}
	return nil
}

// SendTronAdminTx 示例：
// selector 例子 "setNativeTokenEnabled(bool)"
// methodName 例子 "setNativeTokenEnabled"
func SendTronAdminTx(
	ctx context.Context,
	fullNodeURL, ownerBase58, contractBase58, selector, methodName string,
	feeLimit int64,
	privateKeyHex string,
	abiJSON string,
	args ...any,
) (string, error) {
	paramHex, err := encodeTronParams(abiJSON, methodName, args...)
	if err != nil {
		return "", err
	}

	var triggerResp map[string]any
	err = postJSON(ctx, fullNodeURL+"/wallet/triggersmartcontract", map[string]any{
		"owner_address":    ownerBase58,
		"contract_address": contractBase58,
		"function_selector": selector,
		"parameter":        paramHex,
		"fee_limit":        feeLimit,
		"call_value":       0,
		"visible":          true,
	}, &triggerResp)
	if err != nil {
		return "", err
	}

	txObj, ok := triggerResp["transaction"]
	if !ok {
		return "", fmt.Errorf("transaction not found in trigger response")
	}

	var signedResp map[string]any
	err = postJSON(ctx, fullNodeURL+"/wallet/gettransactionsign", map[string]any{
		"transaction": txObj,
		"privateKey":  privateKeyHex,
	}, &signedResp)
	if err != nil {
		return "", err
	}

	var broadcastResp map[string]any
	err = postJSON(ctx, fullNodeURL+"/wallet/broadcasttransaction", signedResp, &broadcastResp)
	if err != nil {
		return "", err
	}
	if result, _ := broadcastResp["result"].(bool); !result {
		return "", fmt.Errorf("broadcast failed: %v", broadcastResp)
	}

	txid, _ := broadcastResp["txid"].(string)
	return txid, nil
}
```

调用示例：
- `setNativeTokenEnabled(true)`：
  `selector = "setNativeTokenEnabled(bool)"`，`methodName = "setNativeTokenEnabled"`，`args = true`
- `setAllowAllTokens(false)`：
  `selector = "setAllowAllTokens(bool)"`，`methodName = "setAllowAllTokens"`，`args = false`
- `setAllowedToken(token, true, 1_000_000)`：
  `selector = "setAllowedToken(address,bool,uint256)"`，`methodName = "setAllowedToken"`，`args = common.HexToAddress(tokenHexAddress), true, big.NewInt(1_000_000)`

安全建议：生产环境不要把私钥直接传给节点接口，建议改为本地离线签名或托管签名服务。

---

## 7. 最小落地建议（直接可用）

- 统一保存：`chain + txHash + packetId + eventName + rawEventJson`
- 创建成功：只认 `PacketCreated.packetId`
- 领取成功：只认 `PacketClaimed.amount`
- 退款成功：只认 `PacketRefunded.amount`
- 签名服务：`authNonce` 做地址维度去重；`deadline` 过期即废弃
