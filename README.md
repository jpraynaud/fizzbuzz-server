# FizzBuzz Server
[![Go Report Card](https://goreportcard.com/badge/github.com/jpraynaud/fizzbuzz-server)](https://goreportcard.com/report/github.com/jpraynaud/fizzbuzz-server)
[![GoDoc](https://godoc.org/github.com/jpraynaud/fizzbuzz-server/pkg/render?status.svg)](https://godoc.org/github.com/jpraynaud/fizzbuzz-server/pkg/render)

This project implements a simple FizzBuzz REST server. 

It exposes 2 endpoints:
* **/render?limit=$limit&int1=$int1&int2=$int2&str1=$str1&str2=$str2** GET endpoint where **limit**, **int1** & **int2** are integer parameters and **str1** & **str2** are string parameters. When called, returns the FizzBuzz string associated with the parameters.
* **/statistics** GET endpoint. When called, returns the most called request parameters from previous endpoint and the number of hits of this request.

---

* [Algorithm](#algorithm)
* [Response](#response)
* [Examples](#examples)
* [Docker](#docker)
* [Install](#install)
* [Build](#build)
* [Run](#run)
* [Tests](#tests)
* [Benchmarks](#benchmarks)
* [Help](#help)
* [Documentation](#documentation)
* [SSL](#ssl)

---

## Algorithm
The **FizzBuzz** algorithm takes **limit**, **int1**, **int2**, **str1** and **str2** as parameters and returns a list of strings with numbers from **1** to **limit** where: 
* all multiples of **int1** are replaced by **str1**.
* all multiples of **int2** are replaced by **str2**.
* all multiples of **int1** and **int2** are replaced by **str1str2**.

### Example 1
* Parameters: **limit**=20, **int1**=3, **int2**=5, **str1**=A, **str2**=B
* Result: *1,2,A,4,B,A,7,8,A,B,11,A,13,14,AB,16,17,A,19,B*

### Example 2
* Parameters: **limit**=30, **int1**=4, **int2**=7, **str1**=AA, **str2**=BBB
* Result: *1,2,3,AA,5,6,BBB,AA,9,10,11,AA,13,BBB,15,AA,17,18,19,AA,BBB,22,23,AA,25,26,27,AABBB,29,30*

## Response
The response is sent in JSON format with 2 fields:
* **error**: a boolean, *true* if an error occurred else *false*.
* **response**: an object that will be:
    * a string for /render endpoint.
    * a nested object for /statistics endpoint.

## Examples
### Example: /render?limit=20&int1=4&int2=7&str1=AA&str2=BBB
**response** returns the FizzBuzz list.
```
{
    "error": false,
    "response": "1,2,3,AA,5,6,BBB,AA,9,10,11,AA,13,BBB,15,AA,17,18,19,AA"
}
```

### Example: /render?limit=Z&int1=4&int2=7&str1=AA&str2=BBB
**response** returns an error message.
```
{
    "error": true,
    "response": "limit parameter must be an integer, value Z was given"
}
```

### Example: /statistics
**response** returns an object with 2 fields:
* **total** is the number of hits for the top request.
* **request** is the top request.
```
{
    "error": false,
    "response": {
        "request": {
            "int1": 4,
            "int2": 7,
            "limit": 20,
            "str1": "AA",
            "str2": "BBB"
        },
        "total": 7
    }
}
```

## Docker

### Build and run Docker container:

```sh
# Git clone
git clone https://github.com/jpraynaud/fizzbuzz-server

# Build Docker container
docker build -t jpraynaud/fizzbuzz-server fizzbuzz-server

# Run Docker container in production
docker run --rm -p 8080:8080 -e SERVER_ADDR='0.0.0.0:8080' -e SERVER_ENV='production' --name='fizzbuzz' -d jpraynaud/fizzbuzz-server

# Show Docker container logs
docker logs -f fizzbuzz

# Kill Docker container
docker kill fizzbuzz

```

### Then access endpoints:

```sh
# Renders FizzBuzz request
curl 'http://0.0.0.0:8080/render?limit=100&int1=3&int2=5&str1=fizz&str2=buzz'

# Get statistics
curl 'http://0.0.0.0:8080/statistics'
```

or:

* Render a request at [http://0.0.0.0:8080/render?limit=100&int1=3&int2=5&str1=fizz&str2=buzz].
* Get statistics at [http://0.0.0.0:8080/statistics].

## Install

With a [correctly configured](https://golang.org/doc/install#testing) Go toolchain (version 1.13+):

```sh
#Download sources
go get -d -t -v github.com/jpraynaud/fizzbuzz-server/cmd/server

# Go to sources directory
cd $GOPATH/src/github.com/jpraynaud/fizzbuzz-server
```

## Build

Build executable:

```sh
# Build
go build -v -o fizzbuzz-server cmd/server/main.go
```

## Run

The server accepts these configuration flags:
* **-address** is the address on which the server is listening (for example *0.0.0.0:8080*).
* **-environment** is the environment of the server (*development* or *production*).
* **-tlscert** is the path of the SSL certificate file.
* **-tlskey** is the path of the SSL private key file.


If these flags are not set, they will respectively default to environment variables:
* **SERVER_ADDR** 
* **SERVER_ENV** 
* **SERVER_TLSCERTFILE**
* **SERVER_TLSKEYFILE**

***If the certificate/private key files are not specified the server will start without TLS.***

### Start server on 0.0.0.0:8080 in development:

```sh
# Start server 
./fizzbuzz-server -address=0.0.0.0:8080 -environment=development

# or 
SERVER_ADDR=0.0.0.0:8080 SERVER_ENV=development ./fizzbuzz-server

# or
go run cmd/server/main.go -address=0.0.0.0:8080 -environment=development

# or
SERVER_ADDR=0.0.0.0:8080 SERVER_ENV=development go run cmd/server/main.go
```

### Start server on 0.0.0.0:8080 in production:

```sh
# Start server 
./fizzbuzz-server -address=0.0.0.0:8080 -environment=production

# or 
SERVER_ADDR=0.0.0.0:8080 SERVER_ENV=production ./fizzbuzz-server

# or
go run cmd/server/main.go -address=0.0.0.0:8080 -environment=production

# or
SERVER_ADDR=0.0.0.0:8080 SERVER_ENV=production go run cmd/server/main.go
```

### Then access endpoints:

```sh
# Renders FizzBuzz request
curl 'http://0.0.0.0:8080/render?limit=100&int1=3&int2=5&str1=fizz&str2=buzz'

# Get statistics
curl 'http://0.0.0.0:8080/statistics'
```

or:

* Render a request at [http://0.0.0.0:8080/render?limit=100&int1=3&int2=5&str1=fizz&str2=buzz].
* Get statistics at [http://0.0.0.0:8080/statistics].

## Tests

Run unit tests:

```sh
# Test with race detection and code coverage
go test -race -cover -v ./...
```

## Benchmarks

Run benchmarks:

```sh
# Benchmark
go test -run="^$" -bench=. ./...
```

Or run server benchmark:
```sh
# Benchmark /render
ab -n 100000 -c 100 -k "http://0.0.0.0:8080/render?limit=100&int1=4&int2=7&str1=AA&str2=BBB"

# and /statistics
ab -n 100000 -c 100 -k "http://0.0.0.0:8080/statistics"
```

## Help

Get help:

```sh
# Help with executable
./fizzbuzz-server --help

# or
go run cmd/server/main.go --help
```

## Documentation

[![GoDoc](https://godoc.org/github.com/jpraynaud/fizzbuzz-server/pkg/render?status.svg)](https://godoc.org/github.com/jpraynaud/fizzbuzz-server/pkg/render)

### Generation of the package documentation:
Generate documentation from source code and access it from [http://0.0.0.0:6060/pkg/github.com/jpraynaud/fizzbuzz-server/].

```sh
# Generate documentation
godoc
```

### Structure explanation:
The project is split in 2 packages:
* **main** package that creates the HTTP server and the router/handlers that serve the endpoints.
* **render** package with:
    * a **Request** that represents a FizzBuzz request (a struct that holds request parameters explained in [Algorithm](#algorithm)).
    * a **Response** that represents a response rendered from a **Request** (a struct that holds a channel of strings and an error).
    * a **Renderer** that processes a **Request** and returns a **Response**, while recording **Statistics**.
    * a **Statistics** that stores statistics (a struct that holds a map of total hits for requests and the top request so far).
    * a **RequestStatistic** that gives the statistic of a request (a struct that holds the **Request** and the total hits).

## SSL

### Self signed SSL certificate generation:

```sh
# Generate Certificate and Private Key
openssl genrsa -out server.key 2048
openssl rsa -in server.key -out server.key
openssl req -sha256 -new -key server.key -out server.csr -subj '/CN=localhost'
openssl x509 -req -sha256 -days 365 -in server.csr -signkey server.key -out server.crt
cat server.crt server.key > cert.pem
```

### Start SSL server on 0.0.0.0:8080 in production:

```sh
# Start server 
./fizzbuzz-server -address=0.0.0.0:8080 -environment=production -tlscert=cert.pem -tlskey=server.key

# or 
SERVER_ADDR=0.0.0.0:8080 SERVER_ENV=production SERVER_TLSCERTFILE=cert.pem SERVER_TLSKEYFILE=server.key ./fizzbuzz-server

# or
go run cmd/server/main.go -address=0.0.0.0:8080 -environment=production -tlscert=cert.pem -tlskey=server.key

# or
SERVER_ADDR=0.0.0.0:8080 SERVER_ENV=production SERVER_TLSCERTFILE=cert.pem SERVER_TLSKEYFILE=server.key go run cmd/server/main.go
```

### Then access endpoints:

```sh
# Renders FizzBuzz request
curl --insecure 'https://localhost:8080/render?limit=100&int1=3&int2=5&str1=fizz&str2=buzz'

# Get statistics
curl --insecure 'https://localhost:8080/statistics'
```





