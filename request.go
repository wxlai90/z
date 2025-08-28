package z

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
)

func (z *Z) BindBody(reqBodyType any) error {
	if z.r.Body == nil {
		return fmt.Errorf("request body is nil")
	}

	defer z.r.Body.Close()
	return json.NewDecoder(z.r.Body).Decode(reqBodyType)
}

func (z *Z) PathValue(key string) string {
	return z.r.PathValue(key)
}

func (z *Z) Query(key string) string {
	return z.r.URL.Query().Get(key)
}

func (z *Z) Header(key string) string {
	return z.r.Header.Get(key)
}

func (z *Z) Cookie(name string) (*http.Cookie, error) {
	return z.r.Cookie(name)
}

func (z *Z) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return z.r.FormFile(key)
}
