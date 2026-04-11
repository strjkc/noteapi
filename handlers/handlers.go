package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/strjkc/noteapi/auth"
	"github.com/strjkc/noteapi/converter"
	"github.com/strjkc/noteapi/internal/queries"
	"github.com/strjkc/noteapi/spellcheck"
	"github.com/strjkc/noteapi/state"
)

// TODO:
// Logging, replace all double failures with loud logging
const (
	INTERNALERROR     = "Internal Server Error"
	FILENOTFOUND      = "The Requested File Can Not Be Found"
	FILENOTSAVED      = "The File Was Not Stored Due to an Internal Error"
	BADREQ            = "Bad Request"
	USEREXISTS        = "Username Already Exists"
	INCONSISTENTSTATE = "There is a Mismatch Between the State of the Storage and the DB"
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

// TODO: if file is not .md then return error
func (h *Handlers) HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	mr, err := r.MultipartReader()
	if err != nil {
		respondWithError(w, 400, BADREQ)
		return
	}
	uid, err := auth.ValidateToken(r.Header.Get("Authorization"))
	if err != nil {
		fmt.Println("Error")
	}
	userID, err := strconv.Atoi(uid)
	if err != nil {
		respondWithError(w, 500, INTERNALERROR)
		return
	}

	user, err := h.State.DbQueries.GetUser(context.Background(), int64(userID))
	if err != nil {
		respondWithError(w, 500, BADREQ)
		return
	}

	fileName, tmpFileName, err := h.State.Storage.StoreFile(mr)
	if err != nil {
		respondWithError(w, 500, FILENOTSAVED)
		return
	}

	dbMDFile, err := h.State.DbQueries.GetFile(context.Background(), queries.GetFileParams{Name: fileName, UserID: user.ID})
	if err == nil {
		err = h.State.Storage.RenameFile(fileName, "backup_"+fileName)
		if err != nil {
			respondWithError(w, 500, FILENOTSAVED)
			return
		}

		err = h.State.Storage.RenameFile(tmpFileName, fileName)
		if err != nil {
			respondWithError(w, 500, FILENOTSAVED)
			return
		}

		updateFile := queries.UpdateFileParams{
			UpdatedAt: string(time.Now().Format(time.RFC3339)),
			ID:        dbMDFile.ID,
		}
		err = h.State.DbQueries.UpdateFile(context.Background(), updateFile)
		if err != nil {
			err := h.State.Storage.DeleteFile(fileName)
			if err != nil {
				fmt.Println("file exists on disk but we could not delete it")
				respondWithError(w, 500, INTERNALERROR)
				return
			}
			// err
			err = h.State.Storage.RenameFile("backup_"+fileName, fileName)
			if err != nil {
				fmt.Println("file exists on disk but we could not delete it")
				respondWithError(w, 500, INTERNALERROR)
				return
			}
		}

	} else {
		err := h.State.Storage.RenameFile(tmpFileName, fileName)
		if err != nil {
			// remove temp file
			err := h.State.Storage.DeleteFile(tmpFileName)
			if err != nil {
				// TODO: handle in a  better way
				fmt.Println("major error")
			}
			respondWithError(w, 500, FILENOTSAVED)
			return
		}
		// we update the db with the new updated ad
		// TODO:
		fileData := queries.CreateFileParams{
			Name:      fileName,
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
			UserID:    user.ID,
		}
		_, err = h.State.DbQueries.CreateFile(context.Background(), fileData)
		if err != nil {
			err := h.State.Storage.DeleteFile(fileName)
			if err != nil {
				// TODO: handle this in a better way
				fmt.Println("Something mayor went wrong, the file exists on disk but could not be deleted!")
			}
			fmt.Println("Erorr saving to db: ", err)
			respondWithError(w, 500, FILENOTSAVED)
			return
		}
	}
	w.WriteHeader(201)
}

