#!/usr/bin/env bash
set -euo pipefail

# 测试 open-im-server 手机号 is_snd 接口（/phone/get_sn_info、/phone/set_sn_info）
#
# 用法：
#   查询（无需 token，已加入 GinParseToken 白名单）：
#     ./scripts/phone_sn_api.sh get <phone>
#
#   写入（需要用户 token）：
#     OPENIM_TOKEN="<用户 token>" ./scripts/phone_sn_api.sh set <phone> <userID整数> <is_snd:0|1|true|false>
#
# 环境变量（可覆盖）：
#   OPENIM_API_ADDR  默认 http://127.0.0.1:10002
#   OPENIM_TOKEN     set 时必填（Header: token）
#   OPERATION_ID     默认自动生成

OPENIM_API_ADDR="${OPENIM_API_ADDR:-http://127.0.0.1:10002}"
OPENIM_TOKEN="${OPENIM_TOKEN:-}"
OPERATION_ID="${OPERATION_ID:-phone_sn_$(date +%s)_$RANDOM}"

ACTION="${1:-}"
PHONE="${2:-}"
USER_ID="${3:-}"
IS_SND_RAW="${4:-}"

die() {
  echo "ERROR: $*" >&2
  exit 1
}

usage() {
  cat <<'EOF'
用法：
  查询（无需 token）：
    ./scripts/phone_sn_api.sh get <phone>

  写入（需要 OPENIM_TOKEN）：
    OPENIM_TOKEN="<用户token>" ./scripts/phone_sn_api.sh set <phone> <userID整数> <is_snd:0|1|true|false>

环境变量：
  OPENIM_API_ADDR   默认 http://127.0.0.1:10002
  OPENIM_TOKEN      set 时必填
  OPERATION_ID      可选，默认自动生成
EOF
}

curl_post() {
  local path=$1
  local json_body=$2
  local with_token=${3:-0}
  local -a hdrs=(
    -H "Content-Type: application/json"
    -H "operationID: ${OPERATION_ID}"
  )
  if [[ "$with_token" == "1" ]]; then
    [[ -n "$OPENIM_TOKEN" ]] || die "set 接口需要环境变量 OPENIM_TOKEN（用户 token）"
    hdrs+=(-H "token: ${OPENIM_TOKEN}")
  fi
  curl -sS "${hdrs[@]}" -X POST "${OPENIM_API_ADDR}${path}" -d "${json_body}"
}

pretty_json() {
  if command -v jq >/dev/null 2>&1; then
    jq .
  else
    cat
  fi
}

case "$ACTION" in
  get)
    [[ -n "$PHONE" ]] || die "用法: $0 get <phone>"
    body=$(printf '{"phone":"%s"}' "$PHONE")
    echo "==> POST ${OPENIM_API_ADDR}/phone/get_sn_info"
    echo "    body: ${body}"
    resp=$(curl_post "/phone/get_sn_info" "$body" 0)
    echo "$resp" | pretty_json
    ;;
  set)
    [[ -n "$PHONE" && -n "$USER_ID" && -n "$IS_SND_RAW" ]] || die "用法: $0 set <phone> <userID> <is_snd:0|1|true|false>"
    case "$IS_SND_RAW" in
      1|true|True|TRUE) is_snd=true ;;
      0|false|False|FALSE) is_snd=false ;;
      *) die "is_snd 必须是 0、1、true 或 false" ;;
    esac
    # userID 必须为 JSON 数字
    if ! [[ "$USER_ID" =~ ^-?[0-9]+$ ]]; then
      die "userID 必须为整数"
    fi
    body=$(printf '{"phone":"%s","userID":%s,"is_snd":%s}' "$PHONE" "$USER_ID" "$is_snd")
    echo "==> POST ${OPENIM_API_ADDR}/phone/set_sn_info"
    echo "    body: ${body}"
    resp=$(curl_post "/phone/set_sn_info" "$body" 1)
    echo "$resp" | pretty_json
    ;;
  ""|-h|--help|help)
    usage
    exit 0
    ;;
  *)
    die "未知动作: $ACTION，使用 -h 查看帮助"
    ;;
esac
