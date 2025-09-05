# z

A minimalistic wrapper around Go's `net/http` for Go > 1.22.

## Features

- Minimalistic and easy to use
- Built on top of `net/http`
- Supports path parameters (Go > 1.22)
- Middleware support
- Helper functions for common tasks

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

- `BindBody(reqBodyType any) error`: Binds the request body to a struct.
- `PathValue(key string) string`: Gets a path parameter by key.
- `Query(key string) string`: Gets a query parameter by key.
- `Header(key string) string`: Gets a request header by key.
- `Cookie(name string) (*http.Cookie, error)`: Gets a cookie by name.
- `FormFile(key string) (multipart.File, *multipart.FileHeader, error)`: Gets a file from a multipart form.

### Response Helpers

- `String(statusCode int, respStr string)`: Sends a string response.
- `JSON(statusCode int, respJSON any)`: Sends a JSON response.
- `Ok(body string)`: Sends a string response with a 200 status code.
- `OkJSON(data interface{})`: Sends a JSON response with a 200 status code.
- `SetHeader(key, value string)`: Sets a response header.
- `SetCookie(cookie *http.Cookie)`: Sets a cookie.
- `Error(err error, code int)`: Sends an error response.
- `Redirect(url string, code int)`: Redirects to a URL with the given status code.
- `ServeFile(filename string, forceDownload bool)`: Serves a file from disk; when `forceDownload` is true, sets `Content-Disposition` to trigger a download.

### Escape Hatches

When you need to break out of the z framework's abstractions and access the underlying Go `net/http` objects:

- `ResponseWriter() http.ResponseWriter`: Returns the underlying `http.ResponseWriter`.
- `Request() *http.Request`: Returns the underlying `*http.Request`.

#### Example Usage

```go
app.GET("/string", func(z *z.Z) {
	z.String(200, "plain text response")
})

app.GET("/json", func(z *z.Z) {
	z.JSON(200, map[string]string{"message": "hello"})
})

app.GET("/ok", func(z *z.Z) {
	z.Ok("Everything is OK!")
})

app.GET("/okjson", func(z *z.Z) {
	z.OkJSON(map[string]string{"status": "ok"})
})

app.GET("/header", func(z *z.Z) {
	z.SetHeader("X-Custom-Header", "value")
	z.Ok("Header set!")
})

app.GET("/cookie", func(z *z.Z) {
	z.SetCookie(&http.Cookie{Name: "token", Value: "abc123"})
	z.Ok("Cookie set!")
})

app.GET("/error", func(z *z.Z) {
	z.Error(fmt.Errorf("something went wrong"), 500)
})

app.GET("/redirect", func(z *z.Z) {
	z.Redirect("/new-location", http.StatusFound)
})

app.GET("/file", func(z *z.Z) {
	// true/false for force download
	z.ServeFile("/path/to/test.txt", true)
})
```

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
