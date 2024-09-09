#!/usr/bin/env sh
set -o errexit -o nounset

cd "$(dirname "$0")/.."

mkdir -p target

go build -o target/consul-ip-finder cmd/main.go 
