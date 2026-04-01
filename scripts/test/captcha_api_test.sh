#!/usr/bin/env bash
# ============================================================
# Captcha API 接口测试脚本
#
# 覆盖接口：
#   POST /captcha/generate  —— 生成滑块验证码
#   POST /captcha/verify    —— 验证滑块验证码
#
# 依赖：curl / jq
# 用法：
#   chmod +x captcha_api_test.sh
#   ./captcha_api_test.sh
#   ./captcha_api_test.sh --host http://127.0.0.1:10002
# ============================================================

set -euo pipefail

# ──────────────────────────────────────────────
# 可配置参数（可通过环境变量覆盖）
# ──────────────────────────────────────────────
HOST="${HOST:-http://127.0.0.1:10002}"
ADMIN_USER_ID="${ADMIN_USER_ID:-imAdmin}"
ADMIN_SECRET="${ADMIN_SECRET:-openIM123}"
PLATFORM_ID="${PLATFORM_ID:-1}"       # 1=iOS 2=Android 3=Windows ...

# 命令行参数解析
while [[ $# -gt 0 ]]; do
  case "$1" in
    --host) HOST="$2"; shift 2 ;;
    --admin-user-id) ADMIN_USER_ID="$2"; shift 2 ;;
    --admin-secret) ADMIN_SECRET="$2"; shift 2 ;;
    *) echo "未知参数: $1"; exit 1 ;;
  esac
done

# ──────────────────────────────────────────────
# 颜色输出
# ──────────────────────────────────────────────
RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; CYAN='\033[0;36m'; NC='\033[0m'

PASS=0; FAIL=0

pass()    { echo -e "${GREEN}  [PASS]${NC} $1"; PASS=$((PASS+1)); }
fail()    { echo -e "${RED}  [FAIL]${NC} $1"; FAIL=$((FAIL+1)); }
info()    { echo -e "${CYAN}  [INFO]${NC} $1"; }
section() { echo -e "\n${YELLOW}══ $1 ══${NC}"; }

# ──────────────────────────────────────────────
# 生成唯一 operationID（每次调用递增）
# ──────────────────────────────────────────────
_OP_SEQ=0
new_op_id() {
  (( _OP_SEQ++ ))
  echo "captcha-test-$$-${_OP_SEQ}"
}

# ──────────────────────────────────────────────
# 断言工具函数
# ──────────────────────────────────────────────
assert_err_code() {
  local resp="$1" expected="$2" desc="$3"
  local actual
  actual=$(echo "$resp" | jq -r '.errCode // "null"')
  if [[ "${actual}" == "${expected}" ]]; then
    pass "${desc} (errCode=${actual})"
  else
    fail "${desc} - expected errCode=${expected}, got errCode=${actual}"
    info "resp: ${resp}"
  fi
}

assert_not_empty() {
  local resp="$1" jq_path="$2" desc="$3"
  local val
  val=$(echo "$resp" | jq -r "$jq_path // empty")
  if [[ -n "${val}" && "${val}" != "null" ]]; then
    pass "${desc} (val=${val:0:40}...)"
  else
    fail "${desc} - '${jq_path}' is empty or null"
    info "resp: ${resp}"
  fi
}

assert_eq() {
  local resp="$1" jq_path="$2" expected="$3" desc="$4"
  local actual
  # 不使用 // empty：jq 的 // 运算符会把布尔 false 视为 false 并走替代分支
  actual=$(echo "$resp" | jq -r "$jq_path")
  if [[ "${actual}" == "${expected}" ]]; then
    pass "${desc} (val=${actual})"
  else
    fail "${desc} - expected=${expected}, got=${actual}"
    info "resp: ${resp}"
  fi
}

# errCode 非 0 即通过
assert_err_nonzero() {
  local resp="$1" desc="$2"
  local actual
  actual=$(echo "$resp" | jq -r '.errCode // "null"')
  if [[ "${actual}" != "0" && "${actual}" != "null" ]]; then
    pass "${desc} (errCode=${actual})"
  else
    fail "${desc} - expected errCode!=0, got errCode=${actual}"
    info "resp: ${resp}"
  fi
}

# ──────────────────────────────────────────────
# 前置：获取 Admin Token
# ──────────────────────────────────────────────
section "前置：获取 Admin Token"

TOKEN_RESP=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -H "operationID: $(new_op_id)" \
  -d "{\"secret\":\"${ADMIN_SECRET}\",\"platformID\":${PLATFORM_ID},\"userID\":\"${ADMIN_USER_ID}\"}" \
  "${HOST}/auth/get_admin_token")

info "Token 响应: $TOKEN_RESP"

ERR_CODE=$(echo "$TOKEN_RESP" | jq -r '.errCode // "null"')
if [[ "$ERR_CODE" != "0" ]]; then
  echo -e "${RED}[ERROR]${NC} 获取 Admin Token 失败 (errCode=$ERR_CODE)，中止测试"
  exit 1
fi

