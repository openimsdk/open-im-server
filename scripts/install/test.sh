#!/usr/bin/env bash

# Copyright © 2023 OpenIM. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# OpenIM RPC Service Test Control Script
# 
# This control script is designed to conduct various tests on the OpenIM RPC services.
# It includes functions to perform smoke tests, API tests, and comprehensive service tests.
# The script is intended to be used in a Linux environment with appropriate permissions and
# environmental variables set.
# 
# It provides robust error handling and logging to facilitate debugging and service monitoring.
# Functions within the script can be called directly or passed as arguments to perform
# systematic testing, ensuring the integrity of the RPC services.
# 
# Test Functions:
# - openim::test::smoke: Runs basic tests to ensure the fundamental functionality of the service.
# - openim::test::api: Executes a series of API tests covering authentication, user, friend, 
#   group, and message functionalities.
# - openim::test::test: Performs a complete test suite, invoking utility checks and all defined
#   test cases, and reports on their success.
#

# The root of the build/dist directory
IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..
[[ -z ${COMMON_SOURCED} ]] && source ${IAM_ROOT}/scripts/install/common.sh

# API Server API Address:Port
INSECURE_OPENIMAPI=${IAM_APISERVER_HOST}:${API_OPENIM_PORT}
INSECURE_OPENIMAUTO=${OPENIM_RPC_AUTH_HOST}:${OPENIM_AUTH_PORT}

CCURL="curl -f -s -XPOST" # Create
UCURL="curl -f -s -XPUT" # Update
RCURL="curl -f -s -XGET" # Retrieve
DCURL="curl -f -s -XDELETE" # Delete

openim::test::check_error() {
  local response=$1
  local err_code=$(echo "$response" | jq '.errCode')
  openim::log::status "Response from user registration: $response"
  if [[ "$err_code" != "0" ]]; then
    openim::log::error_exit "Error occurred: $response, You can read the error code in the API documentation https://docs.openim.io/restapi/errcode"
  else
    openim::log::success "Operation was successful."
  fi
}

# The `openim::test::auth` function serves as a test suite for authentication-related operations.
function openim::test::auth() {
  # 1. Retrieve and set the authentication token.
  openim::test::get_token
  
  # 2. Force logout the test user from a specific platform.
  openim::test::force_logout
  
  # Log the completion of the auth test suite.
  openim::log::success "Auth test suite completed successfully."
}

#################################### Auth Module ####################################

# Define a function to get a token (Admin Token)
openim::test::get_token() {
  token_response=$(${CCURL} "${OperationID}" "${Header}" http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/auth/user_token \
      -d'{"secret": "'"$SECRET"'","platformID": 1,"userID": "openIM123456"}')
  token=$(echo $token_response | grep -Po 'token[" :]+\K[^"]+')
  echo "$token"
}

Header="-HContent-Type: application/json"
OperationID="-HoperationID: 1646445464564"
Token="-Htoken: $(openim::test::get_token)"

# Forces a user to log out from the specified platform by user ID.
openim::test::force_logout() {
  local request_body=$(cat <<EOF
{
  "platformID": 2,
  "userID": "4950983283"
}
EOF
  )
  echo "Requesting force logout for user: $request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/auth/force_logout" -d "${request_body}")

  openim::test::check_error "$response"
}


#################################### User Module ####################################

# Registers a new user with provided user ID, nickname, and face URL using the API.
openim::test::user_register() {
  # Assign the parameters to local variables, with defaults if not provided
  local user_id="${1:-${TEST_USER_ID}}"
  local nickname="${2:-cubxxw}"
  local face_url="${3:-https://github.com/cubxxw}"

  # Create the request body using the provided or default values
  local request_body=$(cat <<EOF
{
  "secret": "${SECRET}",
  "users": [
    {
      "userID": "${user_id}",
      "nickname": "${nickname}",
      "faceURL": "${face_url}"
    }
  ]
}
EOF
)

  echo "Request body for user registration: $request_body"

  # Send the registration request
  local user_register_response=$(${CCURL} "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/user/user_register" -d "${request_body}")

  # Check for errors in the response
  openim::test::check_error "$user_register_response"
}

