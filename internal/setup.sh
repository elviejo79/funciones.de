#!/bin/sh

CURRENT=$(pwd)
cd $(pwd)/../gopath
export GOPATH=$(pwd)
echo $GOPATH
cd $CURRENT

export PATH="$PATH:$GOPATH/bin"
