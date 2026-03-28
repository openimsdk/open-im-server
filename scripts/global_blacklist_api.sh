#!/usr/bin/env bash
set -euo pipefail

# 统一通过 API 新链路管理全局黑名单（按 nickname）
#
# 用法：
#   1) 添加
#      ./scripts/global_blacklist_api.sh add "alice,bob" [reason]
#
#   2) 删除
#      ./scripts/global_blacklist_api.sh remove "alice,bob"
#
#   3) 查询
#      ./scripts/global_blacklist_api.sh list [pageNumber] [showNumber]
#
# 环境变量（可覆盖）：
#   OPENIM_API_ADDR   默认: http://127.0.0.1:10002
#   ADMIN_TOKEN       管理员 token（如未提供则自动调用 /auth/get_admin_token 获取）
#   OPENIM_SECRET     获取管理员 token 所需 secret，默认: openIM123
#   ADMIN_USER_ID     获取管理员 token 所需 userID，默认: imAdmin

OPENIM_API_ADDR="${OPENIM_API_ADDR:-http://127.0.0.1:10002}"
ADMIN_TOKEN="${ADMIN_TOKEN:-}"
OPENIM_SECRET="${OPENIM_SECRET:-openIM123}"
ADMIN_USER_ID="${ADMIN_USER_ID:-imAdmin}"
OPERATION_ID="${OPERATION_ID:-gb_$(date +%s)_$RANDOM}"

ACTION="${1:-}"
NICKNAMES_RAW="${2:-}"
REASON="${3:-manual_by_api_script}"
PAGE_NUMBER="${2:-1}"
SHOW_NUMBER="${3:-20}"

die() {
  echo "ERROR: $*" >&2
  exit 1
}

trim() {
  local s="$1"
  s="${s#"${s%%[![:space:]]*}"}"
  s="${s%"${s##*[![:space:]]}"}"
  printf '%s' "$s"
}

nicknames_csv_to_json_array() {
  local csv="$1"
  local arr_json="["
  local first=1
  local item

  IFS=',' read -r -a _items <<< "$csv"
  for item in "${_items[@]}"; do
    item="$(trim "$item")"
    [[ -z "$item" ]] && continue
    if [[ $first -eq 1 ]]; then
      arr_json="${arr_json}\"${item}\""
      first=0
    else
      arr_json="${arr_json},\"${item}\""
    fi
  done
  arr_json="${arr_json}]"

  if [[ "$arr_json" == "[]" ]]; then
    die "nicknames 为空，请传入逗号分隔昵称，如 \"alice,bob\""
  fi
  printf '%s' "$arr_json"
}

get_admin_token() {
  local uid body resp token last_resp
  local -a candidates=("${ADMIN_USER_ID}" "openIM123456" "imAdmin")
  last_resp=""

  for uid in "${candidates[@]}"; do
    body="{\"secret\":\"${OPENIM_SECRET}\",\"userID\":\"${uid}\"}"
    resp="$(curl -sS -X POST "${OPENIM_API_ADDR}/auth/get_admin_token" \
      -H "Content-Type: application/json" \
      -H "operationID: ${OPERATION_ID}" \
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
  die "自动获取管理员 token 失败，请检查 OPENIM_API_ADDR/OPENIM_SECRET/ADMIN_USER_ID（当前: ${ADMIN_USER_ID}），或直接设置 ADMIN_TOKEN"
}

call_api() {
  local path="$1"
  local body="$2"
  local token="$3"

  curl -sS -X POST "${OPENIM_API_ADDR}${path}" \
    -H "Content-Type: application/json" \
    -H "operationID: ${OPERATION_ID}" \
    -H "token: ${token}" \
    -d "$body"
}

if [[ -z "$ACTION" ]]; then
  cat <<'EOF'
用法:
  添加: ./scripts/global_blacklist_api.sh add "alice,bob" [reason]
  删除: ./scripts/global_blacklist_api.sh remove "alice,bob"
  查询: ./scripts/global_blacklist_api.sh list [pageNumber] [showNumber]
EOF
  exit 1
fi

if [[ -z "$ADMIN_TOKEN" ]]; then
  echo "ADMIN_TOKEN 未设置，尝试自动获取管理员 token..."
  ADMIN_TOKEN="$(get_admin_token)"
fi

case "$ACTION" in
  add)
    [[ -z "$NICKNAMES_RAW" ]] && die "add 需要 nicknames 参数"
    NICKNAMES_JSON="$(nicknames_csv_to_json_array "$NICKNAMES_RAW")"
    BODY="{\"nicknames\":${NICKNAMES_JSON},\"reason\":\"${REASON}\"}"
    echo ">>> POST /user/add_global_blacklist"
    call_api "/user/add_global_blacklist" "$BODY" "$ADMIN_TOKEN"
    ;;

  remove)
    [[ -z "$NICKNAMES_RAW" ]] && die "remove 需要 nicknames 参数"
    NICKNAMES_JSON="$(nicknames_csv_to_json_array "$NICKNAMES_RAW")"
    BODY="{\"nicknames\":${NICKNAMES_JSON}}"
    echo ">>> POST /user/remove_global_blacklist"
    call_api "/user/remove_global_blacklist" "$BODY" "$ADMIN_TOKEN"
    ;;

  list)
    BODY="{\"pagination\":{\"pageNumber\":${PAGE_NUMBER},\"showNumber\":${SHOW_NUMBER}}}"
    echo ">>> POST /user/get_global_blacklist"
    call_api "/user/get_global_blacklist" "$BODY" "$ADMIN_TOKEN"
    ;;

  *)
    die "不支持的 action: ${ACTION}（仅支持 add/remove/list）"
    ;;
esac

echo
