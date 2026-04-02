package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/strjkc/noteapi/spellcheck"
)

const (
	INTERNALERROR = "Internal Server Error"
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

func HandleSpellCheck(w http.ResponseWriter, r *http.Request) {
	parser := spellcheck.NewParser()
	errors, err := parser.SpellCheckerService(r.Body)
	if err != nil {
		respondWithError(w, 500, INTERNALERROR)
	}
	fmt.Println(errors)
	json := marshalJson(errors)
	sendJson(w, 200, json)
}

func sendJson(w http.ResponseWriter, status int, json []byte) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(json)
}

func marshalJson(s any) []byte {
	json, err := json.Marshal(s)
	if err != nil {
		fmt.Println("An error occured marshaling json")
		return nil
	}
	return json
}
