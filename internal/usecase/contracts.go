package usecase

import (
	"time"

	"zahne/internal/entity"
)

type PatientUseCase interface {
	CreatePatient(name, cpf, phone, email string, dateOfBirth time.Time) (entity.Patient, error)
	GetPatientByID(id string) (entity.Patient, error)
	GetAllPatients() ([]entity.Patient, error)
	UpdatePatient(id, name, cpf, phone, email string, dateOfBirth time.Time) (entity.Patient, error)
	DeletePatient(id string) error
}
