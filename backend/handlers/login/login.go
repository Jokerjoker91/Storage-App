package login

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Jokerjoker91/Storage-App/handlers/auth"

	"github.com/golang-jwt/jwt/v5"
)

// Initialize allowedUsers map
var allowedUsers map[string]string

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse includes the token
type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

// InitializeLogin loads allowed users from the environment
func InitializeLogin() {
	println("login/InitializeLogin")
	allowedUsers = make(map[string]string)

	usersEnv := os.Getenv("ALLOWED_USERS")
	if usersEnv == "" {
		log.Fatal("ALLOWED_USERS is not defined in the .env file")
	}

	for _, user := range strings.Split(usersEnv, ",") {
		parts := strings.Split(user, ":")
		if len(parts) != 2 {
			log.Fatalf("Invalid user entry in ALLOWED_USERS: %s", user)
		}
		allowedUsers[parts[0]] = parts[1]
	}
}

// LoginHandler handles user login and JWT generation
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	println("login/LoginHandler")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Decode the login request
	var loginReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	password, exists := allowedUsers[loginReq.Email]
	if !exists || password != loginReq.Password {
		response := LoginResponse{Success: false, Message: "Invalid email or password"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"email": loginReq.Email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(auth.JwtSecret)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

		// --- Added Logging for Debugging ---
		log.Printf("Generated JWT Token: %s\n, %+v", signedToken, claims)

	response := LoginResponse{Success: true, Message: "Login successful", Token: signedToken}
	json.NewEncoder(w).Encode(response)
}