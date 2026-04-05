package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/strjkc/noteapi/handlers"
	"github.com/strjkc/noteapi/state"
	"github.com/strjkc/noteapi/storage"
)

func main() {
	err := godotenv.Load()
	storageDir := os.Getenv("STORAGEDIR")
	storage := storage.NewLocalStorage(storageDir)
	state := state.NewState(storage)
	handlers := handlers.NewHandlers(state)
	if err != nil {
		panic("No .env found")
	}
	port := os.Getenv("PORT")
	mux := http.NewServeMux()
	mux.HandleFunc("POST /spellcheck", handlers.HandleSpellCheck)
	mux.HandleFunc("POST /upload", handlers.HandleFileUpload)
	mux.HandleFunc("GET /fileAsHtml/{filename}", handlers.HandleGetHtml)
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	fmt.Printf("Listening on port: %s", port)
	server.ListenAndServe()
}
