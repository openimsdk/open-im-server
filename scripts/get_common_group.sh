#!/usr/bin/env bash
set -euo pipefail

# ====== 按需修改 ======
API_BASE="${API_BASE:-http://127.0.0.1:10002}"   # 你的 open-im-api 地址
SELF_USER_ID="${SELF_USER_ID:-3932647710}"           # 当前登录用户（拿 token 的用户）
#FRIEND_USER_ID="${FRIEND_USER_ID:-4391832441}"       # 要查询共同群的好友
FRIEND_USER_ID="${FRIEND_USER_ID:-9607566286}"       # 要查询共同群的好友
PLATFORM_ID="${PLATFORM_ID:-2}"                  # 1=iOS, 2=Android, 3=Windows...
ADMIN_USER_ID="${ADMIN_USER_ID:-imAdmin}"        # 管理员账号（用于签发用户 token）
ADMIN_SECRET="${ADMIN_SECRET:-openIM123}"                 # 配置中的 share.secret
DEBUG="${DEBUG:-0}"                              # DEBUG=1 打印请求/响应明细
# =====================

debug_log() {
  if [[ "${DEBUG}" == "1" ]]; then
    echo "[DEBUG] $*"
  fi
}

print_json_safe() {
  local raw="${1:-}"
  if echo "${raw}" | jq -e . >/dev/null 2>&1; then
    echo "${raw}" | jq .
  else
    echo "${raw}"
  fi
}

# 1) 先拿 user token（如果你已有 token，可跳过这一步，直接 export TOKEN=xxx）
if [[ -z "${TOKEN:-}" ]]; then
  if [[ -z "${ADMIN_SECRET}" ]]; then
    echo "缺少 ADMIN_SECRET，请先导出：export ADMIN_SECRET='你的share.secret'"
    exit 1
  fi

  echo "获取管理员 token: ${ADMIN_USER_ID}"
  OP_ID_ADMIN="op_admin_$(date +%s)"
  debug_log "POST ${API_BASE}/auth/get_admin_token"
  debug_log "operationID: ${OP_ID_ADMIN}"
  debug_log "admin req body: {\"userID\":\"${ADMIN_USER_ID}\",\"secret\":\"***\"}"
  ADMIN_RESP=$(
    curl -sS -X POST "${API_BASE}/auth/get_admin_token" \
      -H 'Content-Type: application/json' \
      -H "operationID: ${OP_ID_ADMIN}" \
      -d "$(cat <<JSON
{
  "userID": "${ADMIN_USER_ID}",
  "secret": "${ADMIN_SECRET}"
}
JSON
)"
  )
  debug_log "admin raw resp: ${ADMIN_RESP}"
  ADMIN_TOKEN="$(echo "${ADMIN_RESP}" | jq -r '.data.token // empty')"
  debug_log "admin token parsed: ${ADMIN_TOKEN:-<empty>}"
  if [[ -z "${ADMIN_TOKEN}" ]]; then
    echo "获取管理员 token 失败，响应如下："
    print_json_safe "${ADMIN_RESP}"
    exit 1
  fi

  echo "获取用户 token: ${SELF_USER_ID}"
  OP_ID_USER="op_user_$(date +%s)"
  debug_log "POST ${API_BASE}/auth/get_user_token"
  debug_log "operationID: ${OP_ID_USER}"
  debug_log "user req body: {\"userID\":\"${SELF_USER_ID}\",\"platformID\":${PLATFORM_ID}}"
  USER_RESP=$(
    curl -sS -X POST "${API_BASE}/auth/get_user_token" \
      -H 'Content-Type: application/json' \
      -H "operationID: ${OP_ID_USER}" \
      -H "token: ${ADMIN_TOKEN}" \
      -d "$(cat <<JSON
{
  "userID": "${SELF_USER_ID}",
  "platformID": ${PLATFORM_ID}
}
JSON
)"
  )
  debug_log "user raw resp: ${USER_RESP}"
  TOKEN="$(echo "${USER_RESP}" | jq -r '.data.token // empty')"
  debug_log "user token parsed: ${TOKEN:-<empty>}"
fi

if [[ -z "${TOKEN}" ]]; then
  echo "获取用户 token 失败，响应如下："
  print_json_safe "${USER_RESP:-}"
  echo "提示：请确认 SELF_USER_ID 用户已注册存在，或手动传入 TOKEN 后重试。"
  exit 1
fi

OP_ID="op_$(date +%s)"

# 2) 调共同群接口
echo "查询共同群: self=${SELF_USER_ID}, friend=${FRIEND_USER_ID}"
REQ_BODY="$(cat <<JSON
{
  "friendUserID": "${FRIEND_USER_ID}"
}
JSON
)"
debug_log "POST ${API_BASE}/group/get_common_groups_with_friend"
debug_log "operationID: ${OP_ID}"
debug_log "group req body: ${REQ_BODY}"
GROUP_RESP="$(curl -sS -X POST "${API_BASE}/group/get_common_groups_with_friend" \
  -H 'Content-Type: application/json' \
  -H "token: ${TOKEN}" \
  -H "operationID: ${OP_ID}" \
  -d "${REQ_BODY}")"
debug_log "group raw resp: ${GROUP_RESP}"
print_json_safe "${GROUP_RESP}"
