#!/bin/sh

export eTAG="latest-dev"
echo $1
if [ $1 ] ; then
  eTAG=$1
fi

# This file extracts the contents of the forge4flow/forge-cli:${eTAG} image
# to the bin folder.
# Run ./build.sh first to create the image itself with the various
# binaries within it.
docker create --name forge-cli forge4flow/forge-cli:${eTAG} && \
  mkdir -p ./bin && \
  docker cp forge-cli:/home/app/forge-cli ./bin && \
  docker cp forge-cli:/home/app/forge-cli-darwin ./bin && \
  docker cp forge-cli:/home/app/forge-cli-darwin-arm64 ./bin && \
  docker cp forge-cli:/home/app/forge-cli-armhf ./bin && \
  docker cp forge-cli:/home/app/forge-cli-arm64 ./bin && \
  docker cp forge-cli:/home/app/forge-cli.exe ./bin && \
  docker rm -f forge-cli
