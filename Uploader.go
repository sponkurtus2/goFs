package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type UploadHandler struct{}

type gofileUploadResponse struct {
	Status string `json:"status"`
	Data   struct {
		DownloadPage string `json:"downloadPage"`
	} `json:"data"`
}

func (uh *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(300 << 20) // 300 max MB per File
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

	downloadLink, err := saveFile(f, fh.Filename)
	if err != nil {
		log.Fatal(err)
	}

	uploadedFile = &MyFile{
		FileName:     fh.Filename,
		FileSize:     fh.Size,
		DownloadLink: downloadLink,
	}

	http.Redirect(w, r, "/fileInfo", http.StatusSeeOther)
}

func saveFile(file io.Reader, filename string) (string, error) {
	url := "https://upload.gofile.io/uploadfile"
	key := os.Getenv("API_TOKEN")

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("Failed to create form file: %w", err)
	}

	// Copy file data to the form
	if _, err := io.Copy(part, file); err != nil {
		return "", fmt.Errorf("failed to copy file data: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	req, err := http.NewRequest("POST", url, &body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	authorizationHeader := fmt.Sprintf("Bearer %s", key)
	req.Header.Add("Authorization", authorizationHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}
	defer resp.Body.Close()

	var uploadResp gofileUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if uploadResp.Status != "ok" {
		return "", fmt.Errorf("invalid response status: %s", uploadResp.Status)
	}

	return uploadResp.Data.DownloadPage, nil
}
