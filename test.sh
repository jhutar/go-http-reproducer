#!/bin/sh

./server &

payloadSize=10000
threadsCount=100
iterationsCount=10
./client $payloadSize $threadsCount $iterationsCount
