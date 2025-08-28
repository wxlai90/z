package z

import (
	"log"
	"net/http"
)

type App struct {
	mux         *http.ServeMux
	middlewares []MiddlewareFunc
}

type Z struct {
	rw http.ResponseWriter
	r  *http.Request
}

func (app *App) Start(port string) {
	log.Printf("Running on %s\n", port)
	log.Fatalln(http.ListenAndServe(port, app.mux))
}

func (app *App) Use(middlewareFunc MiddlewareFunc) {
	app.middlewares = append(app.middlewares, middlewareFunc)
}

func New() *App {
	return &App{
		mux:         http.NewServeMux(),
		middlewares: []MiddlewareFunc{},
	}
}
