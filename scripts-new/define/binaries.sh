#!/bin/bash

# Define an associative array to store the binaries and their counts.
# The count for openim-msgtransfer is set to 4, all others are set to 1.
declare -A binaries=(
  [openim-test]=2
  [openim-no-port]=2
  )

