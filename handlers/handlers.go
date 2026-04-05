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
	FILENOTSAVED  = "The File Was Not Stored Due to an Internal Error"
	BADREQ        = "Bad Request"
)

type Handlers struct {
	State *state.State
}

func NewHandlers(state *state.State) *Handlers {
	h := Handlers{State: state}
	return &h
}

// TODO: i should not mix html and json apis, apis that return json are under /api/ apis that return html are under something else
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
	err := h.State.Storage.StoreFile(r.Body, "newFile.md")
	if err != nil {
		respondWithError(w, 500, FILENOTSAVED)
		return
	}
	w.WriteHeader(201)
}

func (h *Handlers) HandleGetHtml(w http.ResponseWriter, r *http.Request) {
	fileName := r.PathValue("filename")
	if fileName == "" {
		respondWithError(w, 400, BADREQ)
	}
	// TODO:
	// convert, sotre it to disk, update db and send back
	htmlFilePath := h.State.Storage.GetFilePath(fileName + ".html")
	if htmlFilePath != "" {
		http.ServeFile(w, r, htmlFilePath)
		return
	}
	path := h.State.Storage.GetFilePath(fileName + ".md")
	if path == "" {
		respondWithError(w, 404, FILENOTFOUND)
		return
	}
	htmlFilePath, err := converter.ConvertToHtml(path, fileName)
	// update the db with the file version
	if err != nil {
		respondWithError(w, 500, "Unable to fetch html file")
		return
	}
	http.ServeFile(w, r, htmlFilePath)
}
