package ginhttp

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"zahne/internal/postgres"
	"zahne/patient"
	patientpostgres "zahne/patient/postgres"
)

// NewRouter creates and configures a Gin router with CORS, health check, and API routes.
// Returns an http.Handler so the caller remains decoupled from Gin internals.
func NewRouter(dbService postgres.Service) http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// K8S Probe
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, dbService.Health())
	})

	db := dbService.GetDB()

	patientRepo := patientpostgres.NewRepository(db)
	patientUseCase := patient.NewService(patientRepo)
	patientHandler := NewPatientHandler(patientUseCase)

	apiV1 := r.Group("/api/v1")
	{
		apiV1.POST("/patient", patientHandler.Create)
		apiV1.GET("/patient/:id", patientHandler.Get)
		apiV1.GET("/patients", patientHandler.GetAll)
		apiV1.PUT("/patient/:id", patientHandler.Update)
		apiV1.DELETE("/patient/:id", patientHandler.Delete)
	}

	return r
}
