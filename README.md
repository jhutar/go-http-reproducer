HTTP/2 server & client experiment
=================================

Taken from <https://posener.github.io/http2/> and <https://stackoverflow.com/questions/64814173/how-do-i-use-sans-with-openssl-instead-of-common-name>.

Run server:

    $ go run server.go
    2024/06/11 15:41:58 Serving on https://0.0.0.0:8000
    [...]

And client:

    $ go run client.go
    2024/06/11 16:11:14 Request finished with status 200 and took 2.038082ms
    2024/06/11 16:11:14 Request finished with status 200 and took 136.928µs
    2024/06/11 16:11:14 Request finished with status 200 and took 119.144µs
    2024/06/11 16:11:14 Request finished with status 200 and took 100.929µs
    2024/06/11 16:11:14 Request finished with status 200 and took 144.858µs
    2024/06/11 16:11:14 Request finished with status 200 and took 74.27µs
    2024/06/11 16:11:14 Request finished with status 200 and took 97.951µs
    2024/06/11 16:11:14 Request finished with status 200 and took 105.434µs
    2024/06/11 16:11:14 Request finished with status 200 and took 69.572µs
    2024/06/11 16:11:14 Request finished with status 200 and took 100.601µs
    2024/06/11 16:11:14 Average duration: 0.000299

To run the test in UBI9 container:

    $ podman build -f Containerfile.1.19 . -t go-http-reproducer:1-19   # UBI9 base
    $ podman build -f Containerfile.1.21 . -t go-http-reproducer:1-21   # UBI9 based
    $ podman build -f Containerfile.1.22 . -t go-http-reproducer:1-22   # Alpine base
    $ podman run -ti --rm go-http-reproducer:1-19
    $ podman run -ti --rm go-http-reproducer:1-21
    $ podman run -ti --rm go-http-reproducer:1-22

Because HTTP/2 needs TLS, we generated certificate with (cert is part of the repo):

    openssl req -newkey rsa:2048 -nodes -keyout server.key -x509 -days 365 -out server.crt -addext "subjectAltName = DNS:localhost"