TOKEN=$(echo "$TOKEN_RESP" | jq -r '.data.token')
info "获取到 token: ${TOKEN:0:40}..."

# ──────────────────────────────────────────────
# 用例 1：生成验证码 —— 正常流程
# ──────────────────────────────────────────────
section "用例 1 / POST /captcha/generate —— 正常生成验证码"

GEN_RESP=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -H "token: ${TOKEN}" \
  -H "operationID: $(new_op_id)" \
  -d '{}' \
  "${HOST}/captcha/generate")

info "响应摘要: $(echo "${GEN_RESP}" | jq -c '{errCode,errMsg,data:{captchaID:.data.captchaID,expireAt:.data.expireAt}}')"

GEN_ERR=$(echo "${GEN_RESP}" | jq -r '.errCode // "null"')
GEN_MSG=$(echo "${GEN_RESP}" | jq -r '.errMsg // ""')

# 检测服务端是否因缺少背景图资源而报 500
if [[ "${GEN_ERR}" == "500" && "${GEN_MSG}" == *"background"* ]]; then
  fail "用例 1 跳过 - captcha 服务未配置背景图资源 (errMsg=${GEN_MSG})"
  info "修复方式：在 captcha.Start() 中通过 slide.NewBuilder().SetBackground(...).Make() 注入背景图"
  CAPTCHA_ID=""
  EXPIRE_AT=""
else
  assert_err_code  "${GEN_RESP}" "0" "生成验证码 errCode 应为 0"
  assert_not_empty "${GEN_RESP}" ".data.captchaID"   "captchaID 非空"
  assert_not_empty "${GEN_RESP}" ".data.masterImage" "masterImage(背景图 Base64) 非空"
  assert_not_empty "${GEN_RESP}" ".data.tileImage"   "tileImage(滑块图 Base64) 非空"
  assert_not_empty "${GEN_RESP}" ".data.expireAt"    "expireAt(过期 Unix 时间戳) 非空"
  CAPTCHA_ID=$(echo "${GEN_RESP}" | jq -r '.data.captchaID')
  EXPIRE_AT=$(echo  "${GEN_RESP}" | jq -r '.data.expireAt')
  info "captchaID = ${CAPTCHA_ID}"
  info "expireAt  = ${EXPIRE_AT}"
fi

# ──────────────────────────────────────────────
# 用例 2：生成验证码 —— 不携带 Token
# ──────────────────────────────────────────────
section "用例 2 / POST /captcha/generate —— 无 Token 应被鉴权中间件拦截"

NO_TOKEN_RESP=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -H "operationID: $(new_op_id)" \
  -d '{}' \
  "${HOST}/captcha/generate")

info "响应: $NO_TOKEN_RESP"
assert_err_nonzero "$NO_TOKEN_RESP" "无 Token 被鉴权中间件拦截"

# ──────────────────────────────────────────────
# 用例 3：验证验证码 —— 坐标错误（x=999, y=999）
# ──────────────────────────────────────────────
section "用例 3 / POST /captcha/verify —— 坐标错误，success 应为 false"

if [[ -z "${CAPTCHA_ID}" ]]; then
  fail "用例 3 跳过 - 依赖用例 1 生成的 captchaID，但用例 1 未成功"
else
  VERIFY_WRONG_RESP=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -H "token: ${TOKEN}" \
    -H "operationID: $(new_op_id)" \
    -d "{\"captchaID\":\"${CAPTCHA_ID}\",\"x\":999,\"y\":999}" \
    "${HOST}/captcha/verify")
  info "响应: ${VERIFY_WRONG_RESP}"
  assert_err_code "${VERIFY_WRONG_RESP}" "0"     "验证请求本身成功 errCode=0"
  assert_eq       "${VERIFY_WRONG_RESP}" ".data.success" "false" "坐标错误时 success=false"
fi

# ──────────────────────────────────────────────
# 用例 4：验证验证码 —— 重复使用同一 captchaID
#         用例 3 已消耗该 ID（verify_time 已被 FindOneAndUpdate 写入），
#         再次调用服务端 filter 匹配不到记录，应返回错误
# ──────────────────────────────────────────────
section "用例 4 / POST /captcha/verify —— 重复使用同一 captchaID（幂等），应返回错误"

if [[ -z "${CAPTCHA_ID}" ]]; then
  fail "用例 4 跳过 - 依赖用例 1 生成的 captchaID，但用例 1 未成功"
else
  VERIFY_REUSE_RESP=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -H "token: ${TOKEN}" \
    -H "operationID: $(new_op_id)" \
    -d "{\"captchaID\":\"${CAPTCHA_ID}\",\"x\":0,\"y\":0}" \
    "${HOST}/captcha/verify")
  info "响应: ${VERIFY_REUSE_RESP}"
  assert_err_nonzero "${VERIFY_REUSE_RESP}" "重复使用 captchaID 被拒绝"
fi

