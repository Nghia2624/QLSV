package http

import (
	"net/http"
	"qlsvgo/internal/repository"
	"qlsvgo/pkg/utils"
	"strings"
)

var UserRepo repository.UserRepository

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := UserRepo.GetAll()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, users)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/users/")
	if id == "" || id == "users" || id == "/" {
		utils.Error(w, http.StatusBadRequest, "invalid user ID")
		return
	}
	user, err := UserRepo.GetByID(id)
	if err != nil {
		utils.Error(w, http.StatusNotFound, "user not found")
		return
	}
	utils.JSON(w, http.StatusOK, user)
}
