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

set -e
set -o pipefail

. $(dirname ${BASH_SOURCE})/lib/init.sh

trap 'openim::util::onCtrlC' INT

print_with_delay() {
  text="$1"
  delay="$2"

  for i in $(seq 0 $((${#text}-1))); do
    printf "${text:$i:1}"
    sleep $delay
  done
  printf "\n"
}

print_progress() {
  total="$1"
  delay="$2"

  printf "["
  for i in $(seq 1 $total); do
    printf "#"
    sleep $delay
  done
  printf "]\n"
}

function openim_logo() {
    # Set text color to cyan for header and URL
    echo -e "\033[0;36m"

    # Display fancy ASCII Art logo
    # look http://patorjk.com/software/taag/#p=display&h=1&v=1&f=Doh&t=OpenIM
    print_with_delay '
                                                                                                                      
                                                                                                                      
     OOOOOOOOO                                                               IIIIIIIIIIMMMMMMMM               MMMMMMMM
   OO:::::::::OO                                                             I::::::::IM:::::::M             M:::::::M
 OO:::::::::::::OO                                                           I::::::::IM::::::::M           M::::::::M
O:::::::OOO:::::::O                                                          II::::::IIM:::::::::M         M:::::::::M
O::::::O   O::::::Oppppp   ppppppppp       eeeeeeeeeeee    nnnn  nnnnnnnn      I::::I  M::::::::::M       M::::::::::M
O:::::O     O:::::Op::::ppp:::::::::p    ee::::::::::::ee  n:::nn::::::::nn    I::::I  M:::::::::::M     M:::::::::::M
O:::::O     O:::::Op:::::::::::::::::p  e::::::eeeee:::::een::::::::::::::nn   I::::I  M:::::::M::::M   M::::M:::::::M
O:::::O     O:::::Opp::::::ppppp::::::pe::::::e     e:::::enn:::::::::::::::n  I::::I  M::::::M M::::M M::::M M::::::M
O:::::O     O:::::O p:::::p     p:::::pe:::::::eeeee::::::e  n:::::nnnn:::::n  I::::I  M::::::M  M::::M::::M  M::::::M
O:::::O     O:::::O p:::::p     p:::::pe:::::::::::::::::e   n::::n    n::::n  I::::I  M::::::M   M:::::::M   M::::::M
O:::::O     O:::::O p:::::p     p:::::pe::::::eeeeeeeeeee    n::::n    n::::n  I::::I  M::::::M    M:::::M    M::::::M
O::::::O   O::::::O p:::::p    p::::::pe:::::::e             n::::n    n::::n  I::::I  M::::::M     MMMMM     M::::::M
O:::::::OOO:::::::O p:::::ppppp:::::::pe::::::::e            n::::n    n::::nII::::::IIM::::::M               M::::::M
 OO:::::::::::::OO  p::::::::::::::::p  e::::::::eeeeeeee    n::::n    n::::nI::::::::IM::::::M               M::::::M
   OO:::::::::OO    p::::::::::::::pp    ee:::::::::::::e    n::::n    n::::nI::::::::IM::::::M               M::::::M
     OOOOOOOOO      p::::::pppppppp        eeeeeeeeeeeeee    nnnnnn    nnnnnnIIIIIIIIIIMMMMMMMM               MMMMMMMM
                    p:::::p                                                                                           
                    p:::::p                                                                                           
                   p:::::::p                                                                                          
                   p:::::::p                                                                                          
                   p:::::::p                                                                                          
                   ppppppppp                                                                                          
                                                                                                                      
    ' 0.0001

    # Display product URL
    print_with_delay "Discover more and contribute at: https://github.com/openimsdk/open-im-server" 0.01

    # Reset text color back to normal
    echo -e "\033[0m"

    # Set text color to green for product description
    echo -e "\033[1;32m"

    print_with_delay "Open-IM-Server: Reinventing Instant Messaging" 0.01
    print_progress 50 0.02

    print_with_delay "Open-IM-Server is not just a product; it's a revolution. It's about bringing the power of seamless," 0.01
    print_with_delay "real-time messaging to your fingertips. And it's about joining a global community of developers, dedicated to pushing the boundaries of what's possible." 0.01

    print_progress 50 0.02

    # Reset text color back to normal
    echo -e "\033[0m"

    # Set text color to yellow for the Slack link
    echo -e "\033[1;33m"

    print_with_delay "Join our developer community on Slack: https://join.slack.com/t/openimsdk/shared_invite/zt-22720d66b-o_FvKxMTGXtcnnnHiMqe9Q" 0.01

    # Reset text color back to normal
    echo -e "\033[0m"
}

function main() {
    openim_logo
}
main "$@"
