#!/bin/bash

apt-get update
apt-get install -y software-properties-common

# Install Docker
apt-get install -y docker.io
ln -sf /usr/bin/docker.io /usr/local/bin/docker

# Install Go
add-apt-repository ppa:gophers/go
apt-get update
apt-get install -y golang

# Set Go ENV variable
export GOPATH=/go
echo 'GOPATH=/go' > /etc/profile.d/go.sh
echo 'PATH=$GOPATH/bin:$PATH' >> /etc/profile.d/go.sh

# Install ResourceD
cd $GOPATH/src/github.com/resourced/resourced && go get ./...
cd $GOPATH/src/github.com/resourced/resourced && go install github.com/resourced/resourced
mkdir -p /resourced
echo 'RESOURCED_DB=/resourced/db' > /etc/profile.d/resourced.sh
