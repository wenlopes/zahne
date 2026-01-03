package integration_tests

import (
	"database/sql"
	"testing"
	"time"
)

// testPatientData represents patient data used in tests
type testPatientData struct {
	ID          string
	Name        string
	CPF         string
	Phone       string
	Email       string
	DateOfBirth time.Time
}

// insertTestPatient inserts a test patient into the database
func insertTestPatient(t *testing.T, db *sql.DB, patient testPatientData) {
	query := `
		INSERT INTO patients (id, name, cpf, phone, email, date_of_birth, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	now := time.Now()
	_, err := db.Exec(query,
		patient.ID,
		patient.Name,
		patient.CPF,
		patient.Phone,
		patient.Email,
		patient.DateOfBirth,
		now,
		now,
	)
	if err != nil {
		t.Fatalf("failed to insert test patient: %v", err)
	}
}
