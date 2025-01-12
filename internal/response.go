package internal

import (
	"encoding/json"
	"log"
	"net/http"
)

func responseError(resw http.ResponseWriter, respCode int, message string) {

	log.Println("Response Error:", message)

	responseJson(resw, respCode, map[string]string{"error": message})
}

func responseJson(resw http.ResponseWriter, respCode int, payload interface{}) {
	resp, _ := json.Marshal(payload)

	log.Println("Response Json:", string(resp))

	resw.Header().Set("Content-Type", "application/json")
	resw.WriteHeader(respCode)
	resw.Write(resp)
}
