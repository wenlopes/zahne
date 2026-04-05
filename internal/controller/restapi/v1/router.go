package v1

import (
	"zahne/patient"
	patientpostgres "zahne/patient/postgres"
	"zahne/internal/postgres"

	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.Engine, dbService postgres.Service) {
	// Get database connection
	db := dbService.GetDB()

	// Initialize repositories
	patientRepo := patientpostgres.NewRepository(db)

	// Initialize use cases
	patientUseCase := patient.NewService(patientRepo)

	// Initialize handlers
	patientHandler := NewPatientHandler(patientUseCase)

	apiV1 := r.Group("/api/v1")
	{
		// Patient routes
		apiV1.POST("/patient", patientHandler.Create)
		apiV1.GET("/patient/:id", patientHandler.Get)
		apiV1.GET("/patients", patientHandler.GetAll)
		apiV1.PUT("/patient/:id", patientHandler.Update)
		apiV1.DELETE("/patient/:id", patientHandler.Delete)
	}
}
