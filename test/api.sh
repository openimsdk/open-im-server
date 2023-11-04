#!/usr/bin/env bash

# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.


# The root of the build/dist directory
IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..
[[ -z ${COMMON_SOURCED} ]] && source ${IAM_ROOT}/scripts/install/common.sh

INSECURE_APISERVER=${IAM_APISERVER_HOST}:${IAM_APISERVER_INSECURE_BIND_PORT}
INSECURE_AUTHZSERVER=${IAM_AUTHZ_SERVER_HOST}:${IAM_AUTHZ_SERVER_INSECURE_BIND_PORT}

Header="-HContent-Type: application/json"
CCURL="curl -f -s -XPOST" # Create
UCURL="curl -f -s -XPUT" # Update
RCURL="curl -f -s -XGET" # Retrieve
DCURL="curl -f -s -XDELETE" # Delete

openim::test::login()
{
  ${CCURL} "${Header}" http://${INSECURE_APISERVER}/login \
    -d'{"username":"admin","password":"Admin@2021"}' | grep -Po 'token[" :]+\K[^"]+'
}

openim::test::user()
{
  token="-HAuthorization: Bearer $(openim::test::login)"

  # 1. 如果有 colin、mark、john 用户先清空
  ${DCURL} "${token}" http://${INSECURE_APISERVER}/v1/users/colin; echo
  ${DCURL} "${token}" http://${INSECURE_APISERVER}/v1/users/mark; echo
  ${DCURL} "${token}" http://${INSECURE_APISERVER}/v1/users/john; echo

  # 2. 创建 colin、mark、john 用户
  ${CCURL} "${Header}" "${token}" http://${INSECURE_APISERVER}/v1/users \
    -d'{"password":"User@2021","metadata":{"name":"colin"},"nickname":"colin","email":"colin@foxmail.com","phone":"1812884xxxx"}'; echo

  # 3. 列出所有用户
  ${RCURL} "${token}" "http://${INSECURE_APISERVER}/v1/users?offset=0&limit=10"; echo

  # 4. 获取 colin 用户的详细信息
  ${RCURL} "${token}" http://${INSECURE_APISERVER}/v1/users/colin; echo

  # 5. 修改 colin 用户
  ${UCURL} "${Header}" "${token}" http://${INSECURE_APISERVER}/v1/users/colin \
    -d'{"nickname":"colin","email":"colin_modified@foxmail.com","phone":"1812884xxxx"}'; echo

  # 6. 删除 colin 用户
  ${DCURL} "${token}" http://${INSECURE_APISERVER}/v1/users/colin; echo

  # 7. 批量删除用户
  ${DCURL} "${token}" "http://${INSECURE_APISERVER}/v1/users?name=mark&name=john"; echo
  openim::log::info "$(echo -e '\033[32mcongratulations, /v1/user test passed!\033[0m')"
}

openim::test::secret()
{
  token="-HAuthorization: Bearer $(openim::test::login)"

  # 1. 如果有 secret0 密钥先清空
  ${DCURL} "${token}" http://${INSECURE_APISERVER}/v1/secrets/secret0; echo

  # 2. 创建 secret0 密钥
  ${CCURL} "${Header}" "${token}" http://${INSECURE_APISERVER}/v1/secrets \
    -d'{"metadata":{"name":"secret0"},"expires":0,"description":"admin secret"}'; echo

  # 3. 列出所有密钥
  ${RCURL} "${token}" http://${INSECURE_APISERVER}/v1/secrets; echo

  # 4. 获取 secret0 密钥的详细信息
  ${RCURL} "${token}" http://${INSECURE_APISERVER}/v1/secrets/secret0; echo

  # 5. 修改 secret0 密钥
  ${UCURL} "${Header}" "${token}" http://${INSECURE_APISERVER}/v1/secrets/secret0 \
    -d'{"expires":0,"description":"admin secret(modified)"}'; echo

  # 6. 删除 secret0 密钥
  ${DCURL} "${token}" http://${INSECURE_APISERVER}/v1/secrets/secret0; echo
  openim::log::info "$(echo -e '\033[32mcongratulations, /v1/secret test passed!\033[0m')"
}

