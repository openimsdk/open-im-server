
#!/usr/bin/env bash



OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/
source "${OPENIM_ROOT}/lib/util.sh"
source "${OPENIM_ROOT}/define/binaries.sh"
source "${OPENIM_ROOT}/lib/path.sh"



for binary in "${!binaries[@]}"; do
  expected_count=${binaries[$binary]}
  full_path=$(get_bin_full_path "$binary")
  kill_binary "$full_path"
done


