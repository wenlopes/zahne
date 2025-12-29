package v1

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.Engine) {
	dentistHandler := NewDentistHandler()

	apiV1 := r.Group("/api/v1")
	{
		// Hello World route
		apiV1.GET("/hello", HelloWorldHandler)

		// Dentist routes
		apiV1.POST("/dentist", dentistHandler.Create)
		apiV1.GET("/dentist/:id", dentistHandler.Get)
	}
}
