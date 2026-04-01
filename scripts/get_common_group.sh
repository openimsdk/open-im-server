#!/usr/bin/env bash
set -euo pipefail

# ====== 按需修改 ======
API_BASE="${API_BASE:-http://127.0.0.1:10002}"   # 你的 open-im-api 地址
SELF_USER_ID="${SELF_USER_ID:-5694418935}"           # 当前登录用户（拿 token 的用户）
#FRIEND_USER_ID="${FRIEND_USER_ID:-1971806090}"       # 要查询共同群的好友
FRIEND_USER_ID="${FRIEND_USER_ID:-1011009748}"       # 要查询共同群的好友
PLATFORM_ID="${PLATFORM_ID:-2}"                  # 1=iOS, 2=Android, 3=Windows...
ADMIN_USER_ID="${ADMIN_USER_ID:-imAdmin}"        # 管理员账号（用于签发用户 token）
ADMIN_SECRET="${ADMIN_SECRET:-openIM123}"                 # 配置中的 share.secret
DEBUG="${DEBUG:-1}"                              # DEBUG=1 打印请求/响应明细
# RecordNotFoundError（errCode=1004）常见于 get_user_token：
# 服务端会查用户是否存在（user RPC GetDesignateUsers）；若 SELF_USER_ID 未注册，
# 返回空列表后 rpcli.firstValue 会包装为 ErrRecordNotFound（errDlt: record not found）。
# 处理：先注册该用户，或 export SELF_USER_ID=已存在用户，或 export TOKEN=已有用户 token 跳过拉 token。
#
# HTTP 404 + 响应体 "404 page not found"（Gin）：当前连上的 API 进程路由表里没有该路径。
# 本仓库已注册 POST /group/get_common_groups_with_friend（见 internal/api/router.go）。
# 处理：用当前代码重新编译/替换镜像并重启 openim-api，或确认 API_BASE 指向的就是带该路由的实例（无错误路径前缀/反代截断）。
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
  USER_ERR_CODE="$(echo "${USER_RESP:-}" | jq -r '.errCode // empty')"
  if [[ "${USER_ERR_CODE}" == "1004" ]]; then
    echo ""
    echo "【排查】errCode 1004 (RecordNotFoundError)：当前请求的 userID 在用户库中不存在。"
    echo "  - 服务端路径：auth GetUserToken → user GetDesignateUsers → 未命中则空结果 → record not found"
    echo "  - 请先将 SELF_USER_ID=${SELF_USER_ID} 注册进系统，或改用已存在用户，或: export TOKEN='你的用户token'"
  else
    echo "提示：请确认 SELF_USER_ID 已注册、ADMIN_SECRET 与部署一致，或手动 export TOKEN 后重试。"
  fi
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
GROUP_BODY="$(mktemp)"
GROUP_HTTP_CODE="$(
  curl -sS -o "${GROUP_BODY}" -w "%{http_code}" -X POST "${API_BASE}/group/get_common_groups_with_friend" \
    -H 'Content-Type: application/json' \
    -H "token: ${TOKEN}" \
    -H "operationID: ${OP_ID}" \
    -d "${REQ_BODY}"
)"
GROUP_RESP="$(cat "${GROUP_BODY}")"
rm -f "${GROUP_BODY}"
debug_log "group HTTP status: ${GROUP_HTTP_CODE}"
debug_log "group raw resp: ${GROUP_RESP}"
print_json_safe "${GROUP_RESP}"
if [[ "${GROUP_HTTP_CODE}" == "404" ]] || [[ "${GROUP_RESP}" == "404 page not found" ]]; then
  echo ""
  echo "【排查】HTTP 404：Gin 未匹配到路由，通常表示当前运行的 openim-api 版本过旧，不含 get_common_groups_with_friend。"
  echo "  - 期望路径: POST ${API_BASE}/group/get_common_groups_with_friend"
  echo "  - 请用本仓库代码重新构建并重启 API，或核对 API_BASE / 网关是否多删、少拼了路径前缀。"
  exit 1
fi
if [[ "${GROUP_HTTP_CODE}" != "200" ]]; then
  echo ""
  echo "【提示】HTTP 状态码: ${GROUP_HTTP_CODE}（非 200），请结合响应体与网关/鉴权配置排查。"
  exit 1
fi
