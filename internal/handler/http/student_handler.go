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

var StudentUC *usecase.StudentUsecase // Injected ở main

func CreateStudentHandler(w http.ResponseWriter, r *http.Request) {
	var req model.Student
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" || req.Email == "" {
		utils.Error(w, http.StatusBadRequest, "name and email required")
		return
	}
	if err := StudentUC.Create(&req); err != nil {
		logger.Error("create student: %v", err)
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusCreated, req)
}

func GetStudentHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/students/")
	student, err := StudentUC.QueryRepo.GetByID(id)
	if err != nil {
		if err == errors.ErrNotFound {
			utils.Error(w, http.StatusNotFound, "student not found")
			return
		}
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, student)
}

func UpdateStudentHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/students/")
	var req model.Student
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.ID = id
	if err := StudentUC.Update(&req); err != nil {
		logger.Error("update student: %v", err)
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, req)
}

func DeleteStudentHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/students/")
	if err := StudentUC.Delete(id); err != nil {
		logger.Error("delete student: %v", err)
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, map[string]string{"result": "deleted"})
}

func GetAllStudentsHandler(w http.ResponseWriter, r *http.Request) {
	students, err := StudentUC.QueryRepo.GetAll()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, students)
}
