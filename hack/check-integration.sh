#!/usr/bin/env bash

set -o errexit
set -o nounset

run_cmd="go test -race -p 1 -timeout 120s -tags=integration -covermode=atomic"

if [ $# -eq 0 ]; then
    $run_cmd $(go list ./... | grep -v vendor)
    exit
fi

pkg=$1
testname=$2

[[ $# -eq 3 ]] && args=$3 || args=""

echo "Running test pkg: $pkg name: $testname"
$run_cmd $pkg --run $testname $args
