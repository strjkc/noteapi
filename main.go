package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/strjkc/noteapi/handlers"
	"github.com/strjkc/noteapi/internal/queries"
	"github.com/strjkc/noteapi/spellcheck"
	"github.com/strjkc/noteapi/state"
	"github.com/strjkc/noteapi/storage"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("No .env found")
	}

	db, err := sql.Open("sqlite3", "notes.db")
	if err != nil {
		panic("Unable to open db connection")
	}
	dbQueries := queries.New(db)
	// storageDir := os.Getenv("STORAGEDIR")
	dictDir := os.Getenv("DICTDIR")
	port := os.Getenv("PORT")
	// storage := storage.NewLocalStorage(storageDir)
	storage := storage.NewS3Storage()
	wmf := spellcheck.NewWordMapFactory(dictDir)
	state := state.NewState(storage, wmf, dbQueries)
	handlers := handlers.NewHandlers(state)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /app/get-html/{filename}", handlers.HandleGetHtml)
	mux.HandleFunc("POST /api/spellcheck/{locale}", handlers.HandleSpellCheck)
	mux.HandleFunc("POST /api/upload", handlers.HandleFileUpload)
	mux.HandleFunc("DELETE /api/files/{filename}", handlers.HandleRemoveFile)
	mux.HandleFunc("POST /api/users", handlers.HandleCreateUser)
	mux.HandleFunc("PUT /api/users", handlers.HandleUpdateUser)
	mux.HandleFunc("DELETE /api/users", handlers.HandleRemoveUser)
	mux.HandleFunc("POST /api/login", handlers.HandleLogin)
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	fmt.Printf("Listening on port: %s", port)
	server.ListenAndServe()
}
