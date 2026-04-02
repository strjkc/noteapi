package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func respondWithError(w http.ResponseWriter, status int, message string) {
	error := struct {
		Message string `json:"message"`
	}{
		Message: message,
	}
	json, err := json.Marshal(error)
	if err != nil {
		fmt.Println("An error occured marshaling json")
		return
	}

	sendJson(w, status, json)
}

func marshalJson(s any) []byte {
	json, err := json.Marshal(s)
	if err != nil {
		fmt.Println("An error occured marshaling json")
		return nil
	}
	return json
}

func sendJson(w http.ResponseWriter, status int, json []byte) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(json)
}
