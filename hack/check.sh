#!/bin/bash

set -o errexit
set -o nounset

run_cmd="go test -timeout 10s -race -tags=unit -covermode=atomic"

if [ $# -eq 0 ]; then
    for d in $(go list ./... | grep -v vendor); do
        $run_cmd $d
    done
    exit
fi

pkg=$1
test_case=$2

[[ $# -eq 3 ]] && args=$3 || args=""

echo "Running test pkg:" $pkg " name: " $test_case
$run_cmd $args $pkg --run $test_case
