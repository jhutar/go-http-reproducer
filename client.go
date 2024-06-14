package main

import (
	"crypto/tls"
	"crypto/x509"
	"bytes"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
	"sync"

	"golang.org/x/net/http2"
)

var (
	durationTotal     float64
	durationCounter   int
	durationTotalMutex sync.Mutex
	durationCounterMutex sync.Mutex
)

func doRequest(client *http.Client, method string, serverURL string, payload *bytes.Reader) (float64, error) {
	// Start time measurement
	startTime := time.Now()

	// Create a POST request with the payload
	req, err := http.NewRequest(method, serverURL, payload)
	if err != nil {
		return 0.0, err
	}

	// Send the request and get the response
	resp, err := client.Do(req)
	if err != nil {
		return 0.0, err
	}
	defer resp.Body.Close()

	// Read the response body (optional)
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0.0, err
	}

	// Calculate duration
	duration := time.Since(startTime)

	log.Printf("Request finished with status %d and took %v", resp.StatusCode, duration)
	return duration.Seconds(), nil
}

// Function to perform the request and update counters
func doRequestThread(iterationsCount int, tlsConfig *tls.Config, serverURL string, payloadSize int, wg *sync.WaitGroup) {
	defer wg.Done()

	// Configure transport to enable HTTP/2
	tr := &http2.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: tr}

	for i := 0; i < iterationsCount; i++ {
		doRequestOne(client, serverURL, payloadSize)
	}
}

// Function to perform one request
func doRequestOne(client *http.Client, serverURL string, payloadSize int) {
	// Create random payload
	payload := make([]byte, payloadSize)
	_, err := rand.Read(payload)
	if err != nil {
		panic(err)
	}

	duration, err := doRequest(client, http.MethodPost, serverURL, bytes.NewReader(payload))
	if err != nil {
		panic(err)
	}

	// Update counters within the goroutine
	durationTotalMutex.Lock()
	defer durationTotalMutex.Unlock()
	durationTotal += duration
	durationCounterMutex.Lock()
	defer durationCounterMutex.Unlock()
	durationCounter++
}

func main() {
	// Server address
	serverURL := "https://localhost:8000"

	// Payload size in bytes
	payloadSize := 10000

	// Number of requests to do
	requestsCount := 5000

	// Number of iterations
	iterationsCount := 1

	// Create a pool with the server certificate since it is not signed
	// by a known CA
	caCert, err := ioutil.ReadFile("server.crt")
	if err != nil {
		log.Fatalf("Reading server certificate: %s", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create TLS configuration with the certificate of the server
	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}

	var wg sync.WaitGroup

	for i := 1; i <= requestsCount; i++ {
		wg.Add(1)
		go doRequestThread(iterationsCount, tlsConfig, serverURL, payloadSize, &wg)
	}

	wg.Wait()

	log.Printf("Requests count: %d", durationCounter)
	log.Printf("Average duration: %f", durationTotal / float64(durationCounter))
}
