#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

echo "time,threadsCount,iteration,error rate, average duration" >>results.csv

for threadsCountK in $( seq 1 20 ); do
    for iteration in $( seq 1 5 ); do
        echo "=== threadsCount:${threadsCountK}000 iteration:${iteration} ==="
        output="$( go run client.go 10000 "${threadsCountK}000" 1 2>&1 | tail -n 3 )"
        echo "$output"
        error_rate="$( echo "$output" | grep "Failure rate" | cut -d : -f 4 )"
        avg_duration="$( echo "$output" | grep "Successful average duration" | cut -d : -f 4 )"
        echo "\"$( date --utc -Ins )\",${threadsCountK}000,${iteration},${error_rate},${avg_duration}" | tee -a results.csv
    done
done