openim::test::policy()
{
  token="-HAuthorization: Bearer $(openim::test::login)"

  # 1. 如果有 policy0 策略先清空
  ${DCURL} "${token}" http://${INSECURE_APISERVER}/v1/policies/policy0; echo

  # 2. 创建 policy0 策略
  ${CCURL} "${Header}" "${token}" http://${INSECURE_APISERVER}/v1/policies \
    -d'{"metadata":{"name":"policy0"},"policy":{"description":"One policy to rule them all.","subjects":["users:<peter|ken>","users:maria","groups:admins"],"actions":["delete","<create|update>"],"effect":"allow","resources":["resources:articles:<.*>","resources:printer"],"conditions":{"remoteIPAddress":{"type":"CIDRCondition","options":{"cidr":"192.168.0.1/16"}}}}}'; echo

  # 3. 列出所有策略
  ${RCURL} "${token}" http://${INSECURE_APISERVER}/v1/policies; echo

  # 4. 获取 policy0 策略的详细信息
  ${RCURL} "${token}" http://${INSECURE_APISERVER}/v1/policies/policy0; echo

  # 5. 修改 policy0 策略
  ${UCURL} "${Header}" "${token}" http://${INSECURE_APISERVER}/v1/policies/policy0 \
    -d'{"policy":{"description":"One policy to rule them all(modified).","subjects":["users:<peter|ken>","users:maria","groups:admins"],"actions":["delete","<create|update>"],"effect":"allow","resources":["resources:articles:<.*>","resources:printer"],"conditions":{"remoteIPAddress":{"type":"CIDRCondition","options":{"cidr":"192.168.0.1/16"}}}}}'; echo

  # 6. 删除 policy0 策略
  ${DCURL} "${token}" http://${INSECURE_APISERVER}/v1/policies/policy0; echo
  openim::log::info "$(echo -e '\033[32mcongratulations, /v1/policy test passed!\033[0m')"
}

openim::test::apiserver()
{
  openim::test::user
  openim::test::secret
  openim::test::policy
  openim::log::info "$(echo -e '\033[32mcongratulations, openim-apiserver test passed!\033[0m')"
}

openim::test::authz()
{
  token="-HAuthorization: Bearer $(openim::test::login)"

  # 1. 如果有 authzpolicy 策略先清空
  ${DCURL} "${token}" http://${INSECURE_APISERVER}/v1/policies/authzpolicy; echo

  # 2. 创建 authzpolicy 策略
  ${CCURL} "${Header}" "${token}" http://${INSECURE_APISERVER}/v1/policies \
    -d'{"metadata":{"name":"authzpolicy"},"policy":{"description":"One policy to rule them all.","subjects":["users:<peter|ken>","users:maria","groups:admins"],"actions":["delete","<create|update>"],"effect":"allow","resources":["resources:articles:<.*>","resources:printer"],"conditions":{"remoteIPAddress":{"type":"CIDRCondition","options":{"cidr":"192.168.0.1/16"}}}}}'; echo

  # 3. 如果有 authzsecret 密钥先清空
  ${DCURL} "${token}" http://${INSECURE_APISERVER}/v1/secrets/authzsecret; echo

  # 4. 创建 authzsecret 密钥
  secret=$(${CCURL} "${Header}" "${token}" http://${INSECURE_APISERVER}/v1/secrets -d'{"metadata":{"name":"authzsecret"},"expires":0,"description":"admin secret"}')
  secretID=$(echo ${secret} | grep -Po 'secretID[" :]+\K[^"]+')
  secretKey=$(echo ${secret} | grep -Po 'secretKey[" :]+\K[^"]+')

  # 5. 生成 token
  token=$(iamctl jwt sign ${secretID} ${secretKey})

  # 6. 调用 /v1/authz 完成资源授权。
  # 注意这里要 sleep 3s 等待 openim-authz-server 将新建的密钥同步到其内存中
  echo "wait 3s to allow openim-authz-server to sync information into its memory ..."
  sleep 3
  ret=`$CCURL "${Header}" -H"Authorization: Bearer ${token}" http://${INSECURE_AUTHZSERVER}/v1/authz \
    -d'{"subject":"users:maria","action":"delete","resource":"resources:articles:ladon-introduction","context":{"remoteIPAddress":"192.168.0.5"}}' | grep -Po 'allowed[" :]+\K\w+'`

  if [ "$ret" != "true" ];then
    return 1
  fi

  openim::log::info "$(echo -e '\033[32mcongratulations, /v1/authz test passed!\033[0m')"
}

openim::test::authzserver()
{
  openim::test::authz
  openim::log::info "$(echo -e '\033[32mcongratulations, openim-authz-server test passed!\033[0m')"
}

openim::test::pump()
{
  ${RCURL} http://${IAM_PUMP_HOST}:7070/healthz | egrep -q 'status.*ok' || {
    openim::log::error "cannot access openim-pump healthz api, openim-pump maybe down"
      return 1
    }

  openim::test::real_pump_test

  openim::log::info "$(echo -e '\033[32mcongratulations, openim-pump test passed!\033[0m')"
}

