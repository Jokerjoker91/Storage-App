package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

// JwtSecret stores the generated JWT secret
var JwtSecret []byte

// GenerateJWTSecret derives a secret key from ALLOWED_USERS
func GenerateJWTSecret() []byte {
	println("auth/GenerateJWTSecret")
	allowedUsers := os.Getenv("ALLOWED_USERS")
	if allowedUsers == "" {
		log.Fatal("ALLOWED_USERS is not defined in the .env file")
	}

	// Create a SHA-256 hash of ALLOWED_USERS
	hash := sha256.Sum256([]byte(allowedUsers))

	// Convert to a hex-encoded string
	secret := hex.EncodeToString(hash[:])
	return []byte(secret)
}

// InitializeAuth initializes the JWT secret and loads the .env file
func InitializeAuth() {
	println("auth/InitializeAuth")
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Generate and store the JWT secret
	JwtSecret = GenerateJWTSecret()
	log.Println("JWT Secret generated successfully.")
}

// AuthMiddleware validates the JWT
func AuthMiddleware(next http.Handler) http.Handler {
	println("auth/AuthMiddleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, fmt.Sprintf("Missing Authorization header: %s", authHeader), http.StatusUnauthorized)
			return
		}

		// Extract the token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]

		// Parse and validate the token with claims
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return JwtSecret, nil
		})

		// Handle token parsing errors
		if err != nil {
			log.Printf("Error parsing token: %v\n", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Validate claims
		email, ok := claims["email"].(string)
		if !ok || email == "" {
			http.Error(w, "Invalid token claims: missing email", http.StatusUnauthorized)
			return
		}

		// Validate expiration (if present)
		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				http.Error(w, "Token has expired", http.StatusUnauthorized)
				return
			}
		}

		// Pass request to the next handler
		next.ServeHTTP(w, r)
	})
}
