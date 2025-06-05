package rest

import (
	"encoding/json"
	"net/http"
)

type ResponsePayload struct {
	Result interface{} `json:"result"`
	Errors []string    `json:"errors"`
}

func (ur *ResponsePayload) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (r *ResponsePayload) MarshalJSON() ([]byte, error) {
	type Alias ResponsePayload
	cp := Alias(*r)

	return json.Marshal(cp)
}
