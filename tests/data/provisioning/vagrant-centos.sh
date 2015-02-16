#!/bin/bash

yum install -y golang

# Setup Go
export GOPATH=/go
rm -rf $GOPATH/pkg/linux_amd64
echo 'GOPATH=/go' > /etc/profile.d/go.sh
echo 'PATH=$GOPATH/bin:$PATH' >> /etc/profile.d/go.sh

# Install supervisord
yum install -y supervisor

# Install ResourceD
mkdir -p $GOPATH/src/github.com/resourced/resourced && cd $GOPATH/src/github.com/resourced/resourced
GOPATH=/go go get ./... && GOPATH=/go go install github.com/resourced/resourced
mkdir -p /resourced && echo 'RESOURCED_DB=/resourced/db' > /etc/profile.d/resourced.sh

# Place ENV variables in /home/vagrant/.bashrc
if ! grep -Fxq "# Go and ResourceD Evironment Variables" /home/vagrant/.bashrc ; then
    echo -e "\n# Go and ResourceD Evironment Variables" >> /home/vagrant/.bashrc
    echo -e ". /etc/profile.d/go.sh" >> /home/vagrant/.bashrc
    echo -e ". /etc/profile.d/resourced.sh" >> /home/vagrant/.bashrc
fi
