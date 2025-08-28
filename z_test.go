package z

import (
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
