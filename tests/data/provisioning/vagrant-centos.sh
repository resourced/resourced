#!/bin/bash

# yum -y update
yum -y install git golang

# Setup Go
export GOPATH=/go
rm -rf $GOPATH/pkg/linux_amd64
echo 'GOPATH=/go' > /etc/profile.d/go.sh
echo 'PATH=$GOPATH/bin:$PATH' >> /etc/profile.d/go.sh

# Place ENV variables in /home/vagrant/.bashrc
if ! grep -Fxq "# Go and ResourceD Evironment Variables" /home/vagrant/.bashrc ; then
    echo -e "\n# Go and ResourceD Evironment Variables" >> /home/vagrant/.bashrc
    echo -e ". /etc/profile.d/go.sh" >> /home/vagrant/.bashrc
    echo -e ". /etc/profile.d/resourced.sh" >> /home/vagrant/.bashrc
fi

GOPATH=/go go get github.com/tools/godep

# Install ResourceD
cd $GOPATH/src/github.com/resourced/resourced
GOPATH=/go go get ./...
GOPATH=/go godep go build
mv $GOPATH/src/github.com/resourced/resourced/resourced /go/bin/resourced

# SYSTEMD
# Setup ResourceD on port :55555
rm -f /etc/systemd/user/resourced.service && cp /go/src/github.com/resourced/resourced/tests/data/script-init/systemd/resourced.service /etc/systemd/user/resourced.service
systemctl enable /etc/systemd/user/resourced.service
systemctl start resourced
