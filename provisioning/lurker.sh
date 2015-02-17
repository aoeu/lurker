#!/usr/bin/env bash

# Provisioning script for maestro
# because ansible is too fucking complicated

echo "Provisioning..."
sudo apt-get -y update
sudo apt-get -y install build-essential libssl-dev
sudo apt-get -y install pkg-config
sudo apt-get -y install git

# install go
echo "Installing golang 1.4.1"
wget https://storage.googleapis.com/golang/go1.4.1.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.4.1.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc 
echo 'export GOPATH=$HOME/go' >> ~/.bashrc 
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc 
sudo chown -R vagrant:vagrant /home/vagrant/go

# install libxml for gokogiri
sudo apt-get -y install libxml2-dev


# Some obstinate sites may have too much JS for their own good.
# Use this to proxy lurker's requests to statically render the site.
# go install github.com/sourcegraph/webloop/...
# static-reverse-proxy

echo "Good Hunting"