#!/bin/bash
set -euo pipefail

SCRIPT_DIR=$(readlink -f $(dirname ${BASH_SOURCE[0]}))

${SCRIPT_DIR}/migration -action=up -app=crane
exec ${SCRIPT_DIR}/octopus