package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// Mocked file content for testing getFileContent function
var mockFileContent = []byte("This is a mock file content.")

// Helper function to create a temporary folder for testing and return its path
func createTempFolderForTest(t *testing.T) string {
	tmpDir := t.TempDir()
	testFolder := filepath.Join(tmpDir, "test_folder")
	err := os.Mkdir(testFolder, 0755)
	if err != nil {
		t.Fatalf("Failed to create test folder: %s", err)
	}
	return testFolder
}

func TestGetFolderMetadataHandler(t *testing.T) {
	// Create a temporary folder for testing
	testFolder := createTempFolderForTest(t)
	defer os.RemoveAll(testFolder)

	// Create some dummy files and subfolders in the test folder
	dummyFiles := []string{"file1.txt", "file2.txt"}
	dummySubfolders := []string{"subfolder1", "subfolder2"}
	for _, filename := range dummyFiles {
		filePath := filepath.Join(testFolder, filename)
		file, err := os.Create(filePath)
		if err != nil {
			t.Fatalf("Failed to create dummy file: %s", err)
		}
		file.Close()
	}
	for _, folderName := range dummySubfolders {
		folderPath := filepath.Join(testFolder, folderName)
		err := os.Mkdir(folderPath, 0755)
		if err != nil {
			t.Fatalf("Failed to create dummy subfolder: %s", err)
		}
	}

	// Create a mock request to the /folder-metadata endpoint
	req := httptest.NewRequest("GET", "/folder-metadata", nil)

	// Create a mock response writer
	recorder := httptest.NewRecorder()

	// Set the folder path query parameter
	q := req.URL.Query()
	q.Add("path", testFolder)
	req.URL.RawQuery = q.Encode()

	// Call the getFolderMetadataHandler function with the mock request and response writer
	getFolderMetadataHandler(recorder, req)

	// Check the response status code
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, recorder.Code)
	}

	// Parse the JSON response
	var metadata []FileMetadata
	err := json.NewDecoder(recorder.Body).Decode(&metadata)
	if err != nil {
		t.Fatalf("Failed to decode JSON response: %s", err)
	}

	// Check the metadata count (expecting dummyFiles and dummySubfolders)
	expectedItemCount := len(dummyFiles) + len(dummySubfolders)
	if len(metadata) != expectedItemCount {
		t.Errorf("Expected %d items in metadata, but got %d", expectedItemCount, len(metadata))
	}
}

func TestGetFileContentHandler(t *testing.T) {
	// Create a temporary folder for testing
	testFolder := createTempFolderForTest(t)
	defer os.RemoveAll(testFolder)

	// Create a dummy file in the test folder
	fileName := "test_file.txt"
	filePath := filepath.Join(testFolder, fileName)
	err := os.WriteFile(filePath, mockFileContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create dummy file: %s", err)
	}

	// Create a mock request to the /file-content endpoint with the desired fileName query parameter
	req := httptest.NewRequest("GET", fmt.Sprintf("/file-content?fileName=%s", fileName), nil)

	// Create a mock response writer
	recorder := httptest.NewRecorder()

	// Call the getFileContentHandler function with the mock request and response writer
	getFileContentHandler(recorder, req)

	// Check the response status code
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, recorder.Code)
	}

	// Check the response headers
	contentDisposition := recorder.Header().Get("Content-Disposition")
	expectedContentDisposition := fmt.Sprintf("attachment; filename=%s", fileName)
	if contentDisposition != expectedContentDisposition {
		t.Errorf("Expected Content-Disposition header: %s, but got: %s", expectedContentDisposition, contentDisposition)
	}

	// Check the response body (file content)
	if !bytes.Equal(recorder.Body.Bytes(), mockFileContent) {
		t.Errorf("Expected file content: %v, but got: %v", mockFileContent, recorder.Body.Bytes())
	}
}
