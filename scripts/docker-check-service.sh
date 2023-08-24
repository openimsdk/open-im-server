#!/usr/bin/env bash

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

OPENIM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${OPENIM_ROOT}/scripts/install/common.sh"

cd "$OPENIM_ROOT"

if command -v docker-compose &> /dev/null; then
    docker-compose ps
else
    docker compose ps
fi

progress() {
    local _main_pid="$1"
    local _length=20
    local _ratio=1
    local _colors=("31" "32" "33" "34" "35" "36" "37")
    local _wave=("▁" "▂" "▃" "▄" "▅" "▆" "▇" "█" "▇" "▆" "▅" "▄" "▃" "▂")

    while pgrep -P "$_main_pid" &> /dev/null; do
        local _mark='>'
        local _progress_bar=
        for ((i = 1; i <= _length; i++)); do
            if ((i > _ratio)); then
                _mark='-'
            fi
            _progress_bar="${_progress_bar}${_mark}"
        done

        local _color_idx=$((_ratio % ${#_colors[@]}))
        local _color_prefix="\033[${_colors[_color_idx]}m"
        local _reset_suffix="\033[0m"

        local _wave_idx=$((_ratio % ${#_wave[@]}))
        local _wave_progress=${_wave[_wave_idx]}

        printf "Progress: ${_color_prefix}${_progress_bar}${_reset_suffix} ${_wave_progress}   Countdown: %2ds \r" "$_countdown"
        ((_ratio++))
        ((_ratio > _length)) && _ratio=1
        sleep 0.1
    done
}

countdown() {
    local _duration="$1"

    for ((i = _duration; i >= 1; i--)); do
        printf "\rCountdown: %2ds \r" "$i"
        sleep 1
    done
    printf "\rCountdown: %2ds \r" "$_duration"
}

do_sth() {
    echo "++++++++++++++++++++++++"
    progress $$ &
    local _progress_pid=$!
    local _countdown=30

    countdown "$_countdown" &
    local _countdown_pid=$!

    sleep 30

    kill "$_progress_pid" "$_countdown_pid"

    "${SCRIPTS_ROOT}/check-all.sh"
    echo -e "${PURPLE_PREFIX}=========> Check docker-compose status ${COLOR_SUFFIX} \n"
}

set -e

do_sth &
do_sth_pid=$(jobs -p | tail -1)

progress "${do_sth_pid}" &
progress_pid=$(jobs -p | tail -1)

wait "${do_sth_pid}"
printf "Progress: done                \n"
