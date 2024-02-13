#!/bin/sh

./bin/forge-cli build # --squash=true

docker images |head -n 4
