#!/bin/bash
set -x
set -e

#
# This script creates both Darwin and Linux binary distribution complete with default config files.
#

#
# VERSION is required
#
: "${VERSION?You must set VERSION}"

CURRENT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT_DIR=$(dirname $CURRENT_DIR)

cd $ROOT_DIR
git checkout tests/resourced-configs/*.toml

echo
echo "Creating directories"
mkdir -p $ROOT_DIR/resourced-darwin-$VERSION
mkdir -p $ROOT_DIR/resourced-darwin-$VERSION/{readers,writers,executors,tags,access-tokens}

mkdir -p $ROOT_DIR/resourced-linux-$VERSION
mkdir -p $ROOT_DIR/resourced-linux-$VERSION/{readers,writers,executors,tags,access-tokens}

echo
echo "Copying default readers"
cp -r $ROOT_DIR/tests/resourced-configs/readers/* $ROOT_DIR/resourced-darwin-$VERSION/readers/
cp -r $ROOT_DIR/tests/resourced-configs/readers/* $ROOT_DIR/resourced-linux-$VERSION/readers/
rm -f $ROOT_DIR/resourced-{darwin,linux}-$VERSION/readers/{docker,mysql,redis,varnish,haproxy,darwin,mcrouter}-*

echo
echo "Copying default writers"
cp -r $ROOT_DIR/tests/resourced-configs/writers/resourced-master-* $ROOT_DIR/resourced-darwin-$VERSION/writers/
cp -r $ROOT_DIR/tests/resourced-configs/writers/resourced-master-* $ROOT_DIR/resourced-linux-$VERSION/writers/

echo
echo "Copying default general.toml"
cp -r $ROOT_DIR/tests/resourced-configs/general.toml $ROOT_DIR/resourced-darwin-$VERSION/
cp -r $ROOT_DIR/tests/resourced-configs/general.toml $ROOT_DIR/resourced-linux-$VERSION/

echo
echo "Compiling Darwin binary"
GOOS=darwin go build && mv $ROOT_DIR/resourced $ROOT_DIR/resourced-darwin-$VERSION/

echo
echo "Compiling Linux binary"
GOOS=linux go build && mv $ROOT_DIR/resourced $ROOT_DIR/resourced-linux-$VERSION/


echo
echo "Compressing Darwin build"
tar -zcvf resourced-darwin-$VERSION.tar.gz -C $ROOT_DIR/resourced-darwin-$VERSION .

echo
echo "Compressing Linux build"
tar -zcvf resourced-linux-$VERSION.tar.gz -C $ROOT_DIR/resourced-linux-$VERSION .

rm -rf $ROOT_DIR/resourced-{darwin,linux}-$VERSION
