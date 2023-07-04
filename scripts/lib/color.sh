#!/usr/bin/env bash

# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

#Define color variables
#Feature
COLOR_NORMAL='\033[0m';COLOR_BOLD='\033[1m';COLOR_DIM='\033[2m';COLOR_UNDER='\033[4m';
COLOR_ITALIC='\033[3m';COLOR_NOITALIC='\033[23m';COLOR_BLINK='\033[5m';
COLOR_REVERSE='\033[7m';COLOR_CONCEAL='\033[8m';COLOR_NOBOLD='\033[22m';
COLOR_NOUNDER='\033[24m';COLOR_NOBLINK='\033[25m';

#Front color
COLOR_BLACK='\033[30m';COLOR_RED='\033[31m';COLOR_GREEN='\033[32m';COLOR_YELLOW='\033[33m';
COLOR_BLUE='\033[34m';COLOR_MAGENTA='\033[35m';COLOR_CYAN='\033[36m';COLOR_WHITE='\033[37m';

#background color
COLOR_BBLACK='\033[40m';COLOR_BRED='\033[41m';
COLOR_BGREEN='\033[42m';COLOR_BYELLOW='\033[43m';
COLOR_BBLUE='\033[44m';COLOR_BMAGENTA='\033[45m';
COLOR_BCYAN='\033[46m';COLOR_BWHITE='\033[47m';

# Print colors you can use
iam::color::print_color()
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
