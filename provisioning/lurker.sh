#!/usr/bin/env bash

# Provisioning script for maestro
# because ansible is too fucking complicated

echo "Provisioning..."
sudo apt-get -y update
sudo apt-get -y install build-essential libssl-dev
sudo apt-get -y install git
# install go
wget https://storage.googleapis.com/golang/go1.3.3.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.3.3.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc 
echo 'export GOPATH=$HOME/go' >> ~/.bashrc 
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc 
sudo chown -R vagrant:vagrant /home/vagrant/go
echo "Good Hunting"