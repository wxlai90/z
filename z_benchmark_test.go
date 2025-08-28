package z

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkGet(b *testing.B) {
	app := New()
	app.GET("/", func(z *Z) {
		z.Ok("ok")
	})
	req, _ := http.NewRequest("GET", "/", nil)
	rw := httptest.NewRecorder()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.ServeHTTP(rw, req)
	}
}

func BenchmarkGetWithSingleMiddleware(b *testing.B) {
	app := New()
	app.Use(func(next HandlerFunc) HandlerFunc {
		return func(z *Z) {
			next(z)
		}
	})
	app.GET("/", func(z *Z) {
		z.Ok("ok")
	})
	req, _ := http.NewRequest("GET", "/", nil)
	rw := httptest.NewRecorder()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.ServeHTTP(rw, req)
	}
}

func BenchmarkGetWithMultipleMiddlewares(b *testing.B) {
	app := New()
	app.Use(func(next HandlerFunc) HandlerFunc {
		return func(z *Z) {
			next(z)
		}
	})
	app.Use(func(next HandlerFunc) HandlerFunc {
		return func(z *Z) {
			next(z)
		}
	})
	app.GET("/", func(z *Z) {
		z.Ok("ok")
	})
	req, _ := http.NewRequest("GET", "/", nil)
	rw := httptest.NewRecorder()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.ServeHTTP(rw, req)
	}
}

func BenchmarkPostWithJSONBinding(b *testing.B) {
	app := New()
	app.POST("/", func(z *Z) {
		var p struct{}
		z.BindBody(&p)
		z.Ok("ok")
	})
	body := `{"name":"test"}`
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(body))
	rw := httptest.NewRecorder()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.ServeHTTP(rw, req)
	}
}
