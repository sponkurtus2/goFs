package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const Addr = ":8080"

type MyFile struct {
	fileName string
	fileSize int64
}

type UploadHandler struct{ file *MyFile }

func (uh *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(10 << 20)

	f, fh, err := r.FormFile("file")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fmt.Printf("File name -> %s, File Size -> %d\n", fh.Filename, fh.Size)

	uploadedFile := &MyFile{
		fileName: fh.Filename,
		fileSize: fh.Size,
	}

	fmt.Fprintf(w, "Uploades file with name: %s, and size: %d\n", uploadedFile.fileName, uploadedFile.fileSize)
}

func openIndexView(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://localhost:63343/microServices/front/", http.StatusSeeOther)
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", openIndexView)
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
		fmt.Printf("Server failed: %s\n", err)
	}

}
