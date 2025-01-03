package upload

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

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

func DecodeFilename(encodedFilename string) string {
    decodedFilename, err := url.QueryUnescape(encodedFilename)
    if err != nil {
        log.Printf("Error decoding filename %s: %v", encodedFilename, err)
        return encodedFilename // Fall back to the encoded name if decoding fails
    }
    return decodedFilename
}

// Handler to upload files to Scaleway bucket
func UploadFilesToBucket(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form data
    err := r.ParseMultipartForm(100 << 20) // Limit file size to 10MB
    if err != nil {
        log.Printf("Error parsing form data: %v\n", err)
        jsonErrorResponse(w, http.StatusBadRequest, "Unable to parse form data")
        return
    }

    // Retrieve the files from the request
    files := r.MultipartForm.File["files"]
    if len(files) == 0 {
        log.Println("No files found in the request")
        jsonErrorResponse(w, http.StatusBadRequest, "No files provided")
        return
    }

    folderPath := r.FormValue("folder")

    // Scaleway bucket URL
    bucketURL := "https://long-term-strg-app.s3.fr-par.scw.cloud"

    for _, fileHeader := range files {
        // Decode the filename
        decodedFilename := DecodeFilename(fileHeader.Filename)

        // Open the uploaded file
        file, err := fileHeader.Open()
        if err != nil {
            log.Printf("Failed to open file %s: %v\n", decodedFilename, err)
            jsonErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to process file %s", decodedFilename))
            return
        }
        defer file.Close()

       // Construct the file path for the bucket
       filePath := fmt.Sprintf("/%s/%s", folderPath, decodedFilename)

        // Create a signed PUT request
        request, err := signer.CreateSignedRequest("PUT", bucketURL+filePath, "fr-par")
        if err != nil {
            log.Printf("Error creating signed request for file %s: %v\n", decodedFilename, err)
            jsonErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error creating signed request: %v", err))
            return
        }

        // Read file content
        fileBytes, err := io.ReadAll(file)
        if err != nil {
            log.Printf("Error reading file %s: %v\n", decodedFilename, err)
            jsonErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error reading file %s", decodedFilename))
            return
        }

        // Attach file content to the request
        request.Body = io.NopCloser(bytes.NewReader(fileBytes))
        request.ContentLength = int64(len(fileBytes))

        log.Printf("Uploading file: %s to path: %s\n", decodedFilename, filePath)

        // Perform the PUT request
        client := &http.Client{}
        resp, err := client.Do(request)
        if err != nil {
            log.Printf("Upload failed for file %s, status: %v\n", decodedFilename, err)
            jsonErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to upload file %s: %v", decodedFilename, err))
            return
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            log.Printf("Upload failed for file %s, status: %v\n", decodedFilename, resp.Status)
            jsonErrorResponse(w, http.StatusForbidden, fmt.Sprintf("Failed to upload file %s: Forbidden", decodedFilename))
            return
        }

        log.Printf("File %s uploaded successfully", decodedFilename)
    }

    // Respond with success
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "message": "Files uploaded successfully",
    })
}