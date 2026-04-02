package handlers

import (
	"fmt"
	"net/http"

	"github.com/strjkc/noteapi/converter"
	"github.com/strjkc/noteapi/spellcheck"
	"github.com/strjkc/noteapi/state"
)

const (
	INTERNALERROR = "Internal Server Error"
	FILENOTFOUND  = "The Requested File Can Not Be Found"
)

type Handlers struct {
	State *state.State
}

func NewHandlers(state *state.State) *Handlers {
	h := Handlers{State: state}
	return &h
}

func (h *Handlers) HandleSpellCheck(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// here i should inject a file path to the NewChecker in order for the dict to be loaded with the correct language
	// also i should cache checkers per locale, or dicts per locale for reuse
	checker := spellcheck.NewChecker()
	parser := spellcheck.NewParser(checker)
	errors, err := parser.SpellCheckerService(r.Body)
	if err != nil {
		respondWithError(w, 500, INTERNALERROR)
	}
	fmt.Println(errors)
	json := marshalJson(errors)
	sendJson(w, 200, json)
}

func (h *Handlers) HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	// TODO
	// is should read from the body reader, and send the data to s3 or to a local directory
	// so i should have an interface called file uploader or something so i can plug in local dir or remote dir
	// after file is uploaded i can just respond with 201 and thats it really
}

func (h *Handlers) HandleGetHtml(w http.ResponseWriter, r *http.Request) {
	fileName := r.PathValue("fileName")
	// TODO:
	// Check if it's empty?
	// we can only have one storage type per instance of the app, so that should be loaded on startup
	// TODO:
	// does it already exist as html? if so does the version of the html match the version of the text file?
	// if yes, send
	// if no, convert again, sotre it to disk, update db and send back
	fileExists := h.State.Storage.FileExists(fileName)
	if !fileExists {
		respondWithError(w, 404, FILENOTFOUND)
		return
	}
	// There should be an err here?
	fileAsHtml := converter.ConvertToHtml(fileName)
	// strip the .txt
	fileNameAsHtml := fmt.Sprint("%s.html", fileName)
	h.State.Storage.StoreFile(fileAsHtml, fileNameAsHtml)
	// update the db with the file version
	http.ServeFile(w, r, string(fileAsHtml))

	// file name should be in the url path
	// first i check if such a file exists in the storage, is the storage local or s3?
	// if it does, we convert and respond back with the file
	// if not we just send 404
}
