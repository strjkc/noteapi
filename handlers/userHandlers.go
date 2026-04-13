package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/mattn/go-sqlite3"
	"github.com/strjkc/noteapi/auth"
	"github.com/strjkc/noteapi/internal/queries"
)

type UserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const (
	INVALIDPASSWORD = "Invalid Password for User"
)

func validateUsername(username string) bool {
	// TODO:
	// must contain letters only and no sopaces, land longer than 0
	return true
}

func (h *Handlers) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var userReq UserReq
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		respondWithError(w, 400, BADREQ)
		return
	}
	if !validateUsername(userReq.Username) {
		respondWithError(w, 400, BADREQ)
		return
	}
	if len(userReq.Password) < 6 {
		respondWithError(w, 400, BADREQ)
		return
	}

	dbUser, err := h.State.DbQueries.GetUserByUname(context.Background(), userReq.Username)
	if err != nil {
		respondWithError(w, 500, INTERNALERROR)
		return
	}
	ok, err := auth.ValidatePassword(userReq.Password, dbUser.Password)
	if err != nil {
		respondWithError(w, 500, INTERNALERROR)
		return
	}
	if !ok {
		respondWithError(w, 400, INVALIDPASSWORD)
		return
	}
	userID := strconv.Itoa(int(dbUser.ID))
	signedToken, err := auth.CreateToken(userID)

	respUser := struct {
		Username string `json:"username"`
		Token    string `json:"token"`
	}{
		Username: dbUser.Username,
		Token:    signedToken,
	}

	respData, err := json.Marshal(respUser)
	if err != nil {
		respondWithError(w, 500, INTERNALERROR)
		return
	}
	sendJson(w, 201, respData)
}

func (h *Handlers) HandleRemoveUser(w http.ResponseWriter, r *http.Request) {
	id, err := auth.ValidateToken(r.Header.Get("Authorization"))
	if err != nil {
		respondWithError(w, 500, INTERNALERROR)
		return
	}
	userID, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(w, 500, INTERNALERROR)
		return
	}
	data, err := h.State.DbQueries.RemoveUser(context.Background(), int64(userID))
	if err != nil {
		respondWithError(w, 500, INTERNALERROR)
		return
	}
	respData, err := json.Marshal(data)
	if err != nil {
		respondWithError(w, 500, INTERNALERROR)
		return
	}
	sendJson(w, 204, respData)
}

func (h *Handlers) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := auth.ValidateToken(r.Header.Get("Authorization"))
	if err != nil {
		respondWithError(w, 500, INTERNALERROR)
		return
	}
	userID, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(w, 500, INTERNALERROR)
		return
	}
	dbUser, err := h.State.DbQueries.GetUser(context.Background(), int64(userID))
	if err != nil {
		respondWithError(w, 400, BADREQ)
		return
	}
	var userReq UserReq
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		respondWithError(w, 400, BADREQ)
		return
	}
	if !validateUsername(userReq.Username) {
		respondWithError(w, 400, BADREQ)
		return
	}
	if len(userReq.Password) < 6 {
		respondWithError(w, 400, BADREQ)
		return
	}

	if dbUser.Username != userReq.Username {
		u2, err := h.State.DbQueries.GetUserByUname(context.Background(), userReq.Username)
		if err != nil {
			// TODO: update user
		}
		respondWithError(w, 400, USEREXISTS)
		return
	}

	// username and password required
	// if token valid, check if user from db.Username == username - update password
	// if db.Username != username check if username exists, if no update username
	// change updated_at
	// return updated data
}

func (h *Handlers) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var user UserReq
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, 400, BADREQ)
		return
	}

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		respondWithError(w, 500, INTERNALERROR)
		return
	}

	newDbUser := queries.CreateUserParams{
		Username:  user.Username,
		Password:  hashedPassword,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}
	dbUser, err := h.State.DbQueries.CreateUser(context.Background(), newDbUser)
	if err != nil {
		var sqlErr sqlite3.Error
		errors.As(err, &sqlErr)
		if sqlErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			respondWithError(w, 409, USEREXISTS)
			return
		}
		fmt.Println("Unable to serialize user", err)
		respondWithError(w, 500, INTERNALERROR)
		return
	}
	respUser := struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}{
		ID:        int(dbUser.ID),
		Username:  dbUser.Username,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}

	respData, err := json.Marshal(respUser)
	if err != nil {
		respondWithError(w, 500, INTERNALERROR)
		return
	}
	sendJson(w, 201, respData)
}
