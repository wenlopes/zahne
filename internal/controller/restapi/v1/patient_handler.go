package v1

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"zahne/patient"

	"github.com/gin-gonic/gin"
)

type PatientHandler struct {
	patientUseCase patient.PatientUseCase
}

func NewPatientHandler(patientUseCase patient.PatientUseCase) *PatientHandler {
	return &PatientHandler{
		patientUseCase: patientUseCase,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type createPatientRequest struct {
	Name        string    `json:"name" binding:"required"`
	CPF         string    `json:"cpf" binding:"required"`
	Phone       string    `json:"phone" binding:"required"`
	Email       string    `json:"email" binding:"required,email"`
	DateOfBirth time.Time `json:"date_of_birth" binding:"required"`
}

type updatePatientRequest struct {
	Name        string    `json:"name" binding:"required"`
	CPF         string    `json:"cpf" binding:"required"`
	Phone       string    `json:"phone" binding:"required"`
	Email       string    `json:"email" binding:"required,email"`
	DateOfBirth time.Time `json:"date_of_birth" binding:"required"`
}

func (h *PatientHandler) Create(c *gin.Context) {
	var req createPatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	p, err := h.patientUseCase.CreatePatient(c.Request.Context(), patient.CreatePatientInput{
		Name:        req.Name,
		CPF:         req.CPF,
		Phone:       req.Phone,
		Email:       req.Email,
		DateOfBirth: req.DateOfBirth,
	})
	if err != nil {
		slog.Error("failed to create patient", "err", err, slog.String("operation", "create_patient"))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, p)
}

func (h *PatientHandler) Get(c *gin.Context) {
	rawID := c.Param("id")
	if rawID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id is required"})
		return
	}

	id, err := strconv.Atoi(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id must be a number"})
		return
	}

	p, err := h.patientUseCase.GetPatient(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, p)
}

func (h *PatientHandler) GetAll(c *gin.Context) {
	patients, err := h.patientUseCase.GetAllPatients(c.Request.Context())
	if err != nil {
		slog.Error("failed to get all patients", "err", err, slog.String("operation", "get_all_patients"))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, patients)
}

func (h *PatientHandler) Update(c *gin.Context) {
	rawID := c.Param("id")
	if rawID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id is required"})
		return
	}

	id, err := strconv.Atoi(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id must be a number"})
		return
	}

	var req updatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	p, err := h.patientUseCase.UpdatePatient(c.Request.Context(), patient.UpdatePatientInput{
		ID:          id,
		Name:        req.Name,
		CPF:         req.CPF,
		Phone:       req.Phone,
		Email:       req.Email,
		DateOfBirth: req.DateOfBirth,
	})
	if err != nil {
		slog.Error("failed to update patient", "err", err, slog.Int("patient_id", id), slog.String("operation", "update_patient"))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, p)
}

func (h *PatientHandler) Delete(c *gin.Context) {
	rawID := c.Param("id")
	if rawID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id is required"})
		return
	}

	id, err := strconv.Atoi(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id must be a number"})
		return
	}

	if err := h.patientUseCase.DeletePatient(c.Request.Context(), id); err != nil {
		slog.Error("failed to delete patient", "err", err, slog.Int("patient_id", id), slog.String("operation", "delete_patient"))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
