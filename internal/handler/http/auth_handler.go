package http

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"qlsvgo/internal/domain/model"
	"qlsvgo/internal/infrastructure/jwt"
	"qlsvgo/internal/usecase"
	"qlsvgo/pkg/utils"
)

var UserUC *usecase.UserUsecase // Injected ở main
var JWTSecret string            // Injected ở main

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	user, err := UserUC.GetByUsername(req.Username)
	if err != nil || user.Password != req.Password {
		utils.Error(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	token, err := jwt.GenerateJWT(JWTSecret, user.ID, user.Role, 24)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "could not generate token")
		return
	}
	utils.JSON(w, http.StatusOK, map[string]string{"token": token})
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("RegisterHandler called")

	// Read the raw body for debugging
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		utils.Error(w, http.StatusBadRequest, "error reading request body")
		return
	}
	log.Printf("Raw request body: %s", string(bodyBytes))

	// Reset body for JSON decoder
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var req model.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("JSON decode error: %v", err)
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	log.Printf("Parsed user: %+v", req)

	if req.Username == "" || req.Password == "" || req.Email == "" || req.Role == "" {
		log.Printf("Validation failed: username='%s', password='%s', email='%s', role='%s'",
			req.Username, req.Password, req.Email, req.Role)
		utils.Error(w, http.StatusBadRequest, "username, password, email, role required")
		return
	}

	if err := UserUC.Register(&req); err != nil {
		log.Printf("UserUC.Register error: %v", err)
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("User registered successfully: %s", req.Username)
	utils.JSON(w, http.StatusCreated, req)
}
