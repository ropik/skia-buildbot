#!/bin/bash
#
# The continue_install script updates the webtry user's copy of depot_tools
# and the buildbot repository.
# It then builds the webtry server outside the jail.
#
# The setup scripts should run this script as the 'webtry' user.
#
# See the README file for detailed installation instructions.

# Install Go

cd

if [ -d go ]; then
  echo Go already installed.
else
  wget https://storage.googleapis.com/golang/go1.3.3.linux-amd64.tar.gz
  tar -xzf go1.3.3.linux-amd64.tar.gz
fi

mkdir ${HOME}/golib
export GOROOT=${HOME}/go
export GOPATH=${HOME}/golib
export PATH=$PATH:$GOROOT/bin

# Install depot_tools.
if [ -d depot_tools ]; then
  (cd depot_tools && git pull);
else
  git clone https://chromium.googlesource.com/chromium/tools/depot_tools.git;
fi
export PATH=$PATH:~/depot_tools

# Checkout the skia code (needed so that webtry can build its template).
fetch skia

cd

# Checkout the buildbot code and its dependencies.
git clone https://skia.googlesource.com/buildbot.git

cd buildbot/webtry

go get -d
./build