# ──────────────────────────────────────────────
# 用例 5：验证验证码 —— captchaID 不存在
# ──────────────────────────────────────────────
section "用例 5 / POST /captcha/verify —— captchaID 不存在，应返回错误"

VERIFY_NOTFOUND_RESP=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -H "token: ${TOKEN}" \
  -H "operationID: $(new_op_id)" \
  -d '{"captchaID":"00000000-0000-0000-0000-000000000000","x":10,"y":10}' \
  "${HOST}/captcha/verify")

info "响应: $VERIFY_NOTFOUND_RESP"
assert_err_nonzero "$VERIFY_NOTFOUND_RESP" "captchaID 不存在时返回错误"

# ──────────────────────────────────────────────
# 用例 6：验证验证码 —— captchaID 为空字符串
# ──────────────────────────────────────────────
section "用例 6 / POST /captcha/verify —— captchaID 为空字符串，应返回错误"

VERIFY_EMPTY_RESP=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -H "token: ${TOKEN}" \
  -H "operationID: $(new_op_id)" \
  -d '{"captchaID":"","x":10,"y":10}' \
  "${HOST}/captcha/verify")

info "响应: $VERIFY_EMPTY_RESP"
assert_err_nonzero "$VERIFY_EMPTY_RESP" "captchaID 为空时返回错误"

# ──────────────────────────────────────────────
# 用例 7：验证验证码 —— 不携带 Token
# ──────────────────────────────────────────────
section "用例 7 / POST /captcha/verify —— 无 Token 应被鉴权中间件拦截"

VERIFY_NOTOKEN_RESP=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -H "operationID: $(new_op_id)" \
  -d "{\"captchaID\":\"${CAPTCHA_ID:-00000000-0000-0000-0000-000000000000}\",\"x\":10,\"y\":10}" \
  "${HOST}/captcha/verify")

info "响应: $VERIFY_NOTOKEN_RESP"
assert_err_nonzero "$VERIFY_NOTOKEN_RESP" "无 Token 被鉴权中间件拦截"

# ──────────────────────────────────────────────
# 用例 8：完整正向链路 —— 新生成 + 用偏差坐标验证
#         服务端不回传正确坐标，用 (0,0) 验证 success=false
#         正确坐标可从 MongoDB 查询：
#           db.captcha.findOne({captcha_id: "<ID>"}, {x:1,y:1})
# ──────────────────────────────────────────────
section "用例 8 / 完整正向链路 —— 新生成验证码 → 坐标偏差验证"

GEN_RESP2=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -H "token: ${TOKEN}" \
  -H "operationID: $(new_op_id)" \
  -d '{}' \
  "${HOST}/captcha/generate")

GEN_ERR2=$(echo "${GEN_RESP2}" | jq -r '.errCode // "null"')
GEN_MSG2=$(echo "${GEN_RESP2}" | jq -r '.errMsg  // ""')

if [[ "${GEN_ERR2}" == "500" && "${GEN_MSG2}" == *"background"* ]]; then
  fail "用例 8 跳过 - captcha 服务未配置背景图资源 (errMsg=${GEN_MSG2})"
else
  CAPTCHA_ID2=$(echo "${GEN_RESP2}" | jq -r '.data.captchaID')
  EXPIRE_AT2=$(echo  "${GEN_RESP2}" | jq -r '.data.expireAt')
  MASTER_LEN=$(echo  "${GEN_RESP2}" | jq -r '.data.masterImage | length')
  TILE_LEN=$(echo    "${GEN_RESP2}" | jq -r '.data.tileImage   | length')

  assert_err_code  "${GEN_RESP2}" "0" "新一轮生成验证码成功"
  assert_not_empty "${GEN_RESP2}" ".data.captchaID" "captchaID2 非空"

  info "captchaID2       = ${CAPTCHA_ID2}"
  info "expireAt         = ${EXPIRE_AT2}"
  info "masterImage 长度 = ${MASTER_LEN} chars(Base64)"
  info "tileImage 长度   = ${TILE_LEN} chars(Base64)"
  info "查询真实坐标: db.captcha.findOne({captcha_id:\"${CAPTCHA_ID2}\"},{x:1,y:1})"

  VERIFY_LINK_RESP=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -H "token: ${TOKEN}" \
    -H "operationID: $(new_op_id)" \
    -d "{\"captchaID\":\"${CAPTCHA_ID2}\",\"x\":0,\"y\":0}" \
    "${HOST}/captcha/verify")

  assert_err_code "${VERIFY_LINK_RESP}" "0"     "验证接口响应正常 errCode=0"
  assert_eq       "${VERIFY_LINK_RESP}" ".data.success" "false" "偏差坐标(0,0) success=false"
fi

# ──────────────────────────────────────────────
# 汇总
# ──────────────────────────────────────────────
echo ""
echo -e "══════════════════════════════════════════"
echo -e " 测试汇总：${GREEN}PASS=${PASS}${NC}  ${RED}FAIL=${FAIL}${NC}"
echo -e "══════════════════════════════════════════"

[[ $FAIL -eq 0 ]]
