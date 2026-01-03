package patient

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	v1 "zahne/internal/controller/restapi/v1"
	"zahne/internal/repository/persistent"
	"zahne/internal/usecase/patient"
)

// setupRouter creates a Gin router with patient endpoints configured
func setupRouter(db *sql.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	patientRepo := persistent.NewPatientRepository(db)
	patientUseCase := patient.New(patientRepo)
	patientHandler := v1.NewPatientHandler(patientUseCase)

	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/patient", patientHandler.Create)
		apiV1.GET("/patient/:id", patientHandler.Get)
		apiV1.GET("/patients", patientHandler.GetAll)
		apiV1.PUT("/patient/:id", patientHandler.Update)
		apiV1.DELETE("/patient/:id", patientHandler.Delete)
	}

	return router
}
