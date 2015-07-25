#!/bin/bash

apt-get update
apt-get install -y software-properties-common

# Install Docker
apt-get install -y docker.io
ln -sf /usr/bin/docker.io /usr/local/bin/docker

# Install Go, godeb will install the latest version of Go.
curl https://godeb.s3.amazonaws.com/godeb-amd64.tar.gz | tar zx -C /usr/local/bin
GOPATH=/go godeb install

# Setup Go
export GOPATH=/go
rm -rf $GOPATH/pkg/linux_amd64
echo 'GOPATH=/go' > /etc/profile.d/go.sh
echo 'PATH=$GOPATH/bin:$PATH' >> /etc/profile.d/go.sh

# Install supervisord
apt-get install -y supervisor
service supervisord start

# Place ENV variables in /home/vagrant/.bashrc
if ! grep -Fxq "# Go and ResourceD Evironment Variables" /home/vagrant/.bashrc ; then
    echo -e "\n# Go and ResourceD Evironment Variables" >> /home/vagrant/.bashrc
    echo -e ". /etc/profile.d/go.sh" >> /home/vagrant/.bashrc
    echo -e ". /etc/profile.d/resourced.sh" >> /home/vagrant/.bashrc
fi

# Install ResourceD
cd $GOPATH/src/github.com/resourced/resourced
GOPATH=/go go get ./... && GOPATH=/go go install github.com/resourced/resourced

# SUPERVISORD
# Setup ResourceD on port :55556
ln -fs /go/src/github.com/resourced/resourced/tests/data/script-init/supervisord/resourced.conf /etc/supervisor/conf.d/
supervisorctl update

# UPSTART
# Setup ResourceD on port :55555
# Log file can be found here: /var/log/upstart/resourced.log
ln -fs /go/src/github.com/resourced/resourced/tests/data/script-init/upstart/resourced.conf /etc/init/
initctl reload-configuration
service resourced start