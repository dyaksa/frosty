package internal

import (
	"encoding/json"
	"net/http"
)

func responseError(resw http.ResponseWriter, respCode int, message string) {
	responseJson(resw, respCode, map[string]string{"error": message})
}

func responseJson(resw http.ResponseWriter, respCode int, payload interface{}) {
	resp, _ := json.Marshal(payload)
	resw.Header().Set("Content-Type", "application/json")
	resw.WriteHeader(respCode)
	resw.Write(resp)
}
