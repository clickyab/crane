#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(readlink -f $(dirname ${BASH_SOURCE[0]}))

echo "export GOPATH=/home/develop/go" >> /home/develop/.zshrc
echo "export GOPATH=/home/develop/go" >> /etc/environment
echo "export PATH=$PATH:/usr/local/go/bin:/home/develop/go/bin" >> /home/develop/.zshrc

cd /home/develop/gad
make -f /home/develop/gad/Makefile mysql-setup
make -f /home/develop/gad/Makefile rabbitmq-setup

sudo -u develop /home/develop/gad/bin/provision_user.sh