package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"zahne/patient"
)

// Interface guard - ensures Repository implements patient.PatientRepository at compile time.
var _ patient.PatientRepository = (*Repository)(nil)

// NewRepository creates a new Repository backed by the given sql.DB.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Repository is the PostgreSQL implementation of patient.PatientRepository.
type Repository struct {
	db *sql.DB
}

func (r *Repository) Create(ctx context.Context, input patient.CreatePatientInput) (patient.Patient, error) {
	query := `
		INSERT INTO patients (name, cpf, phone, email, date_of_birth)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, cpf, phone, email, date_of_birth, created_at, updated_at
	`
	var p patient.Patient
	err := r.db.QueryRowContext(ctx, query,
		input.Name,
		input.CPF,
		input.Phone,
		input.Email,
		input.DateOfBirth,
	).Scan(
		&p.ID,
		&p.Name,
		&p.CPF,
		&p.Phone,
		&p.Email,
		&p.DateOfBirth,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return patient.Patient{}, fmt.Errorf("failed to create patient: %w", err)
	}
	return p, nil
}

func (r *Repository) GetByID(ctx context.Context, id int) (patient.Patient, error) {
	query := `
		SELECT id, name, cpf, phone, email, date_of_birth, created_at, updated_at
		FROM patients
		WHERE id = $1
	`
	var p patient.Patient
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.Name,
		&p.CPF,
		&p.Phone,
		&p.Email,
		&p.DateOfBirth,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return patient.Patient{}, fmt.Errorf("%w", patient.ErrPatientNotFound)
		}
		return patient.Patient{}, fmt.Errorf("failed to get patient: %w", err)
	}
	return p, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]patient.Patient, error) {
	query := `
		SELECT id, name, cpf, phone, email, date_of_birth, created_at, updated_at
		FROM patients
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get patients: %w", err)
	}
	defer rows.Close()

	patients := make([]patient.Patient, 0)
	for rows.Next() {
		var p patient.Patient
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.CPF,
			&p.Phone,
			&p.Email,
			&p.DateOfBirth,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan patient: %w", err)
		}
		patients = append(patients, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating patients: %w", err)
	}

	return patients, nil
}

func (r *Repository) Update(ctx context.Context, input patient.UpdatePatientInput) (patient.Patient, error) {
	query := `
		UPDATE patients
		SET name = $2, cpf = $3, phone = $4, email = $5, date_of_birth = $6, updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, cpf, phone, email, date_of_birth, created_at, updated_at
	`
	var p patient.Patient
	err := r.db.QueryRowContext(ctx, query,
		input.ID,
		input.Name,
		input.CPF,
		input.Phone,
		input.Email,
		input.DateOfBirth,
	).Scan(
		&p.ID,
		&p.Name,
		&p.CPF,
		&p.Phone,
		&p.Email,
		&p.DateOfBirth,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return patient.Patient{}, fmt.Errorf("%w", patient.ErrPatientNotFound)
		}
		return patient.Patient{}, fmt.Errorf("failed to update patient: %w", err)
	}
	return p, nil
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM patients WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete patient: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check delete result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%w", patient.ErrPatientNotFound)
	}
	return nil
}
