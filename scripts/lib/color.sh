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

# this script is used to install the dependencies of the project
#
# Usage: `scripts/color.sh`.
################################################################################

# shellcheck disable=SC2034
if [[ ! -v COLOR_OPEN ]]; then
    COLOR_OPEN=1
fi

# Function for colored echo
openim::color::echo() {
    COLOR=$1
    [ $COLOR_OPEN -eq 1 ]  && echo -e "${COLOR} $(date '+%Y-%m-%d %H:%M:%S') $@ ${COLOR_SUFFIX}"
    shift
}

# Define color variables
# --- Feature --- 
COLOR_NORMAL='\033[0m';COLOR_BOLD='\033[1m';COLOR_DIM='\033[2m';COLOR_UNDER='\033[4m';
COLOR_ITALIC='\033[3m';COLOR_NOITALIC='\033[23m';COLOR_BLINK='\033[5m';
COLOR_REVERSE='\033[7m';COLOR_CONCEAL='\033[8m';COLOR_NOBOLD='\033[22m';
COLOR_NOUNDER='\033[24m';COLOR_NOBLINK='\033[25m';

# --- Front color --- 
COLOR_BLACK='\033[30m';COLOR_RED='\033[31m';COLOR_GREEN='\033[32m';COLOR_YELLOW='\033[33m';
COLOR_BLUE='\033[34m';COLOR_MAGENTA='\033[35m';COLOR_CYAN='\033[36m';COLOR_WHITE='\033[37m';

# --- background color --- 
COLOR_BBLACK='\033[40m';COLOR_BRED='\033[41m';
COLOR_BGREEN='\033[42m';COLOR_BYELLOW='\033[43m';
COLOR_BBLUE='\033[44m';COLOR_BMAGENTA='\033[45m';
COLOR_BCYAN='\033[46m';COLOR_BWHITE='\033[47m';

# --- Color definitions --- 
# Color definitions
COLOR_SUFFIX="\033[0m"      # End all colors and special effects
BLACK_PREFIX="\033[30m"     # Black prefix
RED_PREFIX="\033[31m"       # Red prefix
GREEN_PREFIX="\033[32m"     # Green prefix
YELLOW_PREFIX="\033[33m"    # Yellow prefix
BLUE_PREFIX="\033[34m"      # Blue prefix
SKY_BLUE_PREFIX="\033[36m"  # Sky blue prefix
WHITE_PREFIX="\033[37m"     # White prefix
BOLD_PREFIX="\033[1m"       # Bold prefix
UNDERLINE_PREFIX="\033[4m"  # Underline prefix
ITALIC_PREFIX="\033[3m"     # Italic prefix
CYAN_PREFIX="\033[0;36m"     # Cyan prefix

# --- make demo (run demo) ---
reset=$(tput sgr0)
bold=$(tput bold)
black=$(tput setaf 0)
red=$(tput setaf 1)
green=$(tput bold; tput setaf 2)
yellow=$(tput bold; tput setaf 3)
blue=$(tput bold; tput setaf 6)
timeout=$(if [ "$(uname)" == "Darwin" ]; then echo "1"; else echo "0.1"; fi)


# Print colors you can use
openim::color::print_color()
{
  echo
  echo -e ${bmagenta}--back-color:${normal}
  echo "bblack; bgreen; bblue; bcyan; bred; byellow; bmagenta; bwhite"
  echo
  echo -e ${red}--font-color:${normal}
  echo "black; red; green; yellow; blue; magenta; cyan; white"
  echo
  echo -e ${bold}--font:${normal}
  echo "normal; italic; reverse; nounder; bold; noitalic; conceal; noblink;
  dim; blink; nobold; under"
  echo
}