# Checks if the provided list of user IDs exist in the system.
openim::test::check_user_account() {
  local request_body=$(cat <<EOF
{
  "checkUserIDs": [
    "${1}",
    "${MANAGER_USERID_1}",
    "${MANAGER_USERID_2}",
    "${MANAGER_USERID_3}"
  ]
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/user/account_check" -d "${request_body}")

  openim::test::check_error "$response"
}

# Retrieves a list of users with pagination, limited to a specific number per page.
openim::test::get_users() {
  local request_body=$(cat <<EOF
{
  "pagination": {
    "pageNumber": 1,
    "showNumber": 100
  }
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/user/get_users" -d "${request_body}")

  openim::test::check_error "$response"
}

# Obtains detailed information for a list of user IDs.
openim::test::get_users_info() {
  local request_body=$(cat <<EOF
{
  "userIDs": [
    "${1}",
    "${MANAGER_USERID_1}"
  ]
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/user/get_users_info" -d "${request_body}")

  openim::test::check_error "$response"
}

# Retrieves the online status for a list of user IDs.
openim::test::get_users_online_status() {
  local request_body=$(cat <<EOF
{
  "userIDs": [
    "${TEST_USER_ID}",
    "${MANAGER_USERID_1}",
    "${MANAGER_USERID_2}",
    "${MANAGER_USERID_3}"
  ]
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/user/get_users_online_status" -d "${request_body}")

  openim::test::check_error "$response"
}

# Updates the information for a user, such as nickname and face URL.
openim::test::update_user_info() {
  local request_body=$(cat <<EOF
{
  "userInfo": {
    "userID": "${TEST_USER_ID}",
    "nickname": "openimbot",
    "faceURL": "https://github.com/openimbot"
  }
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/user/update_user_info" -d "${request_body}")

  openim::test::check_error "$response"
}

# Gets the online status for users that a particular user has subscribed to.
openim::test::get_subscribe_users_status() {
  local request_body=$(cat <<EOF
{
  "userID": "${TEST_USER_ID}"
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/user/get_subscribe_users_status" -d "${request_body}")

  openim::test::check_error "$response"
}

# Subscribes to the online status of a list of users for a particular user ID.
openim::test::subscribe_users_status() {
  local request_body=$(cat <<EOF
{
  "userID": "9168684795",
  "userIDs": [
    "7475779354",
    "6317136453",
    "8450535746"
  ],
  "genre": 1
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/user/subscribe_users_status" -d "${request_body}")

  openim::test::check_error "$response"
}

# Sets the global message receiving option for a user, determining their messaging preferences.
openim::test::set_global_msg_recv_opt() {
  local request_body=$(cat <<EOF
{
  "userID": "${TEST_USER_ID}",
  "globalRecvMsgOpt": 0
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/user/set_global_msg_recv_opt" -d "${request_body}")

  openim::test::check_error "$response"
}

# [openim::test::user function description]
# The `openim::test::user` function serves as a test suite for user-related operations. 
# It sequentially invokes all user-related test functions to ensure the API's user operations are functioning correctly.
function openim::test::user() {
  # 1. Register a test user.
  local USER_ID=$RANDOM
  local TEST_USER_ID=$RANDOM
  openim::test::user_register "${USER_ID}" "user01" "new_face_url"
  openim::test::user_register "${TEST_USER_ID}" "user01" "new_face_url"
  # 2. Check if the test user's account exists.
  openim::test::check_user_account "${TEST_USER_ID}"
  
  # 3. Retrieve a list of users.
  openim::test::get_users
  
  # 4. Get detailed information for the test user.
  openim::test::get_users_info "${TEST_USER_ID}"
  
  # 5. Check the online status of the test user.
  openim::test::get_users_online_status
  
  # 6. Update the test user's information.
  openim::test::update_user_info
  
  # 7. Get the status of users subscribed by the test user.
  openim::test::get_subscribe_users_status
  
  # 8. Subscribe the test user to a set of user statuses.
  openim::test::subscribe_users_status
  
  # 9. Set the message receiving option for the test user.
  openim::test::set_global_msg_recv_opt
  
  # Log the completion of the user test suite.
  openim::log::success "User test suite completed successfully."
}

#################################### Friend Module ####################################

# Checks if two users are friends.
openim::test::is_friend() {
  local request_body=$(cat <<EOF
{
  "userID1": "${1}",
  "userID2": "${2}"
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/friend/is_friend" -d "${request_body}")

  openim::test::check_error "$response"
}

# Deletes a friend for a user.
openim::test::delete_friend() {
  local request_body=$(cat <<EOF
{
  "ownerUserID":"${1}",
  "friendUserID":"${2}"
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/friend/delete_friend" -d "${request_body}")

  openim::test::check_error "$response"
}

# Gets the friend application list for a user.
openim::test::get_friend_apply_list() {
  local request_body=$(cat <<EOF
{
  "userID": "${MANAGER_USERID_1}",
  "pagination": {
    "pageNumber": 1,
    "showNumber": 100
  }
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/friend/get_friend_apply_list" -d "${request_body}")

  openim::test::check_error "$response"
}

# Gets the friend list for a user.
openim::test::get_friend_list() {
  local request_body=$(cat <<EOF
{
  "userID": "${1}",
  "pagination": {
    "pageNumber": 1,
    "showNumber": 100
  }
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/friend/get_friend_list" -d "${request_body}")

  openim::test::check_error "$response"
}

# Sets a remark for a friend.
openim::test::set_friend_remark() {
  local request_body=$(cat <<EOF
{
  "ownerUserID": "${1}",
  "friendUserID": "${2}",
  "remark": "remark"
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/friend/set_friend_remark" -d "${request_body}")

  openim::test::check_error "$response"
}

# Adds a friend request.
openim::test::add_friend() {
  local request_body=$(cat <<EOF
{
  "fromUserID": "${1}",
  "toUserID": "${2}",
  "reqMsg": "hello!",
  "ex": ""
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/friend/add_friend" -d "${request_body}")

  openim::test::check_error "$response"
}

# Imports friends for a user.
openim::test::import_friend() {
  local friend_ids=$(printf ', "%s"' "${@:2}")
  friend_ids=${friend_ids:2}
  local request_body=$(cat <<EOF
{
  "ownerUserID": "${1}",
  "friendUserIDs": [${friend_ids}]
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/friend/import_friend" -d "${request_body}")

  openim::test::check_error "$response"
}


# Responds to a friend request.
openim::test::add_friend_response() {
  local request_body=$(cat <<EOF
{
  "fromUserID": "${1}",
  "toUserID": "${2}",
  "handleResult": 1,
  "handleMsg": "agree"
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/friend/add_friend_response" -d "${request_body}")

  openim::test::check_error "$response"
}

# Retrieves the friend application list that the user has applied for.
openim::test::get_self_friend_apply_list() {
  local request_body=$(cat <<EOF
{
  "userID": "${1}",
  "pagination": {
    "pageNumber": ${2},
    "showNumber": ${3}
  }
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/friend/get_self_friend_apply_list" -d "${request_body}")

  openim::test::check_error "$response"
}

# Adds a user to the blacklist.
openim::test::add_black() {
  local request_body=$(cat <<EOF
{
  "ownerUserID": "${1}",
  "blackUserID": "${2}"
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/friend/add_black" -d "${request_body}")

  openim::test::check_error "$response"
}

# Removes a user from the blacklist.
openim::test::remove_black() {
  local request_body=$(cat <<EOF
{
  "ownerUserID": "${1}",
  "blackUserID": "${2}"
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/friend/remove_black" -d "${request_body}")

  openim::test::check_error "$response"
}

# Retrieves the blacklist for a user.
openim::test::get_black_list() {
  local request_body=$(cat <<EOF
{
  "userID": "${1}",
  "pagination": {
    "pageNumber": 1,
    "showNumber": 100
  }
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/friend/get_black_list" -d "${request_body}")

  openim::test::check_error "$response"
}

# [openim::test::friend function description]
# The `openim::test::friend` function serves as a test suite for friend-related operations.
# It sequentially invokes all friend-related test functions to ensure the API's friend operations are functioning correctly.
function openim::test::friend() {
  local FRIEND_USER_ID=$RANDOM
  local BLACK_USER_ID=$RANDOM
  local TEST_USER_ID=$RANDOM
  # Assumes that TEST_USER_ID, FRIEND_USER_ID, and other necessary IDs are set as environment variables before running this suite.
  # 0. Register a friend user.
  openim::test::user_register "${TEST_USER_ID}" "user01" "new_face_url"
  openim::test::user_register "${FRIEND_USER_ID}" "frient01" "new_face_url"
  openim::test::user_register "${BLACK_USER_ID}" "frient02" "new_face_url"
  
  # 1. Check if two users are friends.
  openim::test::is_friend "${TEST_USER_ID}" "${FRIEND_USER_ID}"

  # 2. Send a friend request from one user to another.
  openim::test::add_friend "${TEST_USER_ID}" "${FRIEND_USER_ID}"

  # 3. Respond to a friend request.
  # TODO：
#   openim::test::add_friend_response "${FRIEND_USER_ID}" "${TEST_USER_ID}"

  # 4. Retrieve the friend list of the test user.
  openim::test::get_friend_list "${TEST_USER_ID}"

  # 5. Set a remark for a friend.
  # TODO：
#   openim::test::set_friend_remark "${TEST_USER_ID}" "${FRIEND_USER_ID}"

  # 6. Retrieve the friend application list for the test user.
  openim::test::get_friend_apply_list "${TEST_USER_ID}" 1 100

  # 7. Retrieve the friend application list that the user has applied for.
  openim::test::get_self_friend_apply_list "${TEST_USER_ID}" 1 100

  # 8. Delete a friend.
  # TODO：
#   openim::test::delete_friend "${TEST_USER_ID}" "${FRIEND_USER_ID}"

  # 9. Add a user to the blacklist.
  openim::test::add_black "${TEST_USER_ID}" "${BLACK_USER_ID}"

  # 10. Remove a user from the blacklist.
  openim::test::remove_black "${TEST_USER_ID}" "${BLACK_USER_ID}"

  # 11. Retrieve the blacklist for the test user.
  openim::test::get_black_list "${TEST_USER_ID}"

  # 12. Import friends for the user (Optional).
  # TODO：
#   openim::test::import_friend "${TEST_USER_ID}" "11111114" "11111115"

  # Log the completion of the friend test suite.
  openim::log::success "Friend test suite completed successfully."
}


#################################### Group Module ####################################

# Creates a new group.
openim::test::create_group() {
  local request_body=$(cat <<EOF
{
  "memberUserIDs": [
    "${1}"
  ],
  "adminUserIDs": [
    "${2}"
  ],
  "ownerUserID": "${3}",
  "groupInfo": {
    "groupID": "${4}",
    "groupName": "test-group",
    "notification": "notification",
    "introduction": "introduction",
    "faceURL": "faceURL url",
    "ex": "ex",
    "groupType": 2,
    "needVerification": 0,
    "lookMemberInfo": 0,
    "applyMemberFriend": 0
  }
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/group/create_group" -d "${request_body}")

  openim::test::check_error "$response"
}

# Invites a user to join a group.
openim::test::invite_user_to_group() {
  local request_body=$(cat <<EOF
{
  "groupID": "${1}",
  "invitedUserIDs": [
    "${2}",
    "${3}"
  ],
  "reason": "your reason"
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/group/invite_user_to_group" -d "${request_body}")

  openim::test::check_error "$response"
}

# Transfers the ownership of a group to another user.
openim::test::transfer_group() {
  local request_body=$(cat <<EOF
{
  "groupID":"${1}",
  "oldOwnerUserID":"${2}",
  "newOwnerUserID": "${3}"
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/group/transfer_group" -d "${request_body}")

  openim::test::check_error "$response"
}

# Retrieves information about multiple groups.
openim::test::get_groups_info() {
  local request_body=$(cat <<EOF
{
  "groupIDs": ["${1}", "${2}"]
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/group/get_groups_info" -d "${request_body}")

  openim::test::check_error "$response"
}

# Removes a user from a group.
openim::test::kick_group() {
  local request_body=$(cat <<EOF
{
  "groupID": "${1}",
  "kickedUserIDs": [
    "${2}"
  ],
  "reason": "Bye!"
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/group/kick_group" -d "${request_body}")

  openim::test::check_error "$response"
}

# Retrieves information about group members.
openim::test::get_group_members_info() {
  local request_body=$(cat <<EOF
{
  "groupID": "${1}",
  "userIDs": ["${2}"]
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/group/get_group_members_info" -d "${request_body}")

  openim::test::check_error "$response"
}

# Retrieves a list of group members.
openim::test::get_group_member_list() {
  local request_body=$(cat <<EOF
{
  "groupID": "${1}",
  "pagination": {
    "pageNumber": ${2},
    "showNumber": ${3}
  }
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/group/get_group_member_list" -d "${request_body}")

  openim::test::check_error "$response"
}

# Retrieves a list of groups that a user has joined.
openim::test::get_joined_group_list() {
  local request_body=$(cat <<EOF
{
  "fromUserID": "${1}",
  "pagination": {
    "showNumber": ${2},
    "pageNumber": ${3}
  }
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/group/get_joined_group_list" -d "${request_body}")

  openim::test::check_error "$response"
}


# Sets group member information.
openim::test::set_group_member_info() {
  local request_body=$(cat <<EOF
{
  "members": [
    { 
      "groupID": "${1}",
      "userID": "${2}",
      "nickName": "${3}",
      "faceURL": "${4}",
      "roleLevel": ${5},
      "ex":"${6}"
    }
  ]
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/group/set_group_member_info" -d "${request_body}")

  openim::test::check_error "$response"
}

# Mutes a group.
openim::test::mute_group() {
  local request_body=$(cat <<EOF
{
  "groupID": "${1}"
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/group/mute_group" -d "${request_body}")

  openim::test::check_error "$response"
}

# Cancels the muting of a group.
openim::test::cancel_mute_group() {
  local request_body=$(cat <<EOF
{
  "groupID": "${1}"
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/group/cancel_mute_group" -d "${request_body}")

  openim::test::check_error "$response"
}

# Dismisses a group.
openim::test::dismiss_group() {
  local request_body=$(cat <<EOF
{
  "groupID":"${1}"
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/group/dismiss_group" -d "${request_body}")

  openim::test::check_error "$response"
}

# Cancels muting a group member.
openim::test::cancel_mute_group_member() {
  local request_body=$(cat <<EOF
{
  "groupID": "${1}",
  "userID": "${2}"
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/group/cancel_mute_group_member" -d "${request_body}")

  openim::test::check_error "$response"
}

# Allows a user to join a group.
openim::test::join_group() {
  local request_body=$(cat <<EOF
{
 "groupID": "${1}",
 "reqMessage": "req msg join group",
 "joinSource": 0,
 "inviterUserID": "${2}"
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/group/join_group" -d "${request_body}")

  openim::test::check_error "$response"
}

# Sets group information.
openim::test::set_group_info() {
  local request_body=$(cat <<EOF
{
  "groupInfoForSet": {
    "groupID": "${1}",
    "groupName": "new-name",
    "notification": "new notification",
    "introduction": "new introduction",
    "faceURL": "www.newfaceURL.com",
    "ex": "new ex",
    "needVerification": 1,
    "lookMemberInfo": 1,
    "applyMemberFriend": 1
  }
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/group/set_group_info" -d "${request_body}")

  openim::test::check_error "$response"
}


# Allows a user to quit a group.
openim::test::quit_group() {
  local request_body=$(cat <<EOF
{
  "groupID": "${1}",
  "userID": "${2}"
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/group/quit_group" -d "${request_body}")

  openim::test::check_error "$response"
}

# Retrieves the list of group applications received by the user.
openim::test::get_recv_group_applicationList() {
  local request_body=$(cat <<EOF
{
  "fromUserID": "${1}",
  "pagination": {
    "pageNumber": ${2},
    "showNumber": ${3}
  }
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/group/get_recv_group_applicationList" -d "${request_body}")

  openim::test::check_error "$response"
}

# Responds to a group application.
openim::test::group_application_response() {
  local request_body=$(cat <<EOF
{
  "groupID": "${1}",
  "fromUserID": "${2}",
  "handledMsg": "",
  "handleResult": 1
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/group/group_application_response" -d "${request_body}")

  openim::test::check_error "$response"
}

# Retrieves the list of group applications made by the user.
openim::test::get_user_req_group_applicationList() {
  local request_body=$(cat <<EOF
{
  "userID": "${1}",
  "pagination": {
    "pageNumber": ${2},
    "showNumber": ${3}
  }
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/group/get_user_req_group_applicationList" -d "${request_body}")

  openim::test::check_error "$response"
}

# Mutes a group member for a specified duration.
openim::test::mute_group_member() {
  local request_body=$(cat <<EOF
{
  "groupID": "${1}",
  "userID": "${2}",
  "mutedSeconds": ${3}
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/group/mute_group_member" -d "${request_body}")

  openim::test::check_error "$response"
}

# Retrieves a list of group applications from specific users.
openim::test::get_group_users_req_application_list() {
  local request_body=$(cat <<EOF
{
  "groupID": "${1}",
  "userIDs": [
    "${2}"
  ]
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/group/get_group_users_req_application_list" -d "${request_body}")

  openim::test::check_error "$response"
}

# [openim::test::group function description]
# The `openim::test::group` function serves as a test suite for group-related operations.
# It sequentially invokes all group-related test functions to ensure the API's group operations are functioning correctly.
function openim::test::group() {
  local USER_ID=$RANDOM
  local OTHER_USER1_ID=$RANDOM
  local OTHER_USER2_ID=$RANDOM
  local TEST_USER_ID=$RANDOM

  local GROUP_ID=$RANDOM
  local GROUP_ID2=$RANDOM
  # Assumes that TEST_GROUP_ID, USER_ID, and other necessary IDs are set as environment variables before running this suite.
  # 0. Register a friend user.
  openim::test::user_register "${USER_ID}" "group00" "new_face_url"
  openim::test::user_register "${OTHER_USER1_ID}" "group01" "new_face_url"
  openim::test::user_register "${OTHER_USER2_ID}" "group02" "new_face_url"
  
  # 0. Create a new group.
  openim::test::create_group "$OTHER_USER2_ID" "$OTHER_USER1_ID" "$USER_ID" "$GROUP_ID"
  
  # 1. Invite user to group.
  openim::test::invite_user_to_group "$GROUP_ID" "$MANAGER_USERID_1" "$MANAGER_USERID_2"
  
  # 2. Transfer group ownership.
  openim::test::transfer_group "$GROUP_ID" "$USER_ID" "$OTHER_USER1_ID"
  
  # 3. Get group information.
  openim::test::get_groups_info "$GROUP_ID" "$OTHER_USER1_ID"
  
  # 4. Kick a user from the group.
  openim::test::kick_group "$GROUP_ID" "$OTHER_USER2_ID"
  
  # 5. Get group members info.
  openim::test::get_group_members_info "$GROUP_ID" "$USER_ID"
  
  # 6. Get group member list.
  openim::test::get_group_member_list "$GROUP_ID" 1 100
  
  # 7. Get joined group list.
  openim::test::get_joined_group_list "$USER_ID" 10 1
  
  # 8. Set group member info.
  openim::test::set_group_member_info "$GROUP_ID" "$USER_ID" "New NickName" "New Face URL" 60 "Extra Info"
  
  # 9. Mute group.
  openim::test::mute_group "$GROUP_ID"
  
  # 10. Cancel mute group.
  openim::test::cancel_mute_group "$GROUP_ID"
  
  # 11. Dismiss group.
  openim::test::dismiss_group "$GROUP_ID"

  openim::test::create_group "$OTHER_USER2_ID" "$OTHER_USER1_ID" "$USER_ID" "$GROUP_ID2"
  
  # 12. Cancel mute group member.
  openim::test::cancel_mute_group_member "$GROUP_ID" "$USER_ID"
  
  # 13. Join group.
  # TODO:
#   openim::test::join_group "$GROUP_ID2" "$OTHER_USER2_ID"
  
  # 14. Set group info.
  openim::test::set_group_info "$GROUP_ID2"
  
  # 15. Quit group.
  openim::test::quit_group "$GROUP_ID2" "$OTHER_USER1_ID"
  
  # 16. Get received group application list.
  openim::test::get_recv_group_applicationList "$USER_ID" 1 100
  
  # 17. Group application response.
  # TODO:
#   openim::test::group_application_response "$GROUP_ID2" "$OTHER_USER2_ID"
  
  # 18. Get user requested group application list.
  openim::test::get_user_req_group_applicationList "$USER_ID" 1 100
  
  # 19. Mute group member.
  openim::test::mute_group_member "$GROUP_ID" "$OTHER_USER1_ID" 3600
  
  # 20. Get group users request application list.
  openim::test::get_group_users_req_application_list "$GROUP_ID" "$USER_ID"
  
  # Log the completion of the group test suite.
  openim::log::success "Group test suite completed successfully."
}

#################################### Register And Check Module ####################################

# Define a function to register a user
openim::register_user() {
  user_register_response=$(${CCURL} "${Header}" http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/user/user_register \
    -d'{
      "secret": "openIM123",
      "users": [{"userID": "11111112","nickname": "yourNickname","faceURL": "yourFaceURL"}]
    }')
  
  echo "$user_register_response"
}

# Define a function to check the account
openim::test::check_account() {
  local token=$1
  account_check_response=$(${CCURL} "${Header}" -H"operationID: 1646445464564" -H"token: ${token}" http://localhost:${API_OPENIM_PORT}/user/account_check \
        -d'{
          "checkUserIDs": ["11111111","11111112"]
        }')
  
  echo "$account_check_response"
}

# Define a function to register, get a token and check the account
openim::test::register_and_check() {
  # Register a user
  user_register_response=$(openim::register_user)
  
  if [[ "$user_register_response" == *"\"errCode\": 0"* ]]; then
    echo "User registration successful."
  
    # Get token
    token=$(openim::get_token)
    
    if [[ -n "$token" ]]; then
      echo "Token acquired: $token"
    
      # Check account
      account_check_response=$(openim::check_account "$token")
      
      if [[ "$account_check_response" == *"\"errCode\": 0"* ]]; then
        echo "Account check successful."
      else
        echo "Account check failed."
      fi
    else
      echo "Failed to acquire token."
    fi
  else
    echo "User registration failed."
  fi
}

#################################### Msg Module ####################################

# Sends a message.
openim::test::send_msg() {
  local sendID="${1}"
  local recvID="${2}"
  local groupID="${3}"

  local request_body=$(cat <<EOF
{
  "sendID": "${sendID}",
  "recvID": "${recvID}",
  "groupID": "${groupID}",
  "senderNickname": "openIMAdmin-Gordon",
  "senderFaceURL": "http://www.head.com",
  "senderPlatformID": 1,
  "content": {
    "content": "hello!!"
  },
  "contentType": 101,
  "sessionType": 1,
  "isOnlineOnly": false,
  "notOfflinePush": false,
  "sendTime": $(date +%s)000,
  "offlinePushInfo": {
    "title": "send message",
    "desc": "",
    "ex": "",
    "iOSPushSound": "default",
    "iOSBadgeCount": true
  }
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/msg/send_msg" -d "${request_body}")

  openim::test::check_error "$response"
}

# Revokes a message.
openim::test::revoke_msg() {
  local userID="${1}"
  local conversationID="${2}"
  local seq="${3}"

  local request_body=$(cat <<EOF
{
  "userID": "${userID}",
  "conversationID": "${conversationID}",
  "seq": ${seq}
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/msg/revoke_msg" -d "${request_body}")

  openim::test::check_error "$response"
}


# Clears all messages for a user.
openim::test::user_clear_all_msg() {
  local userID="${1}"

  local request_body=$(cat <<EOF
{
  "userID": "${userID}"
}
EOF
)
  echo "$request_body"

  local response=$(${CCURL} "${Token}" "${OperationID}" "${Header}" "http://${OPENIM_API_HOST}:${API_OPENIM_PORT}/msg/user_clear_all_msg" -d "${request_body}")

  openim::test::check_error "$response"
}

# [openim::test::msg function description]
# The `openim::test::msg` function serves as a test suite for message-related operations.
# It sequentially invokes all message-related test functions to ensure the API's message operations are functioning correctly.
function openim::test::msg()
{
  local SEND_USER_ID="${MANAGER_USERID_1}" # This should be the sender's userID
  local GROUP_ID="" # GroupID if it's a group message
  local USER_ID="$RANDOM"
  openim::test::user_register "${USER_ID}" "msg00" "new_face_url"
  local RECV_USER_ID="${USER_ID}" # Receiver's userID

  # 0. Send a message.
  openim::test::send_msg "${SEND_USER_ID}" "${RECV_USER_ID}" "${GROUP_ID}"
  
  # Assuming message sending was successful and returned a sequence number.
  local SEQ_NUMBER=1 # This should be the actual sequence number of the message sent.
  
  # 1. Revoke a message.
  # TODO：
  # openim::test::revoke_msg "${RECV_USER_ID}" "si_${SEND_USER_ID}_${RECV_USER_ID}" "${SEQ_NUMBER}"

  # 2. Clear all messages for a user.
  openim::test::user_clear_all_msg "${RECV_USER_ID}"

  # Log the completion of the message test suite.
  openim::log::success "Message test suite completed successfully."
}

#################################### Man Module ####################################

# TODO:

openim::test::man() {
  openim::log::info "TODO: openim test man"
}


#################################### Build Module ####################################

# Function: openim::test::smoke
# Purpose: Performs a series of basic tests to validate the core functionality of the system.
# These are preliminary checks to ensure that the most crucial operations like user registration
# and account checking are operational.
openim::test::smoke() {
  openim::register_user
  openim::test::check_account
  openim::test::register_and_check
}

# Function: openim::test::api
# Purpose: This function is a collection of API test calls that cover various aspects of the 
# service such as authentication, user operations, friend management, group interactions, and 
# message handling. It is used to verify the integrity and functionality of API endpoints.
openim::test::api() {
  openim::test::auth
  openim::test::user
  openim::test::friend
  openim::test::group
  openim::test::msg
}

# Function: openim::test::test
# Purpose: This is the comprehensive test function that invokes all individual test functions.
# It ensures that each component of the service is tested, utilizing utility functions for
# environment checking and completing with a success message if all tests pass.
openim::test::test() {
  openim::util::require-jq
  openim::test::smoke
  openim::test::man
  openim::test::api

  openim::log::info "$(echo -e '\033[32mcongratulations, all test passed!\033[0m')"
}

# Main execution logic: This conditional block checks if the script's arguments match any known
# test function patterns and, if so, evaluates the function call. This allows for specific test
# functions to be triggered based on the passed arguments.
if [[ "$*" =~ openim::test:: ]];then
  eval $*
fi