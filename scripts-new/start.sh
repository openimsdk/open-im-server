
#!/usr/bin/env bash


OPENIM_SCRIPTS=$(dirname "${BASH_SOURCE[0]}")/
source "$OPENIM_SCRIPTS/bricks.sh"


# Assuming 'openim::util::host_platform' is defined in one of the sourced scripts or elsewhere.
# If not, you'll need to define it to return the appropriate platform directory name.

# Main function to start binaries

openim::log::print_blue "Starting tools"

result=$(start_tools)
ret_val=$?
if [ $ret_val -ne 0 ]; then
  openim::log::print_red "Some tools failed to start, details are as follows, abort start"
  openim::log::print_red_no_time_stamp "$result"
  exit 1
fi

echo "$result"


openim::log::print_green "All tools executed successfully"
openim::log::print_blue "Starting services"

kill_exist_binaries

result=$(check_binaries_stop)
ret_val=$?

if [ $ret_val -ne 0 ]; then
  openim::log::print_red "Some services running, details are as follows, abort start"
  openim::log::print_red_no_time_stamp "$result"
  exit 1
fi


# Call the main function
result=$(start_binaries)
openim::log::print_blue_two_line "$result"

$OPENIM_SCRIPTS/check.sh

