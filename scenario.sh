#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

###echo "time,threadsCount,iteration,error rate,average duration" >>results.csv
###
###for threadsCountK in $( seq 1 20 ); do
###    for iteration in $( seq 1 5 ); do
###        echo "=== threadsCount:${threadsCountK}000 iteration:${iteration} ==="
###        output="$( go run client.go 10000 "${threadsCountK}000" 1 2>&1 | tail -n 3 )"
###        echo "$output"
###        error_rate="$( echo "$output" | grep "Failure rate" | cut -d : -f 4 )"
###        avg_duration="$( echo "$output" | grep "Successful average duration" | cut -d : -f 4 )"
###        echo "\"$( date --utc -Ins )\",${threadsCountK}000,${iteration},${error_rate},${avg_duration}" | tee -a results.csv
###    done
###done

podman build -f Containerfile.1.19 . -t go-http-reproducer:1-19
podman build -f Containerfile.1.21 . -t go-http-reproducer:1-21
podman build -f Containerfile.1.22 . -t go-http-reproducer:1-22

echo '"time","image","threadsCount","iteration","error rate","average duration"' >>results-container.csv

for image in go-http-reproducer:1-19 go-http-reproducer:1-21 go-http-reproducer:1-22; do
    for threadsCountK in $( seq 1 20 ); do
        for iteration in $( seq 1 3 ); do
            echo "=== threadsCount:${threadsCountK}000 iteration:${iteration} ==="
            output="$( podman run -e payloadSize=10000 -e threadsCount=${threadsCountK}000 -e iterationsCount=10 -ti --rm "$image" 2>&1 | tail -n 3)"
            echo "$output"
            error_rate="$( echo "$output" | grep "Failure rate" | cut -d : -f 4 | sed "s/\s*\([0-9.]\+\)\s*/\1/" )"
            avg_duration="$( echo "$output" | grep "Successful average duration" | cut -d : -f 4 | sed "s/\s*\([0-9.]\+\)\s*/\1/" )"
            echo "\"$( date --utc -Ins )\",\"${image}\",${threadsCountK}000,${iteration},${error_rate},${avg_duration}" | tee -a results-container.csv
        done
    done
done
