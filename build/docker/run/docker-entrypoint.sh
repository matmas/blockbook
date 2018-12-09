#!/bin/dash
set -e

envsubst </opt/blockchaincfg.json.template >/home/blockbook/blockchaincfg.json

exec "$@"
