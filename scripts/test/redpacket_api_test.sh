#!/usr/bin/env bash
# ============================================================
# 红包 HTTP 接口测试：create_order / created_callback
#
# 路由（与 internal/api/router.go 一致）：
#   POST ${HOST}/redpacket/create_order
#   POST ${HOST}/redpacket/created_callback
#
# 鉴权：两接口均不在白名单，需在 Header 携带 token（见 protocol/constant constant.Token = "token"）。
# 追踪：Header 需携带 operationID。
#
# 依赖：curl、jq；自动拉管理员 token 时另需 python3。
#
# 用法示例：
#   chmod +x scripts/test/redpacket_api_test.sh
#   GROUP_ID=你的群ID USER_ID=你的用户ID ./scripts/test/redpacket_api_test.sh
#   ./scripts/test/redpacket_api_test.sh --host http://127.0.0.1:10002 --group-id xxx --try-callback
#
# 说明：
#   - create_order 在 packetType=0（拼手气固定份）时要求 scopeType=GROUP 且当前用户在该群内。
#   - 若 RPC 侧未配置 EVM chain client，created_callback 可走「离线」路径：传任意非空 txHash，
#     并在 body 中提供与订单一致的 packetID（见 internal/rpc/redpacket resolveCreatedPacket EVM 分支）。
#   - 生产环境若已接链，created_callback 需真实上链交易哈希，此时请自行设置 TX_HASH / PACKET_ID。
# ============================================================

set -euo pipefail

HOST="${HOST:-http://127.0.0.1:10002}"
USER_ID="${USER_ID:-5694418935}"
PLATFORM_ID="${PLATFORM_ID:-2}"
ADMIN_TOKEN="${ADMIN_TOKEN:-}"
OPENIM_SECRET="${OPENIM_SECRET:-openIM123}"
ADMIN_USER_ID="${ADMIN_USER_ID:-imAdmin}"
TOKEN="${TOKEN:-}"

GROUP_ID="${GROUP_ID:-}"
CHAIN_TYPE="${CHAIN_TYPE:-EVM}"
CHAIN_ID="${CHAIN_ID:-0}"
SCOPE_TYPE="${SCOPE_TYPE:-GROUP}"
PACKET_TYPE="${PACKET_TYPE:-0}"
CREATOR_WALLET="${CREATOR_WALLET:-0x0000000000000000000000000000000000000001}"
TOKEN_ADDR="${TOKEN_ADDR:-0x0000000000000000000000000000000000000000}"
TOTAL_AMOUNT="${TOTAL_AMOUNT:-100}"
TOTAL_SHARES="${TOTAL_SHARES:-5}"
EXPIRY_AT="${EXPIRY_AT:-0}"
REMARK="${REMARK:-api-test}"

TRY_CALLBACK="${TRY_CALLBACK:-0}"
TX_HASH="${TX_HASH:-0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa}"
CALLBACK_PACKET_ID="${CALLBACK_PACKET_ID:-}"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --host) HOST="$2"; shift 2 ;;
    --user-id) USER_ID="$2"; shift 2 ;;
    --platform-id) PLATFORM_ID="$2"; shift 2 ;;
    --group-id) GROUP_ID="$2"; shift 2 ;;
    --token) TOKEN="$2"; shift 2 ;;
    --try-callback) TRY_CALLBACK="1"; shift ;;
    *)
      echo "未知参数: $1"
      exit 1
      ;;
  esac
done

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "缺少依赖命令: $1"
    exit 1
  }
}

need_cmd curl
need_cmd jq

op_id() {
  echo "redpacket-test-$$-$(date +%s%N)"
}

get_admin_token() {
  local uid body resp token last_resp
  local -a candidates=("${ADMIN_USER_ID}" "openIM123456" "imAdmin")
  last_resp=""

  for uid in "${candidates[@]}"; do
    body="{\"secret\":\"${OPENIM_SECRET}\",\"userID\":\"${uid}\"}"
    resp="$(curl -sS -X POST "${HOST}/auth/get_admin_token" \
      -H "Content-Type: application/json" \
      -H "operationID: $(op_id)" \
      -d "$body")"
    last_resp="$resp"

    token="$(python3 - <<'PY' "$resp"
import json
import sys

raw = sys.argv[1]
try:
    obj = json.loads(raw)
except Exception:
    print("")
    raise SystemExit(0)

token = ""
if isinstance(obj, dict):
    data = obj.get("data")
    if isinstance(data, dict):
        token = data.get("token") or data.get("Token") or ""
    if not token:
        token = obj.get("token") or obj.get("Token") or ""
print(token)
PY
)"
    if [[ -n "$token" ]]; then
      echo "自动获取管理员 token 成功，userID=${uid}" >&2
      printf '%s' "$token"
      return 0
    fi
  done

  echo "get_admin_token raw response: $last_resp" >&2
  echo "自动获取管理员 token 失败，请检查 HOST/OPENIM_SECRET/ADMIN_USER_ID 或直接设置 ADMIN_TOKEN" >&2
  exit 1
}

