package patient

//go:generate mockgen -source=patient.go -destination=mock/mock_patient.go -package=mock

import (
	"context"
	"errors"
	"time"
)

// ErrPatientNotFound is returned when a Patient is not found.
var ErrPatientNotFound = errors.New("patient not found")

// Patient represents the patient domain entity.
type Patient struct {
	ID          int
	Name        string
	CPF         string
	Phone       string
	Email       string
	DateOfBirth time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// PatientRepository defines the persistence contract for the Patient domain.
type PatientRepository interface {
	Create(ctx context.Context, input CreatePatientInput) (Patient, error)
	GetByID(ctx context.Context, id int) (Patient, error)
	GetAll(ctx context.Context) ([]Patient, error)
	Update(ctx context.Context, input UpdatePatientInput) (Patient, error)
	Delete(ctx context.Context, id int) error
}

// PatientUseCase defines the business contract for the Patient domain.
type PatientUseCase interface {
	CreatePatient(ctx context.Context, input CreatePatientInput) (Patient, error)
	GetPatient(ctx context.Context, id int) (Patient, error)
	GetAllPatients(ctx context.Context) ([]Patient, error)
	UpdatePatient(ctx context.Context, input UpdatePatientInput) (Patient, error)
	DeletePatient(ctx context.Context, id int) error
}

// CreatePatientInput contains the data required to create a Patient.
type CreatePatientInput struct {
	Name        string
	CPF         string
	Phone       string
	Email       string
	DateOfBirth time.Time
}

// UpdatePatientInput contains the data required to update a Patient.
type UpdatePatientInput struct {
	ID          int
	Name        string
	CPF         string
	Phone       string
	Email       string
	DateOfBirth time.Time
}
