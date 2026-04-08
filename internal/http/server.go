package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

// New creates a configured HTTP server with the provided handler and standard timeouts.
func New(handler http.Handler) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
