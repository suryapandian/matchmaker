package handlers

import (
	"encoding/json"
	"net/http"
)

const jsonContentType = "application/json; charset=utf-8"

func writeJSONMessage(msg string, code int, w http.ResponseWriter) {
	writeJSONStruct(struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{
		code,
		msg,
	}, code, w)
}

func writeJSONStruct(v interface{}, code int, w http.ResponseWriter) {
	d, err := json.Marshal(v)
	if err != nil {
		writeResponse([]byte("Unable to marshal data. Error: "+err.Error()), http.StatusInternalServerError, jsonContentType, w)
		return
	}
	writeResponse(d, code, jsonContentType, w)
}

func writeResponse(d []byte, code int, contentType string, w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(code)
	w.Write(d)
}
