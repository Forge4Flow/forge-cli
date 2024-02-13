#!/bin/bash

export eTAG="latest-dev"
echo $1
if [ $1 ] ; then
  eTAG=$1
fi

echo Building forge4flow/forge-cli:$eTAG

docker build --build-arg http_proxy=$http_proxy --build-arg https_proxy=$https_proxy --target release -t forge4flow/forge-cli:$eTAG .

echo Building forge4flow/forge-cli:$eTAG-root

docker build --build-arg http_proxy=$http_proxy --build-arg https_proxy=$https_proxy --target root -t forge4flow/forge-cli:$eTAG-root .

if [ $? == 0 ] ; then

  docker create --name forge-cli forge4flow/forge-cli:$eTAG && \
  docker cp forge-cli:/usr/bin/forge-cli . && \
  docker rm -f forge-cli

else
 exit 1
fi
