#!/bin/bash

# Define an associative array to store the binaries and their counts.
# The count for openim-msgtransfer is set to 4, all others are set to 1.
declare -A binaries=(
  [openim-api]=1
  [openim-cmdutils]=1
  [openim-crontask]=1
  [openim-msggateway]=1
  [openim-msgtransfer]=4  # openim-msgtransfer count is 4
  [openim-push]=1
  [openim-rpc-auth]=1
  [openim-rpc-conversation]=1
  [openim-rpc-friend]=1
  [openim-rpc-group]=1
  [openim-rpc-msg]=1
  [openim-rpc-third]=1
  [openim-rpc-user]=1
)

