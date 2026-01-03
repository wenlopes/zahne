package patient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"zahne/integration-tests/shared"
	"zahne/internal/entity"
)

func TestGetPatientByID_Success(t *testing.T) {
	db, cleanup := shared.SetupTestDB(t)
	defer cleanup()

	router := setupRouter(db)

	// Create test patient data
	testPatient := TestPatientData{
		ID:          uuid.New().String(),
		Name:        "John Doe",
		CPF:         "12345678901",
		Phone:       "11999999999",
		Email:       "john.doe@example.com",
		DateOfBirth: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
	}

	InsertTestPatient(t, db, testPatient)

	// Make GET request
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/patient/%s", testPatient.ID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response entity.Patient
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, testPatient.ID, response.ID)
	assert.Equal(t, testPatient.Name, response.Name)
	assert.Equal(t, testPatient.CPF, response.CPF)
	assert.Equal(t, testPatient.Phone, response.Phone)
	assert.Equal(t, testPatient.Email, response.Email)
	assert.Equal(t, testPatient.DateOfBirth.UTC().Format(time.RFC3339), response.DateOfBirth.UTC().Format(time.RFC3339))
}

func TestGetPatientByID_NotFound(t *testing.T) {
	db, cleanup := shared.SetupTestDB(t)
	defer cleanup()

	router := setupRouter(db)

	// Use a random UUID that doesn't exist in the database
	nonExistentID := uuid.New().String()

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/patient/%s", nonExistentID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response, "error")
}

func TestGetPatientByID_EmptyID(t *testing.T) {
	db, cleanup := shared.SetupTestDB(t)
	defer cleanup()

	router := setupRouter(db)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/patient/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// When the ID parameter is empty, Gin will redirect to /api/v1/patient
	// which results in 301 redirect or the request goes to a different route
	// We expect either 301 (redirect) or 404 (not found)
	assert.Contains(t, []int{http.StatusNotFound, http.StatusMovedPermanently}, w.Code)
}
