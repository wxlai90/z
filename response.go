package z

import "encoding/json"

func (z *Z) String(statusCode int, respStr string) {
	z.rw.WriteHeader(statusCode)
	z.rw.Write([]byte(respStr))
}

func (z *Z) JSON(statusCode int, respJSON any) {
	z.rw.WriteHeader(statusCode)
	z.rw.Header().Set("content-type", "application/json")
	json.NewEncoder(z.rw).Encode(respJSON)
}
