package z

import (
	"net/http"
	"testing"
)

func TestAppHandleRegistersRoute(t *testing.T) {
	app := New()
	called := false
	handler := func(z *Z) { called = true }
	app.GET("/test", handler)

	req, _ := http.NewRequest("GET", "/test", nil)
	rw := &mockResponseWriter{}
	app.mux.ServeHTTP(rw, req)
	if !called {
		t.Error("Handler was not called for GET route")
	}
}

func TestAllHTTPMethods(t *testing.T) {
	app := New()
	methods := []struct {
		name     string
		method   string
		register func(string, HandlerFunc)
	}{
		{"GET", http.MethodGet, func(path string, handler HandlerFunc) { app.GET(path, handler) }},
		{"PUT", http.MethodPut, func(path string, handler HandlerFunc) { app.PUT(path, handler) }},
		{"POST", http.MethodPost, func(path string, handler HandlerFunc) { app.POST(path, handler) }},
		{"PATCH", http.MethodPatch, func(path string, handler HandlerFunc) { app.PATCH(path, handler) }},
		{"DELETE", http.MethodDelete, func(path string, handler HandlerFunc) { app.DELETE(path, handler) }},
	}

	for _, m := range methods {
		called := false
		handler := func(z *Z) { called = true }
		m.register("/route", handler)
		req, _ := http.NewRequest(m.method, "/route", nil)
		rw := &mockResponseWriter{}
		app.mux.ServeHTTP(rw, req)
		if !called {
			t.Errorf("Handler not called for %s", m.name)
		}
	}
}

func TestRouteAndAppMiddlewares(t *testing.T) {
	app := New()
	order := []string{}

	app.Use(func(next HandlerFunc) HandlerFunc {
		return func(z *Z) {
			order = append(order, "app")
			next(z)
		}
	})

	routeMw := func(next HandlerFunc) HandlerFunc {
		return func(z *Z) {
			order = append(order, "route")
			next(z)
		}
	}

	handler := func(z *Z) { order = append(order, "handler") }

	app.GET("/mwtest", handler, routeMw)

	req, _ := http.NewRequest("GET", "/mwtest", nil)
	rw := &mockResponseWriter{}
	app.mux.ServeHTTP(rw, req)

	expected := []string{"app", "route", "handler"}
	for i, v := range expected {
		if order[i] != v {
			t.Errorf("Expected %s at position %d, got %s", v, i, order[i])
		}
	}
}

type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() http.Header       { return http.Header{} }
func (m *mockResponseWriter) Write([]byte) (int, error) { return 0, nil }
func (m *mockResponseWriter) WriteHeader(int)           {}
