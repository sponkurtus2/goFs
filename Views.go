package main

import (
	"html/template"
	"net/http"
)

func openIndexView(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(indexPage, "templates/index.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func openUploadedFileView(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(uploadedFilePage, "templates/uploadedFile.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, uploadedFile)
}
