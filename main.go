package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"encoding/json"
)

func getFolderMetadataHandler(w http.ResponseWriter, r *http.Request) {
	// Parse folder path from query parameter
	folderPath := getFolderPathFromArgs()
	if folderPath == "" {
		http.Error(w, "No folder path provided", http.StatusBadRequest)
		return
	}

	// Get folder metadata
	metadata, err := getFolderMetadata(folderPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get folder metadata: %s", err), http.StatusInternalServerError)
		return
	}

	// Respond with JSON encoded metadata
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(metadata); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %s", err), http.StatusInternalServerError)
		return
	}
}

func getFileContentHandler(w http.ResponseWriter, r *http.Request) {
	// Parse file name from query parameter
	fileName := r.URL.Query().Get("fileName")
	folderPath := getFolderPathFromArgs()
	if fileName == "" {
		http.Error(w, "Missing 'fileName' query parameter", http.StatusBadRequest)
		return
	}

	// Get file content
	fileContent, err := getFileContent(filepath.Join(folderPath, fileName))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read file content: %s", err), http.StatusInternalServerError)
		return
	}

	// Respond with file content as octet/stream
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(fileName)))
	w.Header().Set("Content-Type", "application/octet-stream")
	if _, err := w.Write(fileContent); err != nil {
		log.Printf("Failed to write file content to response: %s", err)
	}
}

type FileMetadata struct {
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	IsDirectory  bool      `json:"isDirectory"`
	LastModified time.Time `json:"lastModified"`
}

func getFolderMetadata(folderPath string) ([]FileMetadata, error) {
	var metadata []FileMetadata

	// Walk through the folder and its subdirectories
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root folder itself
		if path == folderPath {
			return nil
		}

		// Collect file/folder metadata
		fileMeta := FileMetadata{
			Name:         info.Name(),
			Size:         info.Size(),
			IsDirectory:  info.IsDir(),
			LastModified: info.ModTime(),
		}
		metadata = append(metadata, fileMeta)

		return nil
	})

	return metadata, err
}

func getFileContent(fileName string) ([]byte, error) {
	return os.ReadFile(fileName)
}

func main() {
		if len(os.Args) < 3 {
			fmt.Println("Usage: ./filebrowser <folderPath> <PORT>")
			os.Exit(1)
		}
	
		port := os.Args[2]
	
		http.HandleFunc("/folder-metadata", getFolderMetadataHandler)
		http.HandleFunc("/file-content", getFileContentHandler)
		log.Printf("Server is running on http://localhost:%s", port)
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}
	
func getFolderPathFromArgs() string {
	if len(os.Args) >= 2 {
		return os.Args[1]
	}
	return ""
}
