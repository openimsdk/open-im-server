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

# A set of helpers for starting/running chat for tests

CHAT_VERSION=${CHAT_VERSION:-1.1.0}
CHAT_HOST=${CHAT_HOST:-127.0.0.1}
CHAT_PORT=${CHAT_PORT:-2379}
# This is intentionally not called CHAT_LOG_LEVEL:
# chat checks that and compains when it is set in addition
# to the command line argument, even when both have the same value.
CHAT_LOGLEVEL=${CHAT_LOGLEVEL:-warn}
export OPENIM_INTEGRATION_CHAT_URL="http://${CHAT_HOST}:${CHAT_PORT}"

openim::chat::validate() {
  # validate if in path
  command -v chat >/dev/null || {
    openim::log::usage "chat must be in your PATH"
    openim::log::info "You can use 'hack/install-chat.sh' to install a copy in third_party/."
    exit 1
  }

  # validate chat port is free
  local port_check_command
  if command -v ss &> /dev/null && ss -Version | grep 'iproute2' &> /dev/null; then
    port_check_command="ss"
  elif command -v netstat &>/dev/null; then
    port_check_command="netstat"
  else
    openim::log::usage "unable to identify if chat is bound to port ${CHAT_PORT}. unable to find ss or netstat utilities."
    exit 1
  fi
  if ${port_check_command} -nat | grep "LISTEN" | grep "[\.:]${CHAT_PORT:?}" >/dev/null 2>&1; then
    openim::log::usage "unable to start chat as port ${CHAT_PORT} is in use. please stop the process listening on this port and retry."
    openim::log::usage "$(${port_check_command} -nat | grep "LISTEN" | grep "[\.:]${CHAT_PORT:?}")"
    exit 1
  fi

  # need set the env of "CHAT_UNSUPPORTED_ARCH" on unstable arch.
  arch=$(uname -m)
  if [[ $arch =~ arm* ]]; then
	  export CHAT_UNSUPPORTED_ARCH=arm
  fi
  # validate installed version is at least equal to minimum
  version=$(chat --version | grep Version | head -n 1 | cut -d " " -f 3)
  if [[ $(openim::chat::version "${CHAT_VERSION}") -gt $(openim::chat::version "${version}") ]]; then
   export PATH="${OPENIM_ROOT}"/third_party/chat:${PATH}
   hash chat
   echo "${PATH}"
   version=$(chat --version | grep Version | head -n 1 | cut -d " " -f 3)
   if [[ $(openim::chat::version "${CHAT_VERSION}") -gt $(openim::chat::version "${version}") ]]; then
    openim::log::usage "chat version ${CHAT_VERSION} or greater required."
    openim::log::info "You can use 'hack/install-chat.sh' to install a copy in third_party/."
    exit 1
   fi
  fi
}

openim::chat::version() {
  printf '%s\n' "${@}" | awk -F . '{ printf("%d%03d%03d\n", $1, $2, $3) }'
}

openim::chat::start() {
  # validate before running
  openim::chat::validate

  # Start chat
  CHAT_DIR=${CHAT_DIR:-$(mktemp -d 2>/dev/null || mktemp -d -t test-chat.XXXXXX)}
  if [[ -d "${ARTIFACTS:-}" ]]; then
    CHAT_LOGFILE="${ARTIFACTS}/chat.$(uname -n).$(id -un).log.DEBUG.$(date +%Y%m%d-%H%M%S).$$"
  else
    CHAT_LOGFILE=${CHAT_LOGFILE:-"/dev/null"}
  fi
  openim::log::info "chat --advertise-client-urls ${OPENIM_INTEGRATION_CHAT_URL} --data-dir ${CHAT_DIR} --listen-client-urls http://${CHAT_HOST}:${CHAT_PORT} --log-level=${CHAT_LOGLEVEL} 2> \"${CHAT_LOGFILE}\" >/dev/null"
  chat --advertise-client-urls "${OPENIM_INTEGRATION_CHAT_URL}" --data-dir "${CHAT_DIR}" --listen-client-urls "${OPENIM_INTEGRATION_CHAT_URL}" --log-level="${CHAT_LOGLEVEL}" 2> "${CHAT_LOGFILE}" >/dev/null &
  CHAT_PID=$!

  echo "Waiting for chat to come up."
  openim::util::wait_for_url "${OPENIM_INTEGRATION_CHAT_URL}/health" "chat: " 0.25 80
  curl -fs -X POST "${OPENIM_INTEGRATION_CHAT_URL}/v3/kv/put" -d '{"key": "X3Rlc3Q=", "value": ""}'
}

