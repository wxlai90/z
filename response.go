package z

import (
	"encoding/json"
	"net/http"
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
