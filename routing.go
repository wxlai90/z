package z

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(z *Z)

func (app *App) handle(method string, path string, handler HandlerFunc, routeMiddlewares ...MiddlewareFunc) {
	finalHandler := handler

	for i := len(routeMiddlewares) - 1; i >= 0; i-- {
		finalHandler = routeMiddlewares[i](finalHandler)
	}

	for i := len(app.middlewares) - 1; i >= 0; i-- {
		finalHandler = app.middlewares[i](finalHandler)
	}

	app.mux.HandleFunc(fmt.Sprintf("%s %s", method, path), func(w http.ResponseWriter, r *http.Request) {
		zHandler := &Z{
			rw: w,
			r:  r,
		}
		finalHandler(zHandler)
	})
}

func (app *App) GET(path string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	app.handle(http.MethodGet, path, handler, middlewares...)
}

func (app *App) PUT(path string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	app.handle(http.MethodPut, path, handler, middlewares...)
}

func (app *App) POST(path string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	app.handle(http.MethodPost, path, handler, middlewares...)
}

func (app *App) PATCH(path string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	app.handle(http.MethodPatch, path, handler, middlewares...)
}

func (app *App) DELETE(path string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	app.handle(http.MethodDelete, path, handler, middlewares...)
}
