package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "embed"
	_ "github.com/joho/godotenv/autoload"
)

//go:embed templates/index.html
var indexPage embed.FS

//go:embed templates/uploadedFile.html
var uploadedFilePage embed.FS

var uploadedFile *MyFile

type MyFile struct {
	FileName     string
	FileSize     int64
	DownloadLink string
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", openIndexView)
	mux.HandleFunc("/fileInfo", openUploadedFileView)
	mux.Handle("/upload", &UploadHandler{})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	s := &http.Server{
		Addr:         port,
		Handler:      mux,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	fmt.Printf("Starting server on: %s\n", port)
	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %s\n", err)
	}
}
