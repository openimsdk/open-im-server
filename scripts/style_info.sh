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

function style-info() {
    COLOR_SUFFIX="\033[0m"  # End all colors and special effects

    BLACK_PREFIX="\033[30m"  # Black prefix
    RED_PREFIX="\033[31m"  # Red prefix
    GREEN_PREFIX="\033[32m"  # Green prefix
    YELLOW_PREFIX="\033[33m"  # Yellow prefix
    BLUE_PREFIX="\033[34m"  # Blue prefix
    PURPLE_PREFIX="\033[35m"  # Purple prefix
    SKY_BLUE_PREFIX="\033[36m"  # Sky blue prefix
    WHITE_PREFIX="\033[37m"  # White prefix
    BOLD_PREFIX="\033[1m"  # Bold prefix
    UNDERLINE_PREFIX="\033[4m"  # Underline prefix
    ITALIC_PREFIX="\033[3m"  # Italic prefix

    CYAN_PREFIX="033[0;36m"  # Cyan prefix

    BACKGROUND_BLACK="\033[40m"  # Black background
    BACKGROUND_RED="\033[41m"  # Red background
    BACKGROUND_GREEN="\033[42m"  # Green background
    BACKGROUND_YELLOW="\033[43m"  # Yellow background
    BACKGROUND_BLUE="\033[44m"  # Blue background
    BACKGROUND_PURPLE="\033[45m"  # Purple background
    BACKGROUND_SKY_BLUE="\033[46m"  # Sky blue background
    BACKGROUND_WHITE="\033[47m"  # White background

    BLINK="\033[5m"  # Blinking effect
    INVERT="\033[7m"  # Invert color
    HIDE="\033[8m"  # Hide text

    GRAY_PREFIX="\033[90m"  # Gray prefix
    LIGHT_RED_PREFIX="\033[91m"  # Light red prefix
    LIGHT_GREEN_PREFIX="\033[92m"  # Light green prefix
    LIGHT_YELLOW_PREFIX="\033[93m"  # Light yellow prefix
    LIGHT_BLUE_PREFIX="\033[94m"  # Light blue prefix
    LIGHT_PURPLE_PREFIX="\033[95m"  # Light purple prefix
    LIGHT_SKY_BLUE_PREFIX="\033[96m"  # Light sky blue prefix
    LIGHT_WHITE_PREFIX="\033[97m"  # Light white prefix

    BACKGROUND_GRAY="\033[100m"  # Gray background
    BACKGROUND_LIGHT_RED="\033[101m"  # Light red background
    BACKGROUND_LIGHT_GREEN="\033[102m"  # Light green background
    BACKGROUND_LIGHT_YELLOW="\033[103m"  # Light yellow background
    BACKGROUND_LIGHT_BLUE="\033[104m"  # Light blue background
    BACKGROUND_LIGHT_PURPLE="\033[105m"  # Light purple background
    BACKGROUND_LIGHT_SKY_BLUE="\033[106m"  # Light sky blue background
    BACKGROUND_LIGHT_WHITE="\033[107m"  # Light white background
}

style-info