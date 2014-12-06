#!/usr/bin/env bash

# Provisioning script for maestro
# because ansible is too fucking complicated

echo "Provisioning..."
sudo apt-get -y update
sudo apt-get -y install build-essential libssl-dev

# install go
wget https://storage.googleapis.com/golang/go1.3.3.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.3.3.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

echo "Good Hunting"