package upload

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Jokerjoker91/Storage-App/handlers/signer"
)

// Struct to represent the file data received from the frontend
type FileData struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func jsonErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// Handler to upload files to Scaleway bucket
func UploadFilesToBucket(w http.ResponseWriter, r *http.Request) {
	// Parse incoming JSON data (files list)
	var requestBody struct {
		Files []FileData `json:"files"`
	}

	// Decode the incoming JSON request body
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		jsonErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to parse request body: %v", err))
		return
	}

	// Iterate over the file list and upload each file
	for _, fileData := range requestBody.Files {
		// Open the file
		file, err := os.Open(fileData.Path)
		if err != nil {
			jsonErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to open file: %v", err))
			return
		}
		defer file.Close()

		// Get file stats (including size)
		fileInfo, err := file.Stat()
		if err != nil {
			jsonErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get file info: %v", err))
			return
		}

		// URL for the Scaleway object storage
		bucketURL := "https://long-term-strg-app.s3.fr-par.scw.cloud"

		// Prepare the file path (uploading to root folder in the bucket)
		filePath := fmt.Sprintf("/storage/%s", fileData.Name)

		// Create a signed request to Scaleway
		request, err := signer.CreateSignedRequest("PUT", bucketURL+filePath, "fr-par")
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating signed request: %v", err), http.StatusInternalServerError)
			return
		}

		// Set the file content in the PUT request
		request.Body = io.NopCloser(file) // Attach file to the request body
		request.ContentLength = fileInfo.Size() // Set file size

		log.Printf("Uploading file: %s to path: %s\n", fileData.Name, filePath)
		log.Printf("Request Headers: %v\n", request.Header)

		// Make the actual PUT request to Scaleway
		client := &http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error uploading file to Scaleway: %v", err), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Check the response status
		if resp.StatusCode != http.StatusOK {
			http.Error(w, fmt.Sprintf("Failed to upload file: %v", resp.Status), http.StatusInternalServerError)
			return
		}

		log.Printf("File %s uploaded successfully", fileData.Name)
	}

	// Respond back with success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Files uploaded successfully",
	})
}