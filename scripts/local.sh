#!/bin/bash
set -x
set -eo pipefail
echo it is from terminal
apt-get update && apt-get install build-essential -y
make -C /crane
