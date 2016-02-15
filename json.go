package rest

import (
	"encoding/json"
	"net/http"
)

// EncoderJSON implementes the Encoder interface, encoding the request source
// as JSON and adding the application/json Content-Type header, and decoding
// the response body.
type EncoderJSON struct{}

func (e EncoderJSON) Encode(req *http.Request, src interface{}) ([]byte, error) {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	return json.Marshal(src)
}

func (e EncoderJSON) Decode(res *http.Response, dst interface{}) error {
	return json.NewDecoder(res.Body).Decode(dst)
}
