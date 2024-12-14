package login

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

// Initialize allowedUsers map
var allowedUsers = make(map[string]string)

// LoadAllowedUsersFromEnv loads allowed users from the .env file
func LoadAllowedUsersFromEnv() {
	users := os.Getenv("ALLOWED_USERS")
	if users == "" {
		panic("ALLOWED_USERS environment variable is not set") // Fail fast if users are not configured
	}

	// Parse users from the ALLOWED_USERS string
	userPairs := strings.Split(users, ",")
	for _, pair := range userPairs {
		parts := strings.Split(pair, ":")
		if len(parts) == 2 {
			allowedUsers[parts[0]] = parts[1]
		} else {
			panic("Invalid format in ALLOWED_USERS environment variable") // Fail fast if the format is invalid
		}
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var loginReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	password, exists := allowedUsers[loginReq.Email]
	if !exists || password != loginReq.Password {
		response := LoginResponse{Success: false, Message: "Invalid email or password"}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := LoginResponse{Success: true, Message: "Login successful"}
	json.NewEncoder(w).Encode(response)
}