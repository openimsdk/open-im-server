#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MODE="${1:-before-rpc}"

case "$MODE" in
  gateway)
    export OPENIM_WS_BREAK_ON_SEND=1
    ;;
  before-rpc)
    export OPENIM_WS_BREAK_BEFORE_SENDMSG_RPC=1
    ;;
  after-rpc)
    export OPENIM_WS_BREAK_AFTER_SENDMSG_RPC=1
    ;;
  *)
    echo "usage: $0 [gateway|before-rpc|after-rpc]" >&2
    exit 1
    ;;
esac

if [ -d /opt/homebrew/opt/openjdk ]; then
  export JAVA_HOME=/opt/homebrew/opt/openjdk
  export PATH="/opt/homebrew/bin:$JAVA_HOME/bin:$PATH"
fi

export PATH="$HOME/go/bin:$PATH"

cd "$ROOT_DIR"

exec dlv debug ./cmd -- -c ./config
