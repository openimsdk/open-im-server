# Ubuntu 22.04 OpenIM Project Development Guide

## TOC
- [Ubuntu 22.04 OpenIM Project Development Guide](#ubuntu-2204-openim-project-development-guide)
  - [TOC](#toc)
  - [1. Setting Up Ubuntu Server](#1-setting-up-ubuntu-server)
  - [1.1 Create `openim` Standard User](#11-create-openim-standard-user)
  - [1.2 Setting up the `openim` User's Shell Environment](#12-setting-up-the-openim-users-shell-environment)
  - [1.3 Installing Dependencies](#13-installing-dependencies)

## 1. Setting Up Ubuntu Server

You can use tools like PuTTY or other SSH clients to log in to your Ubuntu server. Once logged in, a few fundamental configurations are required, such as creating a standard user, adding to sudoers, and setting up the `$HOME/.bashrc` file. The steps are as follows:

## 1.1 Create `openim` Standard User

1. Log in to the Ubuntu system as the `root` user and create a standard user.

Generally, a project will involve multiple developers. Instead of provisioning a server for every developer, many organizations share a single development machine among developers. To simulate this real-world scenario, we'll use a standard user for development. To create the `openim` user:

```
# adduser openim # Create the openim user, which developers will use for login and development.
# passwd openim # Set the login password for openim.
```

Working with a non-root user ensures the system's safety and is a good practice. It's recommended to avoid using the root user as much as possible during everyday development.

1. Add to sudoers.

Often, even standard users need root privileges. Instead of frequently asking the system administrator for the root password, you can add the standard user to the sudoers. This allows them to temporarily gain root access using the sudo command. To add the `openim` user to sudoers:

```

# sed -i '/^root.*ALL=(ALL:ALL).*ALL/a\openim\tALL=(ALL) \tALL' /etc/sudoers
```

## 1.2 Setting up the `openim` User's Shell Environment

1. Log into the Ubuntu system.

Assuming we're using the **openim** user, log in using PuTTY or other SSH clients.

1. Configure the `$HOME/.bashrc` file.

The first step after logging into a new server is to configure the `$HOME/.bashrc` file. It makes the Linux shell more user-friendly by setting environment variables like `LANG` and `PS1`. Here's how the configuration would look:

```
# .bashrc

# User specific aliases and functions

alias rm='rm -i'
alias cp='cp -i'
alias mv='mv -i'

# Source global definitions
if [ -f /etc/bashrc ]; then
    . /etc/bashrc
fi

if [ ! -d $HOME/workspace ]; then
    mkdir -p $HOME/workspace
fi

# User specific environment
export LANG="en_US.UTF-8" 
export PS1='[\u@dev \W]\$ '
export WORKSPACE="$HOME/workspace"
export PATH=$HOME/bin:$PATH

cd $WORKSPACE
```

After updating `$HOME/.bashrc`, run the `bash` command to reload the configurations into the current shell.

## 1.3 Installing Dependencies

The OpenIM project on Ubuntu may have various dependencies. Some are direct, and others are indirect. Installing these in advance prevents issues like missing packages or compile-time errors later on.

1. Install dependencies.

You can use the `apt` command to install the required tools on Ubuntu:

```
$ sudo apt-get update 
$ sudo apt-get install build-essential autoconf automake cmake perl libcurl4-gnutls-dev libtool gcc g++ glibc-doc-reference zlib1g-dev git-lfs telnet lrzsz jq libexpat1-dev libssl-dev
$ sudo apt install libcurl4-openssl-dev
```

1. Install Git.

A higher version of Git ensures compatibility with certain commands like `git fetch --unshallow`. To install a recent version:

```
$ cd /tmp
$ wget --no-check-certificate https://mirrors.edge.kernel.org/pub/software/scm/git/git-2.36.1.tar.gz
$ tar -xvzf git-2.36.1.tar.gz
$ cd git-2.36.1/
$ ./configure
$ make
$ sudo make install
$ git --version          
```

Then, add Git's binary directory to the `PATH`:

```

$ echo 'export PATH=/usr/local/libexec/git-core:$PATH' >> $HOME/.bashrc
```

1. Configure Git.

To set up Git:

```
$ git config --global user.name "Your Name"
$ git config --global user.email "your_email@example.com"
$ git config --global credential.helper store
$ git config --global core.longpaths true
```

Other Git configurations include:

```

$ git config --global core.quotepath off
```

And for handling larger files:

```

$ git lfs install --skip-repo
```

By following the steps in this guide, your Ubuntu 22.04 server should now be set up and ready for OpenIM project development.