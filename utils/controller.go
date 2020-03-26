package utils

import (
	"encoding/json"
	"net/http"
)

// WriteJSONError format the error response and sends it
func WriteJSONError(w http.ResponseWriter, msg string) {
	payload := &HTTPError{
		Msg: msg,
	}

	answer, _ := json.Marshal(payload)
	w.Write(answer)
}

// WriteJSONSuccess format the response and sends it
func WriteJSONSuccess(w http.ResponseWriter, payload interface{}) {
	answer, _ := json.Marshal(payload)
	w.Write(answer)
}

// ReadBody get the request body and unmarshal its content
func ReadBody(r *http.Request) []byte {
	buff := make([]byte, r.ContentLength)
	r.Body.Read(buff)

	return buff
}
