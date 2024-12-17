package upload

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	// Parse the multipart form data
    err := r.ParseMultipartForm(100 << 20) // Limit file size to 10MB (adjust as needed)
    if err != nil {
        jsonErrorResponse(w, http.StatusBadRequest, "Unable to parse form data")
        return
    }

    // Retrieve the files from the request
    files := r.MultipartForm.File["files"]
    if len(files) == 0 {
        jsonErrorResponse(w, http.StatusBadRequest, "No files provided")
        return
    }

    // Scaleway bucket URL
    bucketURL := "https://long-term-strg-app.s3.fr-par.scw.cloud"

    for _, fileHeader := range files {
        // Open the uploaded file
        file, err := fileHeader.Open()
        if err != nil {
            jsonErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to open uploaded file: %v", err))
            return
        }
        defer file.Close()

        // Prepare the upload path (use fileHeader.Filename for the name)
        filePath := fmt.Sprintf("/storage/%s", fileHeader.Filename)

        // Create a signed PUT request
        request, err := signer.CreateSignedRequest("PUT", bucketURL+filePath, "fr-par")
        if err != nil {
            jsonErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error creating signed request: %v", err))
            return
        }

        // Attach file content to the request
        request.Body = file
        request.ContentLength = fileHeader.Size

        log.Printf("Uploading file: %s to path: %s\n", fileHeader.Filename, filePath)

        // Perform the PUT request
        client := &http.Client{}
        resp, err := client.Do(request)
        if err != nil || resp.StatusCode != http.StatusOK {
            log.Printf("Upload failed for file %s, status: %v\n", fileHeader.Filename, resp.Status)
            http.Error(w, fmt.Sprintf("Failed to upload file: %v", err), http.StatusInternalServerError)
            return
        }
        defer resp.Body.Close()

        log.Printf("File %s uploaded successfully", fileHeader.Filename)
    }

    // Respond with success
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "message": "Files uploaded successfully",
    })
}