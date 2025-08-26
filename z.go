package z

import "net/http"

type App struct {
	mux *http.ServeMux
}

func New() *App {
	return &App{
		mux: http.NewServeMux(),
	}
}
