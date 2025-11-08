package http

import (
	"encoding/json"
	"net/http"
	"qlsvgo/internal/domain/model"
	"qlsvgo/internal/usecase"
	"qlsvgo/pkg/errors"
	"qlsvgo/pkg/logger"
	"qlsvgo/pkg/utils"
	"strings"
)

var TeacherUC *usecase.TeacherUsecase // Injected ở main

func CreateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	var req model.Teacher
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" || req.Email == "" {
		utils.Error(w, http.StatusBadRequest, "name and email required")
		return
	}
	if err := TeacherUC.Create(&req); err != nil {
		logger.Error("create teacher: %v", err)
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusCreated, req)
}

func GetTeacherHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/teachers/")
	teacher, err := TeacherUC.QueryRepo.GetByID(id)
	if err != nil {
		if err == errors.ErrNotFound {
			utils.Error(w, http.StatusNotFound, "teacher not found")
			return
		}
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, teacher)
}

func UpdateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/teachers/")
	var req model.Teacher
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.ID = id
	if err := TeacherUC.Update(&req); err != nil {
		logger.Error("update teacher: %v", err)
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, req)
}

func DeleteTeacherHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/teachers/")
	if err := TeacherUC.Delete(id); err != nil {
		logger.Error("delete teacher: %v", err)
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, map[string]string{"result": "deleted"})
}

func GetAllTeachersHandler(w http.ResponseWriter, r *http.Request) {
	teachers, err := TeacherUC.QueryRepo.GetAll()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, teachers)
}
