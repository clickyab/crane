#!/bin/bash
set -euo pipefail

SCRIPT_DIR=$(readlink -f $(dirname ${BASH_SOURCE[0]}))

exec ${SCRIPT_DIR}/ip2location