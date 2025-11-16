// handlers/files.go
package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form, with max file size
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	// Retrieve the file from the form-data
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// to  sanitize and generate a unique file name
	filename := filepath.Base(handler.Filename)    // simple sanitize
	destPath := filepath.Join("uploads", filename) // save in uploads/

	// Create des directory if doesnt exist
	os.MkdirAll("uploads", os.ModePerm)

	dest, err := os.Create(destPath)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	defer dest.Close()

	_, err = io.Copy(dest, file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %s", destPath)
}

func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	filePath := filepath.Join("uploads", filename)

	http.ServeFile(w, r, filePath)
}
