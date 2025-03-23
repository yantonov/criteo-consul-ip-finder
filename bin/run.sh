#!/usr/bin/env sh
set -o errexit -o nounset

cd "$(dirname "$0")/.."

./target/consul-ip-finder $@
