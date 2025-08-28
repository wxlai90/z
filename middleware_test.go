package z

import "testing"

func TestMiddlewareFunc(t *testing.T) {
	called := false
	handler := func(z *Z) { called = true }
	mw := func(next HandlerFunc) HandlerFunc {
		return func(z *Z) {
			called = true
			next(z)
		}
	}
	wrapped := mw(handler)
	wrapped(nil)
	if !called {
		t.Error("MiddlewareFunc did not call handler")
	}
}
