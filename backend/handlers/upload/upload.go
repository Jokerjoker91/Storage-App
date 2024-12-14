package upload

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

type FileDetails struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type UploadRequest struct {
	Files []FileDetails `json:"files"`
}

func init() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }
}

// Function to upload files to Scaleway S3 Glacier
func uploadToS3(bucketName, filePath, fileName string, fileContent []byte) error {
	// Scaleway S3 endpoint
	endpoint := "https://long-term-strg-app.s3.fr-par.scw.cloud"

	// Scaleway access credentials (replace with secure method in production)
	accessKey := os.Getenv("SCW_ACCESS_KEY")
	secretKey := os.Getenv("SCW_SECRET_KEY")

	if accessKey == "" || secretKey == "" {
		return fmt.Errorf("missing Scaleway credentials")
	}

	// Create AWS SDK config for Scaleway S3
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "scw",
			URL:           endpoint,
			SigningRegion: "fr-par",
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("fr-par"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return fmt.Errorf("unable to load SDK config: %v", err)
	}

	// Create S3 client
	s3Client := s3.NewFromConfig(cfg)

	// Set the object key (simulating folder structure)
	objectKey := "Test folder/" + filePath + "/" + fileName

	// Upload object to Scaleway S3
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:       aws.String(bucketName),
		Key:          aws.String(objectKey),
		Body:         bytes.NewReader(fileContent),
		StorageClass: "DEEP_ARCHIVE", // Glacier equivalent for Scaleway
	})

	return err
}

// UploadFolderHandler handles file uploads
func UploadFolderHandler(w http.ResponseWriter, r *http.Request) {

    // Enable CORS headers
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    // Handle preflight OPTIONS request
    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }

	var req UploadRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Upload each file concurrently
	var wg sync.WaitGroup
	for _, file := range req.Files {
		wg.Add(1)
		go func(file FileDetails) {
			defer wg.Done()

			// Read file content
			fileContent, err := os.ReadFile(file.Name)
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", file.Name, err)
				return
			}

			// Upload file to S3 (Scaleway)
			err = uploadToS3("long-term-strg-app", file.Path, file.Name, fileContent)
			if err != nil {
				fmt.Printf("Error uploading %s: %v\n", file.Name, err)
			} else {
				fmt.Printf("Successfully uploaded %s\n", file.Name)
			}
		}(file)
	}

	wg.Wait()

    // Create a response struct
    response := map[string]interface{}{
        "success": true,
        "message": "Folder uploaded successfully",
    }

    // Set content type to JSON
    w.Header().Set("Content-Type", "application/json")

    // Write the response as JSON
    jsonResponse, err := json.Marshal(response)
    if err != nil {
        http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
        return
    }

    // Send the response
    w.Write(jsonResponse)
}