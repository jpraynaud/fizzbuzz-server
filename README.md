# FizzBuzz Server
This project implements a simple FizzBuzz REST server. 

It exposes 2 endpoints:
* **/render?limit=$limit&int1=$int1&int2=$int2&str1=$str1&str2=$str2** endpoint where **limit**, **int1** & **int2** are integer parameters and **str1** & **str2** are string parameters. When called, returns the FizzBuzz string associated with the parameters.
* **/statistics** endpoint. When called, returns the most called request parameters from previous endpoint and the number of hits of this request.

---

* [Algorithm](#algorithm)
* [Response](#response)
* [Examples](#examples)
* [Install](#install)
* [Build](#build)
* [Tests](#tests)
* [Help](#help)
* [Run](#run)
* [Documentation](#documentation)

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

## Install

With a [correctly configured](https://golang.org/doc/install#testing) Go toolchain (version 1.13+):

```sh
#Download sources
go get -u github.com/jpraynaud/fizzbuzz-server

# Go to sources directory
cd $GOPATH/src/github.com/jpraynaud/fizzbuzz-server
```

## Build

Build executable:

```sh
# Build
go build -o fizzbuzz-server cmd/server/main.go
```

## Tests

Run unit tests:

```sh
# Test
go test ./... -race -cover
```

## Help

Get help:

```sh
# Help with executable
./fizzbuzz-server --help

# or
go run cmd/server/main.go --help
```

## Run

The server accepts 2 configuration flags:
* **-address** is the address on which the server is listening (for example *0.0.0.0:8080*).
* **-environment** is the environment of the server (*development* or *production*).

If these flags are not set, they will default to environment variables:
* **SERVER_ADDR** 
* **SERVER_ENV** 

### Start server on 0.0.0.0:8080 in development:

```sh
# Start server 
./fizzbuzz-server -address=0.0.0.0:8080 -environment=development

# or 
SERVER_ADDR=0.0.0.0:8080 SERVER_ENV=development ./fizzbuzz-server

# or
go run cmd/server/main.go  -address=0.0.0.0:8080 -environment=development

# or
SERVER_ADDR=0.0.0.0:8080 SERVER_ENV=development go run cmd/server/main.go
```

### Then access endpoints:

```sh
# Renders FizzBuzz request
curl 'http://0.0.0.0:8080/render?limit=20&int1=4&int2=7&str1=AA&str2=BBB'

# Get Statistics
curl 'http://0.0.0.0:8080/statistics'
```

or:
* Render a request at [http://0.0.0.0:8080/render?limit=20&int1=4&int2=7&str1=AA&str2=BBB].
* Get statistics at [http://0.0.0.0:8080/statistics].

## Documentation

### Generation of the package documentation:
Generate documentation from source code and access it from [http://0.0.0.0:6060/pkg/github.com/jpraynaud/fizzbuzz-server/pkg/render/].

```sh
# Generate documentation
godoc
```

### Structure explaination:
The project is split in 2 packages:
* **main** package that creates the HTTP server and the router/handlers that serve the endpoints.
* **render** package with:
    * a **Request** that represents a FizzBuzz request (a struct that holds request parameters explained in [Algorithm](#algorithm)).
    * a **Response** that represents a response rendered from a **Request** (a struct that holds a channel of strings and an error).
    * a **Renderer** that processes a **Request** and returns a **Response**, while recording **Statistics**.
    * a **Statistics** that stores statistics (a struct that holds a map of total hits for requests and the top request so far).
    * a **RequestStatistic** that gives the statistic of a request (a struct that holds the **Request** and the total hits).


