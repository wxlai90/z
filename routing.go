package z

import (
	"fmt"
	"net/http"
)

func (app *App) GET(path string, handler func(z *Z)) {
	app.mux.HandleFunc(fmt.Sprintf("%s %s", http.MethodGet, path), func(w http.ResponseWriter, r *http.Request) {
		zHandler := &Z{
			rw: w,
			r:  r,
		}

		handler(zHandler)
	})
}

func (app *App) PUT(path string, handler func(z *Z)) {
	app.mux.HandleFunc(fmt.Sprintf("%s %s", http.MethodPut, path), func(w http.ResponseWriter, r *http.Request) {
		zHandler := &Z{
			rw: w,
			r:  r,
		}

		handler(zHandler)
	})
}

func (app *App) POST(path string, handler func(z *Z)) {
	app.mux.HandleFunc(fmt.Sprintf("%s %s", http.MethodPost, path), func(w http.ResponseWriter, r *http.Request) {
		zHandler := &Z{
			rw: w,
			r:  r,
		}

		handler(zHandler)
	})
}

func (app *App) PATCH(path string, handler func(z *Z)) {
	app.mux.HandleFunc(fmt.Sprintf("%s %s", http.MethodPatch, path), func(w http.ResponseWriter, r *http.Request) {
		zHandler := &Z{
			rw: w,
			r:  r,
		}

		handler(zHandler)
	})
}

func (app *App) DELETE(path string, handler func(z *Z)) {
	app.mux.HandleFunc(fmt.Sprintf("%s %s", http.MethodDelete, path), func(w http.ResponseWriter, r *http.Request) {
		zHandler := &Z{
			rw: w,
			r:  r,
		}

		handler(zHandler)
	})
}
