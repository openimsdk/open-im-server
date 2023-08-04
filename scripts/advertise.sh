#!/usr/bin/env bash

set -e
set -o pipefail

trap 'echo "Script interrupted."; exit 1' INT

# Function for colored echo
function color_echo() {
    COLOR=$1
    shift
    echo -e "${COLOR}===> $* ${COLOR_SUFFIX}"
}

# Color definitions
function openim_color() {
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
}

function print_with_delay() {
  text="$1"
  delay="$2"
  color="$3"

  for i in $(seq 0 $((${#text}-1))); do
    printf "${color}${text:$i:1}${COLOR_SUFFIX}"
    sleep $delay
  done
  printf "\n"
}

function print_progress() {
  total="$1"
  delay="$2"
  color="$3"

  printf "${color}["
  for i in $(seq 1 $total); do
    printf "#"
    sleep $delay
  done
  printf "]${COLOR_SUFFIX}\n"
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
    print_with_delay "Discover more and contribute at: https://github.com/OpenIMSDK/Open-IM-Server" 0.01

    # Reset text color back to normal
    echo -e "\033[0m"

    # Set text color to green for product description
    echo -e "\033[1;32m"

    print_with_delay "Open-IM-Server: Reinventing Instant Messaging" 0.01
    print_progress 50 0.02

    print_with_delay "Open-IM-Server is not just a product; it's a revolution. It's about bringing the power of seamless, real-time messaging to your fingertips. And it's about joining a global community of developers, dedicated to pushing the boundaries of what's possible." 0.01

    print_progress 50 0.02

    # Reset text color back to normal
    echo -e "\033[0m"

    # Set text color to yellow for the Slack link
    echo -e "\033[1;33m"

    print_with_delay "Join our developer community on Slack: https://join.slack.com/t/openimsdk/shared_invite/zt-1tmoj26uf-_FDy3dowVHBiGvLk9e5Xkg" 0.01

    # Reset text color back to normal
    echo -e "\033[0m"
}

function main() {
    openim_logo
}
main "$@"