resolve_user_token() {
  if [[ -n "${TOKEN}" ]]; then
    echo "使用环境变量/参数 TOKEN（跳过 get_user_token）" >&2
    return 0
  fi

  need_cmd python3

  if [[ -z "${ADMIN_TOKEN}" ]]; then
    echo "==> ADMIN_TOKEN 未设置，尝试自动获取管理员 token" >&2
    ADMIN_TOKEN="$(get_admin_token)"
  fi

  echo "==> 获取用户 token（userID=${USER_ID}）" >&2
  local TOKEN_RESP
  TOKEN_RESP=$(curl -sS -X POST \
    -H "Content-Type: application/json" \
    -H "operationID: $(op_id)" \
    -H "token: ${ADMIN_TOKEN}" \
    -d "{\"userID\":\"${USER_ID}\",\"platformID\":${PLATFORM_ID}}" \
    "${HOST}/auth/get_user_token")

  local ERR_CODE
  ERR_CODE=$(echo "${TOKEN_RESP}" | jq -r '.errCode // "null"')
  if [[ "${ERR_CODE}" != "0" ]]; then
    echo "获取用户 token 失败: ${TOKEN_RESP}" >&2
    exit 1
  fi
  TOKEN=$(echo "${TOKEN_RESP}" | jq -r '.data.token // empty')
  if [[ -z "${TOKEN}" ]]; then
    echo "token 为空: ${TOKEN_RESP}" >&2
    exit 1
  fi
  echo "用户 token 获取成功" >&2
}

if [[ -z "${GROUP_ID}" ]]; then
  echo "错误：未设置 GROUP_ID。固定份红包（packetType=0）需要 scopeType=GROUP 且 group_id 非空。" >&2
  echo "示例：GROUP_ID=你的群ID USER_ID=在群内的用户 ./scripts/test/redpacket_api_test.sh" >&2
  exit 1
fi

resolve_user_token

echo "==> POST /redpacket/create_order"
CREATE_BODY=$(jq -n \
  --arg chainType "${CHAIN_TYPE}" \
  --argjson chainID "${CHAIN_ID}" \
  --arg groupID "${GROUP_ID}" \
  --arg scopeType "${SCOPE_TYPE}" \
  --argjson packetType "${PACKET_TYPE}" \
  --arg token "${TOKEN_ADDR}" \
  --arg totalAmount "${TOTAL_AMOUNT}" \
  --argjson totalShares "${TOTAL_SHARES}" \
  --argjson expiryAt "${EXPIRY_AT}" \
  --arg remark "${REMARK}" \
  --arg creatorWallet "${CREATOR_WALLET}" \
  '{
    chainType: $chainType,
    chainID: $chainID,
    groupID: $groupID,
    scopeType: $scopeType,
    packetType: $packetType,
    token: $token,
    totalAmount: $totalAmount,
    totalShares: $totalShares,
    expiryAt: $expiryAt,
    remark: $remark,
    creatorWallet: $creatorWallet
  }')

CREATE_RESP=$(curl -sS -X POST \
  -H "Content-Type: application/json" \
  -H "operationID: $(op_id)" \
  -H "token: ${TOKEN}" \
  -d "${CREATE_BODY}" \
  "${HOST}/redpacket/create_order")

echo "${CREATE_RESP}" | jq .

CO_ERR=$(echo "${CREATE_RESP}" | jq -r '.errCode // "null"')
if [[ "${CO_ERR}" != "0" ]]; then
  echo "create_order 失败（errCode=${CO_ERR}）。请确认 USER_ID/TOKEN 对应用户在 GROUP_ID 群内，且 totalAmount 可被 totalShares 整除（固定份）。" >&2
  exit 1
fi

BIZ_ID=$(echo "${CREATE_RESP}" | jq -r '.data.bizID // empty')
if [[ -z "${BIZ_ID}" ]]; then
  echo "create_order 返回 errCode=0 但 data.bizID 为空: ${CREATE_RESP}" >&2
  exit 1
fi
echo "create_order 成功，bizID=${BIZ_ID}"

if [[ "${TRY_CALLBACK}" != "1" ]]; then
  echo "==> 未调用 created_callback（设置 TRY_CALLBACK=1 或传入 --try-callback 以继续）"
  echo "    离线 EVM：可设置 CALLBACK_PACKET_ID（默认用时间戳十进制字符串）；TX_HASH 可用环境变量 TX_HASH 覆盖。"
  exit 0
fi

if [[ -z "${CALLBACK_PACKET_ID}" ]]; then
  CALLBACK_PACKET_ID="$(date +%s)"
fi

echo "==> POST /redpacket/created_callback（bizID=${BIZ_ID}, packetID=${CALLBACK_PACKET_ID}）"
CALLBACK_BODY=$(jq -n \
  --arg bizID "${BIZ_ID}" \
  --arg txHash "${TX_HASH}" \
  --arg packetID "${CALLBACK_PACKET_ID}" \
  --arg groupID "${GROUP_ID}" \
  --arg scopeType "${SCOPE_TYPE}" \
  '{
    bizID: $bizID,
    txHash: $txHash,
    packetID: $packetID,
    groupID: $groupID,
    scopeType: $scopeType
  }')

CALLBACK_RESP=$(curl -sS -X POST \
  -H "Content-Type: application/json" \
  -H "operationID: $(op_id)" \
  -H "token: ${TOKEN}" \
  -d "${CALLBACK_BODY}" \
  "${HOST}/redpacket/created_callback")

echo "${CALLBACK_RESP}" | jq .

CB_ERR=$(echo "${CALLBACK_RESP}" | jq -r '.errCode // "null"')
if [[ "${CB_ERR}" != "0" ]]; then
  echo "created_callback 失败（errCode=${CB_ERR}）。若已配置链上客户端，请使用真实交易哈希或关闭 TRY_CALLBACK。" >&2
  exit 1
fi

echo "created_callback 成功，红包状态应已更新为 ACTIVE（视部署与链配置而定）。"
echo "测试通过: /redpacket/create_order + /redpacket/created_callback"
