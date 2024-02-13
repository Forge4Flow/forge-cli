#!/bin/sh

for f in forge-cli*; do shasum -a 256 $f > $f.sha256; done
