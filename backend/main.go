package main

import (
	"log"
	"net/http"
	"os"

	"storage-app/handlers/getlist"
	"storage-app/handlers/login"
	"storage-app/handlers/upload"

	"github.com/rs/cors"
)


func main() {

	// Load allowed users from .env
	login.LoadAllowedUsersFromEnv()
	
	mux := http.NewServeMux()

	// Serve static files
	fileServer := http.FileServer(http.Dir("../frontend/public")) // Adjust the folder path if needed
	mux.Handle("/", fileServer)

	// Handle login route
	mux.HandleFunc("/login", login.LoginHandler)

	// Handle the upload folder route
	mux.HandleFunc("/upload-folder", upload.UploadFolderHandler)

	// Bucket contents route
	mux.HandleFunc("/api/get-bucket-contents", getlist.GetBucketContentsHandler)

	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins during development
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Wrap the router with the CORS middleware
	handler := c.Handler(mux)

	port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // Default port for local development
    }

	log.Printf("Server starting on port %s...\n", port)

	// Start the server
	http.ListenAndServe("localhost:"+port, handler)
}
