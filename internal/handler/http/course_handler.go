package http

import (
	"encoding/json"
	"net/http"
	"qlsvgo/internal/domain/model"
	"qlsvgo/internal/usecase"
	"qlsvgo/pkg/logger"
	"qlsvgo/pkg/utils"
	"strings"
)

var CourseUC *usecase.CourseUsecase // Injected ở main

func CreateCourseHandler(w http.ResponseWriter, r *http.Request) {
	var req model.Course
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		utils.Error(w, http.StatusBadRequest, "name required")
		return
	}
	if err := CourseUC.Create(&req); err != nil {
		logger.Error("create course: %v", err)
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusCreated, req)
}

func GetCourseHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/courses/")
	course, err := CourseUC.QueryRepo.GetByID(id)
	if err != nil {
		utils.Error(w, http.StatusNotFound, "course not found")
		return
	}
	utils.JSON(w, http.StatusOK, course)
}

func UpdateCourseHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/courses/")
	var req model.Course
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.ID = id
	if err := CourseUC.Update(&req); err != nil {
		logger.Error("update course: %v", err)
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, req)
}

func DeleteCourseHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/courses/")
	if err := CourseUC.Delete(id); err != nil {
		logger.Error("delete course: %v", err)
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, map[string]string{"result": "deleted"})
}

func GetAllCoursesHandler(w http.ResponseWriter, r *http.Request) {
	courses, err := CourseUC.QueryRepo.GetAll()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, courses)
}
