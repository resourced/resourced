#!/bin/bash

apt-get update
apt-get install -y software-properties-common python-setuptools sysstat

# Install Docker
apt-get install -y docker.io
ln -sf /usr/bin/docker.io /usr/local/bin/docker

# Remove all containers
docker stop $(docker ps -a -q)
docker rm $(docker ps -a -q)

# Run one container for testing
docker pull nginx:latest
docker run -d nginx:latest

# Install mysql
export DEBIAN_FRONTEND=noninteractive
apt-get -q -y install mysql-server

# Install Redis
apt-get -q -y install redis-server
service redis-server restart

# Install memcache
apt-get -q -y install memcached
service memcached restart

# Install varnish
apt-get -q -y install varnish
service varnish restart

# Install haproxy
apt-get -q -y install haproxy
service haproxy restart

# Install nginx
apt-get -q -y install nginx
service nginx restart

# Install Go, godeb will install the latest version of Go.
curl https://godeb.s3.amazonaws.com/godeb-amd64.tar.gz | tar zx -C /usr/local/bin
GOPATH=/go godeb install

# Setup Go
export GOPATH=/go
rm -rf $GOPATH/pkg/linux_amd64
echo 'export GOPATH=/go' > /etc/profile.d/go.sh
echo 'export PATH=$GOPATH/bin:$PATH' >> /etc/profile.d/go.sh

# Install supervisord
apt-get install -y supervisor
easy_install superlance
service supervisor restart

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
ln -fs /go/src/github.com/resourced/resourced/tests/script-init/supervisord/resourced.conf /etc/supervisor/conf.d/
supervisorctl update

# UPSTART
# Setup ResourceD on port :55555
# Log file can be found here: /var/log/upstart/resourced.log
ln -fs /go/src/github.com/resourced/resourced/tests/script-init/upstart/resourced.conf /etc/init/
initctl reload-configuration
service resourced restart