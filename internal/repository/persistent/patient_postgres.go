package persistent

import (
	"database/sql"
	"fmt"

	"zahne/internal/entity"
)

type PatientRepository struct {
	db *sql.DB
}

func NewPatientRepository(db *sql.DB) *PatientRepository {
	return &PatientRepository{db: db}
}

func (r *PatientRepository) Create(patient entity.Patient) error {
	query := `
		INSERT INTO patients (id, name, cpf, phone, email, date_of_birth, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(query,
		patient.ID,
		patient.Name,
		patient.CPF,
		patient.Phone,
		patient.Email,
		patient.DateOfBirth,
		patient.CreatedAt,
		patient.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create patient: %w", err)
	}
	return nil
}

func (r *PatientRepository) GetByID(id string) (entity.Patient, error) {
	query := `
		SELECT id, name, cpf, phone, email, date_of_birth, created_at, updated_at
		FROM patients
		WHERE id = $1
	`
	var patient entity.Patient
	err := r.db.QueryRow(query, id).Scan(
		&patient.ID,
		&patient.Name,
		&patient.CPF,
		&patient.Phone,
		&patient.Email,
		&patient.DateOfBirth,
		&patient.CreatedAt,
		&patient.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.Patient{}, fmt.Errorf("patient not found")
		}
		return entity.Patient{}, fmt.Errorf("failed to get patient: %w", err)
	}
	return patient, nil
}

func (r *PatientRepository) GetAll() ([]entity.Patient, error) {
	query := `
		SELECT id, name, cpf, phone, email, date_of_birth, created_at, updated_at
		FROM patients
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get patients: %w", err)
	}
	defer rows.Close()

	patients := make([]entity.Patient, 0)
	for rows.Next() {
		var patient entity.Patient
		err := rows.Scan(
			&patient.ID,
			&patient.Name,
			&patient.CPF,
			&patient.Phone,
			&patient.Email,
			&patient.DateOfBirth,
			&patient.CreatedAt,
			&patient.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan patient: %w", err)
		}
		patients = append(patients, patient)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating patients: %w", err)
	}

	return patients, nil
}

func (r *PatientRepository) Update(patient entity.Patient) error {
	query := `
		UPDATE patients
		SET name = $2, cpf = $3, phone = $4, email = $5, date_of_birth = $6, updated_at = $7
		WHERE id = $1
	`
	result, err := r.db.Exec(query,
		patient.ID,
		patient.Name,
		patient.CPF,
		patient.Phone,
		patient.Email,
		patient.DateOfBirth,
		patient.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update patient: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("patient not found")
	}

	return nil
}

func (r *PatientRepository) Delete(id string) error {
	query := `DELETE FROM patients WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete patient: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("patient not found")
	}

	return nil
}