openim::chat::start_scraping() {
  if [[ -d "${ARTIFACTS:-}" ]]; then
    CHAT_SCRAPE_DIR="${ARTIFACTS}/chat-scrapes"
  else
    CHAT_SCRAPE_DIR=$(mktemp -d -t test.XXXXXX)/chat-scrapes
  fi
  openim::log::info "Periodically scraping chat to ${CHAT_SCRAPE_DIR} ."
  mkdir -p "${CHAT_SCRAPE_DIR}"
  (
    while sleep 30; do
      openim::chat::scrape
    done
  ) &
  CHAT_SCRAPE_PID=$!
}

openim::chat::scrape() {
    curl -s -S "${OPENIM_INTEGRATION_CHAT_URL}/metrics" > "${CHAT_SCRAPE_DIR}/next" && mv "${CHAT_SCRAPE_DIR}/next" "${CHAT_SCRAPE_DIR}/$(date +%s).scrape"
}

openim::chat::stop() {
  if [[ -n "${CHAT_SCRAPE_PID:-}" ]] && [[ -n "${CHAT_SCRAPE_DIR:-}" ]] ; then
    kill "${CHAT_SCRAPE_PID}" &>/dev/null || :
    wait "${CHAT_SCRAPE_PID}" &>/dev/null || :
    openim::chat::scrape || :
    (
      # shellcheck disable=SC2015
      cd "${CHAT_SCRAPE_DIR}"/.. && \
      tar czf chat-scrapes.tgz chat-scrapes && \
      rm -rf chat-scrapes || :
    )
  fi
  if [[ -n "${CHAT_PID-}" ]]; then
    kill "${CHAT_PID}" &>/dev/null || :
    wait "${CHAT_PID}" &>/dev/null || :
  fi
}

openim::chat::clean_chat_dir() {
  if [[ -n "${CHAT_DIR-}" ]]; then
    rm -rf "${CHAT_DIR}"
  fi
}

openim::chat::cleanup() {
  openim::chat::stop
  openim::chat::clean_chat_dir
}

openim::chat::install() {
  (
    local os
    local arch

    os=$(openim::util::host_os)
    arch=$(openim::util::host_arch)

    cd ""${OPENIM_ROOT}"/third_party" || return 1
    if [[ $(readlink chat) == chat-v${CHAT_VERSION}-${os}-* ]]; then
      openim::log::info "chat v${CHAT_VERSION} already installed. To use:"
      openim::log::info "export PATH=\"$(pwd)/chat:\${PATH}\""
      return  #already installed
    fi

    if [[ ${os} == "darwin" ]]; then
      download_file="chat-v${CHAT_VERSION}-${os}-${arch}.zip"
      url="https://github.com/chat-io/chat/releases/download/v${CHAT_VERSION}/${download_file}"
      openim::util::download_file "${url}" "${download_file}"
      unzip -o "${download_file}"
      ln -fns "chat-v${CHAT_VERSION}-${os}-${arch}" chat
      rm "${download_file}"
    elif [[ ${os} == "linux" ]]; then
      url="https://github.com/coreos/chat/releases/download/v${CHAT_VERSION}/chat-v${CHAT_VERSION}-${os}-${arch}.tar.gz"
      download_file="chat-v${CHAT_VERSION}-${os}-${arch}.tar.gz"
      openim::util::download_file "${url}" "${download_file}"
      tar xzf "${download_file}"
      ln -fns "chat-v${CHAT_VERSION}-${os}-${arch}" chat
      rm "${download_file}"
    else
      openim::log::info "${os} is NOT supported."
    fi
    openim::log::info "chat v${CHAT_VERSION} installed. To use:"
    openim::log::info "export PATH=\"$(pwd)/chat:\${PATH}\""
  )
}