# 使用真实的数据测试 openim-pump 是否正常工作
openim::test::real_pump_test()
{
  # 1. 创建访问 openim-authz-server 需要用到的密钥对
  iamctl secret create pumptest &>/dev/null

  # 2. 使用步骤 1 创建的密钥对生成 JWT Token
  authzAccessToken=`iamctl jwt sign njcho8gJQArsq7zr5v1YpG5NcvL0aeuZ38Ti if70HgRgp021iq5ex2l7pfy5XvgtZM3q` # iamctl jwt sign $secretID $secretKey

  # 3. 创建授权策略
  iamctl policy create pumptest '{"metadata":{"name":"policy0"},"policy":{"description":"One policy to rule them all.","subjects":["users:<peter|ken>","users:maria","groups:admins"],"actions":["delete","<create|update>"],"effect":"allow","resources":["resources:articles:<.*>","resources:printer"],"conditions":{"remoteIPAddress":{"type":"CIDRCondition","options":{"cidr":"192.168.0.1/16"}}}}}' &>/dev/null

  # 注意这里要 sleep 3s 等待 openim-authz-server 将新建的密钥和授权策略同步到其内存中
  echo "wait 3s to allow openim-authz-server to sync information into its memory ..."
  sleep 3

  # 4. 访问 /v1/authz 接口进行资源授权
  $CCURL "${Header}" -H"Authorization: Bearer ${token}" http://${INSECURE_AUTHZSERVER}/v1/authz \
    -d'{"subject":"users:maria","action":"delete","resource":"resources:articles:ladon-introduction","context":{"remoteIPAddress":"192.168.0.5"}}' &>/dev/null

  # 这里要 sleep 5s，等待 openim-pump 将 Redis 中的日志，分析并转存到 MongoDB 中
  echo "wait 10s to allow openim-pump analyze and dump authorization log into MongoDB ..."
  sleep 10

  # 5. 查看 MongoDB 中是否有经过解析后的授权日志。
  echo "db.iam_analytics.find()" | mongosh --quiet "${IAM_PUMP_MONGO_URL}" | grep -q "allow access" || {
    openim::log::error "cannot find analyzed authorization log in MongoDB"
      return 1
    }
}

openim::test::watcher()
{
  ${RCURL} http://${IAM_WATCHER_HOST}:5050/healthz | egrep -q 'status.*ok' || {
    openim::log::error "cannot access openim-watcher healthz api, openim-watcher maybe down"
      return 1
    }
  openim::log::info "$(echo -e '\033[32mcongratulations, openim-watcher test passed!\033[0m')"
}

openim::test::iamctl()
{
  iamctl user list | egrep -q admin || {
    openim::log::error "iamctl cannot list users from openim-apiserver"
      return 1
    }
  openim::log::info "$(echo -e '\033[32mcongratulations, iamctl test passed!\033[0m')"
}

openim::test::man()
{
  man openim-apiserver | grep -q 'OPENIM API Server' || {
    openim::log::error "openim man page not installed or may not installed properly"
      return 1
    }
  openim::log::info "$(echo -e '\033[32mcongratulations, man test passed!\033[0m')"
}

openim::test::smoke()
{
  openim::test::apiserver
  openim::test::authzserver
  openim::test::pump
  openim::test::watcher
  openim::test::iamctl
  openim::log::info "$(echo -e '\033[32mcongratulations, smoke test passed!\033[0m')"
}

openim::test::test()
{
  openim::test::smoke
  openim::test::man

  openim::log::info "$(echo -e '\033[32mcongratulations, all test passed!\033[0m')"
}

if [[ "$*" =~ openim::test:: ]];then
  eval $*
fi
