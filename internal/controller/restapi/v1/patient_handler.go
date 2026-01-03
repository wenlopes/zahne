package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"zahne/internal/controller/restapi/v1/request"
	"zahne/internal/usecase"
)

type PatientHandler struct {
	patientUseCase usecase.PatientUseCase
}

func NewPatientHandler(patientUseCase usecase.PatientUseCase) *PatientHandler {
	return &PatientHandler{
		patientUseCase: patientUseCase,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *PatientHandler) Create(c *gin.Context) {
	var req request.CreatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	patient, err := h.patientUseCase.CreatePatient(
		req.Name,
		req.CPF,
		req.Phone,
		req.Email,
		req.DateOfBirth,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, patient)
}

func (h *PatientHandler) Get(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id is required"})
		return
	}

	patient, err := h.patientUseCase.GetPatientByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, patient)
}

func (h *PatientHandler) GetAll(c *gin.Context) {
	patients, err := h.patientUseCase.GetAllPatients()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, patients)
}

func (h *PatientHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id is required"})
		return
	}

	var req request.UpdatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	patient, err := h.patientUseCase.UpdatePatient(
		id,
		req.Name,
		req.CPF,
		req.Phone,
		req.Email,
		req.DateOfBirth,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, patient)
}

func (h *PatientHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "id is required"})
		return
	}

	err := h.patientUseCase.DeletePatient(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
