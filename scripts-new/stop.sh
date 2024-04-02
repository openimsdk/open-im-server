
#!/usr/bin/env bash






OPENIM_SCRIPTS=$(dirname "${BASH_SOURCE[0]}")/
source "$OPENIM_SCRIPTS/bricks.sh"


kill_exist_binaries

result=$(check_binaries_stop)
ret_val=$?

if [ $ret_val -ne 0 ]; then
  echo "$result"
  echo "no stop..."
  exit 1
fi
