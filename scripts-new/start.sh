
#!/usr/bin/env bash


OPENIM_SCRIPTS=$(dirname "${BASH_SOURCE[0]}")/
source "$OPENIM_SCRIPTS/bricks.sh"


# Assuming 'openim::util::host_platform' is defined in one of the sourced scripts or elsewhere.
# If not, you'll need to define it to return the appropriate platform directory name.

# Main function to start binaries


result=$(start_tools)
ret_val=$?


if [ $ret_val -ne 0 ]; then
  echo "tools start failed, abort start"
  echo "$result"
  exit 1
fi

echo "all tools ok"
echo "$result"

kill_exist_binaries

result=$(check_binaries_stop)
ret_val=$?

if [ $ret_val -ne 0 ]; then
  echo "$result"
  echo "Some services running, abort start"
  exit 1
fi


# Call the main function
result=$(start_binaries)
openim::log::print_blue_two_line "$result"

$OPENIM_SCRIPTS/check.sh

