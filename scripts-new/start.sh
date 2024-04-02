
#!/usr/bin/env bash


OPENIM_SCRIPTS=$(dirname "${BASH_SOURCE[0]}")/
source "$OPENIM_SCRIPTS/bricks.sh"


# Assuming 'openim::util::host_platform' is defined in one of the sourced scripts or elsewhere.
# If not, you'll need to define it to return the appropriate platform directory name.

# Main function to start binaries



kill_exist_binaries

result=$(check_binaries_stop)
ret_val=$?

if [ $ret_val -eq 0 ]; then
    echo "All binaries are stopped."
else
    echo "$result"
    echo "abort..."
    exit 1
fi


# Call the main function
start_binaries

check_binaries_running

print_listened_ports_by_binaries


