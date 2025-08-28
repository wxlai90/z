# z

A minimalistic wrapper around Go's `net/http` for Go > 1.22.

## Features

*   Minimalistic and easy to use
*   Built on top of `net/http`
*   Supports path parameters (Go > 1.22)
*   Middleware support
*   Helper functions for common tasks

## Installation

```bash
go get github.com/wxlai90/z
```

## Usage

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/wxlai90/z"
)

func main() {
	app := z.New()

	app.Use(func(next z.HandlerFunc) z.HandlerFunc {
		return func(z *z.Z) {
			fmt.Println("Middleware 1")
			next(z)
		}
	})

	app.GET("/", func(z *z.Z) {
		z.Ok("Hello, World!")
	})

	app.GET("/users/{id}", func(z *z.Z) {
		id := z.PathValue("id")
		z.Ok(fmt.Sprintf("User ID: %s", id))
	})

	log.Fatal(http.ListenAndServe(":8080", app))
}
```

## API

### Request Helpers

*   `BindBody(reqBodyType any) error`: Binds the request body to a struct.
*   `PathValue(key string) string`: Gets a path parameter by key.
*   `Query(key string) string`: Gets a query parameter by key.
*   `Header(key string) string`: Gets a request header by key.
*   `Cookie(name string) (*http.Cookie, error)`: Gets a cookie by name.
*   `FormFile(key string) (multipart.File, *multipart.FileHeader, error)`: Gets a file from a multipart form.

### Response Helpers

*   `String(statusCode int, respStr string)`: Sends a string response.
*   `JSON(statusCode int, respJSON any)`: Sends a JSON response.
*   `Ok(body string)`: Sends a string response with a 200 status code.
*   `OkJSON(data interface{})`: Sends a JSON response with a 200 status code.
*   `SetHeader(key, value string)`: Sets a response header.
*   `SetCookie(cookie *http.Cookie)`: Sets a cookie.
*   `Error(err error, code int)`: Sends an error response.

## Test Results

```
$ go test ./... -cover
ok  	github.com/wxlai90/z	0.680s	coverage: 100.0% of statements
```

## Benchmarks

```
$ go test ./... -bench=.
goos: darwin
goarch: amd64
pkg: github.com/wxlai90/z
cpu: Intel(R) Core(TM) i5-8210Y CPU @ 1.60GHz
BenchmarkGet-4                          	 4132106	       299.1 ns/op
BenchmarkGetWithSingleMiddleware-4      	 4091242	       297.7 ns/op
BenchmarkGetWithMultipleMiddlewares-4   	 3714985	       308.7 ns/op
BenchmarkPostWithJSONBinding-4          	 1665776	       692.5 ns/op
PASS
ok  	github.com/wxlai90/z	6.983s
```
