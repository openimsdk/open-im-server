#!/usr/bin/env bash
# ============================================================
# get_self_login_platforms 接口测试脚本
#
# 覆盖接口：
#   POST /auth/get_user_token
#   POST /user/get_self_login_platforms
#
# 说明：
#   本脚本仅做 HTTP 接口测试，不建立 WS 连接。
# ============================================================

set -euo pipefail

HOST="${HOST:-http://127.0.0.1:10002}"
USER_ID="${USER_ID:-5694418935}"
PLATFORM_ID="${PLATFORM_ID:-2}"
ADMIN_TOKEN="${ADMIN_TOKEN:-}"
OPENIM_SECRET="${OPENIM_SECRET:-openIM123}"
ADMIN_USER_ID="${ADMIN_USER_ID:-imAdmin}"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --host) HOST="$2"; shift 2 ;;
    --user-id) USER_ID="$2"; shift 2 ;;
    --platform-id) PLATFORM_ID="$2"; shift 2 ;;
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
  echo "self-login-platforms-test-$$-$(date +%s%N)"
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

if [[ -z "${ADMIN_TOKEN}" ]]; then
  echo "==> 1) ADMIN_TOKEN 未设置，尝试自动获取管理员 token"
  ADMIN_TOKEN="$(get_admin_token)"
fi

echo "==> 2) 获取用户 token"
TOKEN_RESP=$(curl -sS -X POST \
  -H "Content-Type: application/json" \
  -H "operationID: $(op_id)" \
  -H "token: ${ADMIN_TOKEN}" \
  -d "{\"userID\":\"${USER_ID}\",\"platformID\":${PLATFORM_ID}}" \
  "${HOST}/auth/get_user_token")

ERR_CODE=$(echo "${TOKEN_RESP}" | jq -r '.errCode // "null"')
if [[ "${ERR_CODE}" != "0" ]]; then
  echo "获取 token 失败: ${TOKEN_RESP}"
  exit 1
fi
TOKEN=$(echo "${TOKEN_RESP}" | jq -r '.data.token // empty')
if [[ -z "${TOKEN}" ]]; then
  echo "token 为空: ${TOKEN_RESP}"
  exit 1
fi
echo "token 获取成功, userID=${USER_ID}"

echo "==> 3) 调用 /user/get_self_login_platforms"
RESP=$(curl -sS -X POST \
  -H "Content-Type: application/json" \
  -H "operationID: $(op_id)" \
  -H "token: ${TOKEN}" \
  -d '{}' \
  "${HOST}/user/get_self_login_platforms")

echo "原始响应: ${RESP}"
ERR_CODE=$(echo "${RESP}" | jq -r '.errCode // "null"')
if [[ "${ERR_CODE}" != "0" ]]; then
  echo "接口调用失败: ${RESP}"
  exit 1
fi

echo "==> 4) 校验响应结构"
DATA_TYPE=$(echo "${RESP}" | jq -r '(.data | type) // "null"')
if [[ "${DATA_TYPE}" != "array" ]]; then
  echo "返回 data 不是数组: ${RESP}"
  exit 1
fi

echo "结构校验通过（data 为数组）"
echo "返回 data:"
echo "${RESP}" | jq '.data'

echo "测试通过: get_self_login_platforms"
