package z

import (
	"encoding/json"
	"fmt"
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
