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

var RegistrationUC *usecase.RegistrationUsecase // Injected ở main

func CreateRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	var req model.Registration
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.StudentID == "" || req.ClassID == "" {
		utils.Error(w, http.StatusBadRequest, "student_id and class_id required")
		return
	}
	if err := RegistrationUC.Create(&req); err != nil {
		logger.Error("create registration: %v", err)
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusCreated, req)
}

func GetRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/registrations/")
	reg, err := RegistrationUC.QueryRepo.GetByID(id)
	if err != nil {
		if err == errors.ErrNotFound {
			utils.Error(w, http.StatusNotFound, "registration not found")
			return
		}
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, reg)
}

func UpdateRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/registrations/")
	var req model.Registration
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.ID = id
	if err := RegistrationUC.Update(&req); err != nil {
		logger.Error("update registration: %v", err)
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, req)
}

func DeleteRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/registrations/")
	if err := RegistrationUC.Delete(id); err != nil {
		logger.Error("delete registration: %v", err)
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, map[string]string{"result": "deleted"})
}

func GetAllRegistrationsHandler(w http.ResponseWriter, r *http.Request) {
	regs, err := RegistrationUC.QueryRepo.GetAll()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, regs)
}
