package patient

import (
	"time"

	"github.com/google/uuid"
	"zahne/internal/entity"
	"zahne/internal/repository"
	"zahne/internal/usecase"
)

type patientUseCase struct {
	patientRepo repository.Patient
}

func New(patientRepo repository.Patient) usecase.PatientUseCase {
	return &patientUseCase{
		patientRepo: patientRepo,
	}
}

func (uc *patientUseCase) CreatePatient(name, cpf, phone, email string, dateOfBirth time.Time) (entity.Patient, error) {
	patient := entity.Patient{
		ID:          uuid.New().String(),
		Name:        name,
		CPF:         cpf,
		Phone:       phone,
		Email:       email,
		DateOfBirth: dateOfBirth,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := uc.patientRepo.Create(patient)
	if err != nil {
		return entity.Patient{}, err
	}

	return patient, nil
}

func (uc *patientUseCase) GetPatientByID(id string) (entity.Patient, error) {
	return uc.patientRepo.GetByID(id)
}

func (uc *patientUseCase) GetAllPatients() ([]entity.Patient, error) {
	return uc.patientRepo.GetAll()
}

func (uc *patientUseCase) UpdatePatient(id, name, cpf, phone, email string, dateOfBirth time.Time) (entity.Patient, error) {
	patient := entity.Patient{
		ID:          id,
		Name:        name,
		CPF:         cpf,
		Phone:       phone,
		Email:       email,
		DateOfBirth: dateOfBirth,
		UpdatedAt:   time.Now(),
	}

	err := uc.patientRepo.Update(patient)
	if err != nil {
		return entity.Patient{}, err
	}

	return uc.patientRepo.GetByID(id)
}

func (uc *patientUseCase) DeletePatient(id string) error {
	return uc.patientRepo.Delete(id)
}
