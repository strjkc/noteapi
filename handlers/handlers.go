package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/strjkc/noteapi/converter"
	"github.com/strjkc/noteapi/internal/queries"
	"github.com/strjkc/noteapi/spellcheck"
	"github.com/strjkc/noteapi/state"
)

const (
	INTERNALERROR = "Internal Server Error"
	FILENOTFOUND  = "The Requested File Can Not Be Found"
	FILENOTSAVED  = "The File Was Not Stored Due to an Internal Error"
	BADREQ        = "Bad Request"
	USEREXISTS        = "Username Already Exists"
)

type Handlers struct {
	State *state.State
}

func NewHandlers(state *state.State) *Handlers {
	h := Handlers{State: state}
	return &h
}

func (h *Handlers) HandleSpellCheck(w http.ResponseWriter, r *http.Request) {
	locale := r.PathValue("locale")
	if locale == "" {
		respondWithError(w, 400, BADREQ)
		return
	}

	mr, err := r.MultipartReader()
	if err != nil {
		respondWithError(w, 400, BADREQ)
		return
	}

	wm, err := h.State.WordMapFactory.WordMap(locale)
	if err != nil {
		respondWithError(w, 500, "Locale not supported")
		return
	}
	checker := spellcheck.NewChecker(wm)
	parser := spellcheck.NewParser(checker)
	errors, err := parser.SpellCheckerService(mr)
	if err != nil {
		respondWithError(w, 500, INTERNALERROR)
	}
	fmt.Println(errors)
	json := marshalJson(errors)
	sendJson(w, 200, json)
}

func (h *Handlers) HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	mr, err := r.MultipartReader()
	if err != nil {
		respondWithError(w, 400, BADREQ)
		return
	}
	err = h.State.Storage.StoreFile(mr)
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
	htmlFilePath := h.State.Storage.FileURL(fileName + ".html")
	if htmlFilePath != "" {
		http.ServeFile(w, r, htmlFilePath)
		return
	}
	path := h.State.Storage.StorageDir()
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

func (h *Handlers) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	type UserReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var user UserReq
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, 400, BADREQ)
		return
	}
	_, err := h.State.DbQueries.GetUser(context.Background(), user.Username)
	if err == nil {
		respondWithError(w, 400, USEREXISTS)
		return
	}

	params := argon2id.Params{
		Memory: 128 * 1024,
		Iterations: 10,
		Parallelism: uint8(runtime.NumCPU()),
		SaltLength: 16,
		KeyLength: 32,
	}
	hashedPassword, err := argon2id.CreateHash(user.Password, &params)
	if err != nil {
		respondWithError(w, 500, INTERNALERROR)
		return
	}
	newDbUser := queries.CreateUserParams{
		Username: user.Username,
		Password: hashedPassword,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}
	dbUser, err := h.State.DbQueries.CreateUser(context.Background(), newDbUser)
	if err != nil {
		fmt.Println("Unable to serialize user", err)
		respondWithError(w, 500, INTERNALERROR)
		return
	}

	respUser := struct {
		ID int `json:"id"`
		Username string `json:"username"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}{
		int(dbUser.ID),
		dbUser.Username,
		dbUser.CreatedAt,
		dbUser.UpdatedAt,
	}

	respData, err := json.Marshal(respUser)
	if err != nil {
		respondWithError(w, 500, INTERNALERROR)
		return
	}
	sendJson(w, 201, respData)
}