func (h *Handlers) HandleGetHtml(w http.ResponseWriter, r *http.Request) {
	fileName := r.PathValue("filename")
	if fileName == "" {
		respondWithError(w, 400, BADREQ)
	}
	// TODO: auth the user
	user, err := h.State.DbQueries.GetUser(context.Background(), 1)
	if err != nil {
		respondWithError(w, 500, BADREQ)
		return
	}

	htmlFilePath := h.State.Storage.FileURL(fileName + ".html")
	dbHTMLFile, err := h.State.DbQueries.GetFile(context.Background(), queries.GetFileParams{Name: fileName + ".html", UserID: user.ID})
	if htmlFilePath == "" && err != nil {
		// file doesn't exist on the disk and in the db, we create and serve
		fmt.Println("No db entry, not file on disk")
		path := h.State.Storage.StorageDir()
		if path == "" {
			respondWithError(w, 404, FILENOTFOUND)
			return
		}
		htmlFilePath, err = converter.ConvertToHtml(path, fileName)
		fileData := queries.CreateFileParams{
			Name:      fileName + ".html",
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
			UserID:    user.ID,
		}
		if err != nil {
			respondWithError(w, 500, "Unable to fetch html file")
			return
		}
		_, err = h.State.DbQueries.CreateFile(context.Background(), fileData)
		if err != nil {
			err := h.State.Storage.DeleteFile(fileName + ".html")
			if err != nil {
				// TODO: handle this in a better way
				fmt.Println("something mayor went wrong, we couln not write to the db")
			}
			respondWithError(w, 500, FILENOTSAVED)
			return
		}
		// update the db with the file version

		http.ServeFile(w, r, htmlFilePath)
		return
	} else if htmlFilePath == "" && err == nil {
		err := h.State.DbQueries.Deleted(context.Background(), dbHTMLFile.ID)
		if err != nil {
			// TODO: handle this in a better way
			fmt.Println("something mayor went wrong, we couln not write to the db")
		}
		respondWithError(w, 500, INCONSISTENTSTATE)
		return
	} else if htmlFilePath != "" && err != nil {
		fmt.Printf("Inconsistent state, file: %s is on disk but not in db\n", htmlFilePath)
		respondWithError(w, 500, INCONSISTENTSTATE)
		return
	}
	// if we are here, it means file is on disk and in the db
	// we need to check if it's stale or not

	dbMDFile, err := h.State.DbQueries.GetFile(context.Background(), queries.GetFileParams{Name: fileName + ".md", UserID: user.ID})
	if err != nil {
		// if its not in db, we can assume the html file is the latest representation and we serve it
		http.ServeFile(w, r, htmlFilePath)
		return
	}

	HTMLUpdatedAt, err := time.Parse(time.RFC3339, dbHTMLFile.UpdatedAt)
	if err != nil {
		respondWithError(w, 500, INTERNALERROR)
	}
	MDUpdatedAt, err := time.Parse(time.RFC3339, dbMDFile.UpdatedAt)
	if err != nil {
		respondWithError(w, 500, INTERNALERROR)
	}
	if HTMLUpdatedAt.Before(MDUpdatedAt) {
		fmt.Println("found on disk, older than orig")
		mdFilePath := h.State.Storage.FileURL(fileName + ".md")
		if mdFilePath == "" {
			// TODO: mark md file in db as deleted
			err := h.State.DbQueries.Deleted(context.Background(), dbMDFile.ID)
			if err != nil {
				// TODO: handle this in a better way
				fmt.Println("something mayor went wrong, we couln not write to the db")
			}
			// inconsistent state, md file is in db but not on disk
			respondWithError(w, 500, INCONSISTENTSTATE)
			return

		}

		path := h.State.Storage.StorageDir()
		if path == "" {
			respondWithError(w, 404, FILENOTFOUND)
			return
		}
		htmlFilePath, err = converter.ConvertToHtml(path, fileName)
		fileData := queries.CreateFileParams{
			Name:      fileName + ".html",
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
			UserID:    user.ID,
		}
		if err != nil {
			respondWithError(w, 500, "Unable to fetch html file")
			return
		}
		_, err = h.State.DbQueries.CreateFile(context.Background(), fileData)
		if err != nil {
			err := h.State.Storage.DeleteFile(fileName + ".html")
			if err != nil {
				// TODO: handle in a better way
				fmt.Println("Error")
			}
			respondWithError(w, 500, FILENOTSAVED)
			return
		}
		// update the db with the file version
		http.ServeFile(w, r, htmlFilePath)
		return
		// convert again and update the timestamps
	}
	// we are good we should server

	fmt.Println("found on disk, newer than orig")
	http.ServeFile(w, r, htmlFilePath)
}
