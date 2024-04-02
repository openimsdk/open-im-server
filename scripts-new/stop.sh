
#!/usr/bin/env bash






OPENIM_SCRIPTS=$(dirname "${BASH_SOURCE[0]}")/
source "$OPENIM_SCRIPTS/bricks.sh"


kill_exist_binaries

result=$(check_binaries_stop)
ret_val=$?
if [ $ret_val -ne 0 ]; then
  openim::log::print_red "Some services have not been stopped, details are as follows:"
  openim::log::print_red_no_time_stamp "$result"
  exit 1
fi

openim::log::print_green "All services have been stopped"