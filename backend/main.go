package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Jokerjoker91/Storage-App/handlers/auth"
	"github.com/Jokerjoker91/Storage-App/handlers/getlist"
	"github.com/Jokerjoker91/Storage-App/handlers/login"
	"github.com/Jokerjoker91/Storage-App/handlers/upload"

	"github.com/rs/cors"
)


func main() {
	// Initialize authentication and login
	auth.InitializeAuth()
	login.InitializeLogin()

	// Main multiplexer
	mux := http.NewServeMux()

	// Serve all static files from the "public" directory
	fs := http.FileServer(http.Dir("../frontend/public"))
	mux.Handle("/", fs) // Serve index.html and other static files by default

	// Public route for login
	mux.HandleFunc("/api/login", login.LoginHandler)

	// Sub-mux for protected routes
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/api/upload-folder", upload.UploadFilesToBucket)
	protectedMux.HandleFunc("/api/get-bucket-contents", getlist.GetBucketContentsHandler)
	
	// Protect the home.html page with AuthMiddleware
	// mux.Handle("/home.html", auth.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "../frontend/public/home.html")
	// })))

	// Combine protected routes under AuthMiddleware
	mux.Handle("/api/", auth.AuthMiddleware(protectedMux))

	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins during development
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Wrap the router with the CORS middleware
	handler := c.Handler(mux)

	// Determine port
	port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // Default port for local development
    }

	log.Printf("Server starting on port %s...\n", port)

	// Start the server
	log.Fatal(http.ListenAndServe("localhost:"+port, handler))
}
