package patient

import (
	"database/sql"
	"testing"
	"time"
)

// TestPatientData represents patient data used in tests
type TestPatientData struct {
	ID          int
	Name        string
	CPF         string
	Phone       string
	Email       string
	DateOfBirth time.Time
}

// InsertTestPatient inserts a test patient into the database and sets the auto-generated ID on the data.
func InsertTestPatient(t *testing.T, db *sql.DB, data *TestPatientData) {
	t.Helper()
	query := `
		INSERT INTO patients (name, cpf, phone, email, date_of_birth)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := db.QueryRow(query,
		data.Name,
		data.CPF,
		data.Phone,
		data.Email,
		data.DateOfBirth,
	).Scan(&data.ID)
	if err != nil {
		t.Fatalf("failed to insert test patient: %v", err)
	}
}
