# OpenIM Bash Utility Script

This script offers a variety of utilities and helpers to enhance and simplify operations related to the OpenIM project.

## Table of Contents

- [OpenIM Bash Utility Script](#openim-bash-utility-script)
  - [Table of Contents](#table-of-contents)
  - [brief descriptions of each function](#brief-descriptions-of-each-function)
  - [Introduction](#introduction)
  - [Usage](#usage)
    - [SSH Key Setup](#ssh-key-setup)
  - [openim::util::ensure-gnu-sed](#openimutilensure-gnu-sed)
  - [openim::util::ensure-gnu-date](#openimutilensure-gnu-date)
  - [openim::util::check-file-in-alphabetical-order](#openimutilcheck-file-in-alphabetical-order)
  - [openim::util::require-jq](#openimutilrequire-jq)
  - [openim::util::md5](#openimutilmd5)
  - [openim::util::read-array](#openimutilread-array)
  - [Color Definitions](#color-definitions)
  - [openim::util::desc and related functions](#openimutildesc-and-related-functions)
  - [openim::util::onCtrlC](#openimutilonctrlc)
  - [openim::util::list-to-string](#openimutillist-to-string)
  - [openim::util::remove-space](#openimutilremove-space)
  - [openim::util::gencpu](#openimutilgencpu)
  - [openim::util::gen-os-arch](#openimutilgen-os-arch)
  - [openim::util::download-file](#openimutildownload-file)
  - [openim::util::get-public-ip](#openimutilget-public-ip)
  - [openim::util::extract-tarball](#openimutilextract-tarball)
  - [openim::util::check-port-open](#openimutilcheck-port-open)
  - [openim::util::file-lines-count](#openimutilfile-lines-count)


##  brief descriptions of each function 

**Englist:**
1. `openim::util::ensure-gnu-sed` - Determines if GNU version of `sed` exists on the system and sets its name.
2. `openim::util::ensure-gnu-date` - Determines if GNU version of `date` exists on the system and sets its name.
3. `openim::util::check-file-in-alphabetical-order` - Checks if a file is sorted in alphabetical order.
4. `openim::util::require-jq` - Checks if `jq` is installed.
5. `openim::util::md5` - Outputs the MD5 hash of a file.
6. `openim::util::read-array` - Reads content from standard input into an array.
7. `openim::util::desc` - Displays descriptive information.
8. `openim::util::run::prompt` - Displays a prompt.
9. `openim::util::run::maybe-first-prompt` - Possibly displays the first prompt based on whether it's started or not.
10. `openim::util::run` - Executes a command and captures its output.
11. `openim::util::run::relative` - Returns paths relative to the current script.
12. `openim::util::onCtrlC` - Performs an action when Ctrl+C is pressed.
13. `openim::util::list-to-string` - Converts a list into a string.
14. `openim::util::remove-space` - Removes spaces from a string.
15. `openim::util::gencpu` - Retrieves CPU information.
16. `openim::util::gen-os-arch` - Generates a repository directory based on the operating system and architecture.
17. `openim::util::download-file` - Downloads a file from a URL.
18. `openim::util::get-public-ip` - Retrieves the public IP address of the machine.
19. `openim::util::extract-tarball` - Extracts a tarball to a specified directory.
20. `openim::util::check-port-open` - Checks if a given port is open on the machine.
21. `openim::util::file-lines-count` - Counts the number of lines in a file.



## Introduction

This script is mainly used to validate whether the code is correctly formatted by `gofmt`. Apart from that, it offers utilities like setting up SSH keys, various wait conditions, host and platform detection, documentation generation, etc. 

## Usage

### SSH Key Setup

To set up an SSH key:

```bash
#1. Write IPs in a file, one IP per line. Let's name it hosts-file.
#2. Modify the default username and password in the script.
hosts-file-path="path/to/your/hosts/file"
openim:util::setup_ssh_key_copy "$hosts-file-path" "root" "123"
```

## openim::util::ensure-gnu-sed

Ensures the presence of the GNU version of the `sed` command. Different operating systems may have variations of the `sed` command, and this utility function is used to make sure the script uses the GNU version. If it finds the GNU `sed`, it sets the `SED` variable accordingly. If not found, it checks for `gsed`, which is usually the name of GNU `sed` on macOS. If neither is found, an error message is displayed.



## openim::util::ensure-gnu-date

Similar to the function for `sed`, this function ensures the script uses the GNU version of the `date` command. If it identifies the GNU `date`, it sets the `DATE` variable. On macOS, it looks for `gdate` as an alternative. In the absence of both, an error message is recommended.



## openim::util::check-file-in-alphabetical-order

This function checks if the contents of a given file are sorted in alphabetical order. If not, it provides a command suggestion for the user to sort the file correctly.



## openim::util::require-jq

Verifies the installation of `jq`, a popular command-line JSON parser. If it's not present, a prompt to install it is displayed.



## openim::util::md5

A cross-platform function that computes the MD5 hash of its input. This function takes into account the differences in the `md5` command between macOS and Linux.



## openim::util::read-array

A function designed to read from stdin and populate an array, line by line. It's provided as an alternative to `mapfile -t` and is compatible with bash 3.



## Color Definitions

The script also defines a set of colors to enhance its console output. These include colors like red, yellow, green, blue, cyan, etc., which can be used for better user experience and clear logs.



## openim::util::desc and related functions

These functions seem to aid in building interactive demonstrations or tutorials in the terminal. They use the `pv` utility to control the display rate of the output, emulating typing. There's also functionality to handle user prompts and execute commands while capturing their output.



## openim::util::onCtrlC

Handles the `CTRL+C` command. It terminates background processes of the script when the user interrupts it using `CTRL+C`.



## openim::util::list-to-string

Transforms a list format (like `[10023, 2323, 3434]`) to a space-separated string (`10023 2323 3434`). Also removes unnecessary spaces and characters.



## openim::util::remove-space

Removes spaces from a given string.



## openim::util::gencpu

Fetches the number of CPUs using the `lscpu` command.



## openim::util::gen-os-arch

Identifies the operating system and architecture of the system running the script. This is useful to determine directories or binaries specific to that OS and architecture.



## openim::util::download-file

This function can be used to download a file from a URL. If `curl` is available, it uses `curl`. If not, it falls back to `wget`.

```bash
function openim::util::download-file() {
  local url="$1"
  local dest="$2"

  if command -v curl &>/dev/null; then
    curl -L "${url}" -o "${dest}"
  elif command -v wget &>/dev/null; then
    wget "${url}" -O "${dest}"
  else
    openim::log::error "Neither curl nor wget available. Cannot download file."
    return 1
  fi
}
```



## openim::util::get-public-ip

Fetches the public IP address of the machine.

```bash
function openim::util::get-public-ip() {
  if command -v curl &>/dev/null; then
    curl -s https://ipinfo.io/ip
  elif command -v wget &>/dev/null; then
    wget -qO- https://ipinfo.io/ip
  else
    openim::log::error "Neither curl nor wget available. Cannot fetch public IP."
    return 1
  fi
}
```



## openim::util::extract-tarball

This function extracts a tarball to a specified directory.

```bash
function openim::util::extract-tarball() {
  local tarball="$1"
  local dest="$2"

  mkdir -p "${dest}"
  tar -xzf "${tarball}" -C "${dest}"
}
```



## openim::util::check-port-open

Checks if a given port is open on the local machine.

```bash
function openim::util::check-port-open() {
  local port="$1"
  if command -v nc &>/dev/null; then
    echo -n > /dev/tcp/127.0.0.1/"${port}" 2>&1
    return $?
  elif command -v telnet &>/dev/null; then
    telnet 127.0.0.1 "${port}" 2>&1 | grep -q "Connected"
    return $?
  else
    openim::log::error "Neither nc nor telnet available. Cannot check port."
    return 1
  fi
}
```



## openim::util::file-lines-count

Counts the number of lines in a file.

```bash
function openim::util::file-lines-count() {
  local file="$1"
  if [[ -f "${file}" ]]; then
    wc -l < "${file}"
  else
    openim::log::error "File does not exist: ${file}"
    return 1
  fi
}
```