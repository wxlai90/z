package z

import (
	"net/http"
	"testing"
)

func TestNewApp(t *testing.T) {
	app := New()
	if app == nil {
		t.Fatal("New() returned nil")
	}
	if app.mux == nil {
		t.Error("App.mux is nil")
	}
	if len(app.middlewares) != 0 {
		t.Error("App.middlewares should be empty on init")
	}
}

func TestUseMiddleware(t *testing.T) {
	app := New()
	mw := func(next HandlerFunc) HandlerFunc {
		return next
	}
	app.Use(mw)
	if len(app.middlewares) != 1 {
		t.Error("Middleware not added")
	}
}

func TestServeHTTP(t *testing.T) {
	app := New()
	called := false
	app.mux.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})
	req, _ := http.NewRequest("GET", "/test", nil)
	rw := &mockResponseWriter{}
	app.ServeHTTP(rw, req)
	if !called {
		t.Error("Handler was not called by ServeHTTP")
	}
}

func TestMiddlewareChaining(t *testing.T) {
	app := New()
	var chain []string
	mw1 := func(next HandlerFunc) HandlerFunc {
		return func(z *Z) {
			chain = append(chain, "mw1")
			if next != nil {
				next(z)
			}
		}
	}
	mw2 := func(next HandlerFunc) HandlerFunc {
		return func(z *Z) {
			chain = append(chain, "mw2")
			if next != nil {
				next(z)
			}
		}
	}
	app.Use(mw1)
	app.Use(mw2)
	app.mux.HandleFunc("GET /chain", func(w http.ResponseWriter, r *http.Request) {
		z := &Z{rw: w, r: r}
		final := mw2(mw1(func(z *Z) {
			chain = append(chain, "handler")
		}))
		final(z)
	})
	req, _ := http.NewRequest("GET", "/chain", nil)
	rw := &mockResponseWriter{}
	app.ServeHTTP(rw, req)
	if len(chain) != 3 || chain[0] != "mw2" || chain[1] != "mw1" || chain[2] != "handler" {
		t.Errorf("Middleware chain incorrect: %v", chain)
	}
}
