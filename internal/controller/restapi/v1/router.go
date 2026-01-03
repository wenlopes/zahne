package v1

import (
	"github.com/gin-gonic/gin"
	"zahne/internal/repository/persistent"
	"zahne/internal/usecase/patient"
	"zahne/pkg/postgres"
)

func NewRouter(r *gin.Engine, dbService postgres.Service) {
	// Get database connection
	db := dbService.GetDB()

	// Initialize repositories
	patientRepo := persistent.NewPatientRepository(db)

	// Initialize use cases
	patientUseCase := patient.New(patientRepo)

	// Initialize handlers
	dentistHandler := NewDentistHandler()
	patientHandler := NewPatientHandler(patientUseCase)

	apiV1 := r.Group("/api/v1")
	{
		// Hello World route
		apiV1.GET("/hello", HelloWorldHandler)

		// Dentist routes
		apiV1.POST("/dentist", dentistHandler.Create)
		apiV1.GET("/dentist/:id", dentistHandler.Get)

		// Patient routes
		apiV1.POST("/patient", patientHandler.Create)
		apiV1.GET("/patient/:id", patientHandler.Get)
		apiV1.GET("/patients", patientHandler.GetAll)
		apiV1.PUT("/patient/:id", patientHandler.Update)
		apiV1.DELETE("/patient/:id", patientHandler.Delete)
	}
}
