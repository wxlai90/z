package z

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
)

func (z *Z) String(statusCode int, respStr string) {
	z.rw.WriteHeader(statusCode)
	z.rw.Write([]byte(respStr))
}

func (z *Z) JSON(statusCode int, respJSON any) {
	z.rw.Header().Set("content-type", "application/json")
	z.rw.WriteHeader(statusCode)
	json.NewEncoder(z.rw).Encode(respJSON)
}

func (z *Z) Ok(body string) {
	z.String(http.StatusOK, body)
}

func (z *Z) OkJSON(data interface{}) {
	z.JSON(http.StatusOK, data)
}

func (z *Z) SetHeader(key, value string) {
	z.rw.Header().Set(key, value)
}

func (z *Z) SetCookie(cookie *http.Cookie) {
	http.SetCookie(z.rw, cookie)
}

func (z *Z) Error(err error, code int) {
	http.Error(z.rw, err.Error(), code)
}

func (z *Z) Redirect(url string, code int) {
	http.Redirect(z.rw, z.r, url, code)
}

func (z *Z) ServeFile(filename string, forceDownload bool) {
	if forceDownload {
		name := filepath.Base(filename)
		z.rw.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", name))
	}

	http.ServeFile(z.rw, z.r, filename)
}

func (z *Z) ServeDir(path string) {
	z.GET(path, http.FileServer(http.Dir(path)))
}
