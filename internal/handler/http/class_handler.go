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

var ClassUC *usecase.ClassUsecase // Injected ở main

func CreateClassHandler(w http.ResponseWriter, r *http.Request) {
	var req model.Class
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" || req.CourseID == "" || req.TeacherID == "" {
		utils.Error(w, http.StatusBadRequest, "name, course_id, teacher_id required")
		return
	}
	if err := ClassUC.Create(&req); err != nil {
		logger.Error("create class: %v", err)
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusCreated, req)
}

func GetClassHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/classes/")
	class, err := ClassUC.QueryRepo.GetByID(id)
	if err != nil {
		if err == errors.ErrNotFound {
			utils.Error(w, http.StatusNotFound, "class not found")
			return
		}
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, class)
}

func UpdateClassHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/classes/")
	var req model.Class
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.ID = id
	if err := ClassUC.Update(&req); err != nil {
		logger.Error("update class: %v", err)
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, req)
}

func DeleteClassHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/classes/")
	if err := ClassUC.Delete(id); err != nil {
		logger.Error("delete class: %v", err)
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, map[string]string{"result": "deleted"})
}

func GetAllClassesHandler(w http.ResponseWriter, r *http.Request) {
	classes, err := ClassUC.QueryRepo.GetAll()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, classes)
}
