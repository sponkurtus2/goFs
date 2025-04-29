package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

const Addr = ":8080"

var uploadedFile *MyFile

type MyFile struct {
	FileName string
	FileSize int64
}

type UploadHandler struct{}

func (uh *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	f, fh, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	uploadedFile = &MyFile{
		FileName: fh.Filename,
		FileSize: fh.Size,
	}

	http.Redirect(w, r, "/fileInfo", http.StatusSeeOther)
}

func openIndexView(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func openUploadedFileView(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/uploadedFile.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, uploadedFile)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", openIndexView)
	mux.HandleFunc("/fileInfo", openUploadedFileView)
	mux.Handle("/upload", &UploadHandler{})

	s := &http.Server{
		Addr:         Addr,
		Handler:      mux,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	fmt.Printf("Starting server on: %s\n", Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %s\n", err)
	}
}
