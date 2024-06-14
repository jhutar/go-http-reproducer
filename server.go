package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	// Create a server on port 8000
	// Exactly how you would run an HTTP/1.1 server
	http2Server := http2.Server{
		MaxReadFrameSize: 16384,
	}

	srv := &http.Server{Addr: ":8000", Handler: h2c.NewHandler(http.HandlerFunc(handler), &http2Server)}

	// Start the server with TLS, since we are running HTTP/2 it must be
	// run with TLS.
	// Exactly how you would run an HTTP/1.1 server with TLS connection.
	log.Printf("Serving on https://0.0.0.0:8000")
	log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
}

func handler(w http.ResponseWriter, r *http.Request) {
	count := 0
	size := 0
	readStartTime := time.Now()
	defer func() {
		readEndTime := time.Now()
		duration := readEndTime.Sub(readStartTime)
		if duration.Seconds() > 10 {
			fmt.Println(fmt.Sprintf("GGMGGM17 HandleStream read loop count %d time %s ts %d.%09d", count, duration.String(), readEndTime.Unix(), readEndTime.Nanosecond()))
		}
	}()

	// TODO: minimize garbage, optimize recvBuffer code/ownership
	const readSize = 8196
	var bodyBuffer bytes.Buffer
	for buf := make([]byte, readSize); ; {
		count++
		httpReadStart := time.Now()
		n, err := r.Body.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Request size:", size)
				return
			} else {
				panic(fmt.Errorf("error while reading body %s", err.Error()))
			}
		}
		httpReadEnd := time.Now()
		httpReadDuration := httpReadEnd.Sub(httpReadStart)
		if httpReadDuration.Seconds() > 3 {
			fmt.Println(fmt.Sprintf("GGMGGM18 HandleStream read loop single iter count %d time %s ts %d.%09d req obj %#v", count, httpReadDuration.String(), httpReadEnd.Unix(), httpReadEnd.Nanosecond(), r))
		}
		size += n

		if n > 0 {
			buf = buf[n:]
			_, err := bodyBuffer.Write(buf)

			if err != nil {
				panic(err)
			}
		}

		allocBufStart := time.Now()
		if len(buf) == 0 {
			buf = make([]byte, readSize)
		}
		allocBufDuration := time.Now().Sub(allocBufStart)
		if allocBufDuration.Seconds() > 3 {
			fmt.Println(fmt.Sprintf("GGMGGM19 HandleStream read loop make buf single iter count %d time %s ts %d.%09d", count, httpReadDuration.String()))
		}
	}
}
