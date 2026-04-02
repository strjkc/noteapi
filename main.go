package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/strjkc/noteapi/handlers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("No .env found")
	}
	port := os.Getenv("PORT")
	mux := http.NewServeMux()
	mux.HandleFunc("POST /spellcheck", handlers.HandleSpellCheck)
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	fmt.Printf("Listening on port: %s", port)
	server.ListenAndServe()
}
