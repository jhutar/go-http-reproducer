#!/bin/sh

./server &

payloadSize=${payloadSize:-10000}
threadsCount=${threadsCount:-100}
iterationsCount=${iterationsCount:-10}
./client $payloadSize $threadsCount $iterationsCount
