package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DentistHandler struct{}

func NewDentistHandler() *DentistHandler {
	return &DentistHandler{}
}

func (h *DentistHandler) Create(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Create Dentist request received"
	c.JSON(http.StatusOK, resp)
}

func (h *DentistHandler) Get(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Get Dentist request received"
	c.JSON(http.StatusOK, resp)
}
