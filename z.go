package z

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type App struct {
	mux *http.ServeMux
}

type Z struct {
	rw http.ResponseWriter
	r  *http.Request
}

func (z *Z) String(statusCode int, respStr string) {
	z.rw.WriteHeader(statusCode)
	z.rw.Write([]byte(respStr))
}

func (z *Z) JSON(statusCode int, respJSON any) {
	z.rw.WriteHeader(statusCode)
	z.rw.Header().Set("content-type", "application/json")
	json.NewEncoder(z.rw).Encode(respJSON)
}

func (z *Z) BindBody(reqBodyType any) error {
	if z.r.Body == nil {
		return fmt.Errorf("request body is nil")
	}

	defer z.r.Body.Close()
	return json.NewDecoder(z.r.Body).Decode(reqBodyType)
}

func (app *App) Start(port string) {
	log.Printf("Running on %s\n", port)
	log.Fatalln(http.ListenAndServe(port, app.mux))
}

func New() *App {
	return &App{
		mux: http.NewServeMux(),
	}
}
