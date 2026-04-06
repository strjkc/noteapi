package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/strjkc/noteapi/handlers"
	"github.com/strjkc/noteapi/spellcheck"
	"github.com/strjkc/noteapi/state"
	"github.com/strjkc/noteapi/storage"
)

func main() {
	err := godotenv.Load()
	storageDir := os.Getenv("STORAGEDIR")
	dictDir := os.Getenv("DICTDIR")
	port := os.Getenv("PORT")
	storage := storage.NewLocalStorage(storageDir)
	wmf := spellcheck.NewWordMapFactory(dictDir)
	state := state.NewState(storage, wmf)
	handlers := handlers.NewHandlers(state)
	if err != nil {
		panic("No .env found")
	}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /spellcheck/{locale}", handlers.HandleSpellCheck)
	mux.HandleFunc("POST /upload", handlers.HandleFileUpload)
	mux.HandleFunc("GET /fileAsHtml/{filename}", handlers.HandleGetHtml)
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	fmt.Printf("Listening on port: %s", port)
	server.ListenAndServe()
}
