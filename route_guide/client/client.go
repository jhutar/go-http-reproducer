/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a simple gRPC client that demonstrates how to use gRPC-Go libraries
// to perform unary, client streaming, server streaming and full duplex RPCs.
//
// It interacts with the route guide service whose definition can be found in routeguide/route_guide.proto.
package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"time"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/examples/data"
	pb "google.golang.org/grpc/examples/route_guide/routeguide"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.example.com", "The server name used to verify the hostname returned by the TLS handshake")
	threadCount        = flag.Int("thread_count", 1, "Number of concurrent threads")
	iterationCount     = flag.Int("iteration_count", 10, "Number of requests mad by each thread")
)

var (
	durationTotal     float64
	durationCounter   int
	failuresCounter   int
	durationTotalMutex sync.Mutex
	durationCounterMutex sync.Mutex
	failuresCounterMutex sync.Mutex
)

// printFeature gets the feature for the given point.
func printFeature(client pb.RouteGuideClient, point *pb.Point) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	before := time.Now()
	_, err := client.GetFeature(ctx, point)
	duration := time.Since(before).Seconds()

	if err == nil {
		durationTotalMutex.Lock()
		defer durationTotalMutex.Unlock()
		durationTotal += duration
		durationCounterMutex.Lock()
		defer durationCounterMutex.Unlock()
		durationCounter++
	} else {
		failuresCounterMutex.Lock()
		defer failuresCounterMutex.Unlock()
		failuresCounter++
	}
}

func randomPoint(r *rand.Rand) *pb.Point {
	lat := (r.Int31n(180) - 90) * 1e7
	long := (r.Int31n(360) - 180) * 1e7
	return &pb.Point{Latitude: lat, Longitude: long}
}

func myThread(wg *sync.WaitGroup) {
	defer wg.Done()

	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile = data.Path("x509/ca_cert.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials: %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.NewClient(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewRouteGuideClient(conn)

	for i := 0; i < *iterationCount; i++ {
		printFeature(client, &pb.Point{Latitude: 409146138, Longitude: -746188906})
	}
}

func main() {
	flag.Parse()

	var wg sync.WaitGroup

	start := time.Now()

	for i := 0; i < *threadCount; i++ {
		wg.Add(1)
		go myThread(&wg)
	}

	wg.Wait()

	duration := time.Since(start)

	log.Printf("Average duration: %f / %d = %fs", durationTotal, durationCounter, durationTotal / float64(durationCounter))
	log.Printf("Runtime: %v", duration)
}
