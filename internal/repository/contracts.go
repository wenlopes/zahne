package repository

import "zahne/internal/entity"

type (
	// Dentist
	DentistReader interface {
		GetByID(id int)
	}

	DentistWriter interface {
		Store()
	}

	Dentist interface {
		DentistReader
		DentistWriter
	}
)

type (
	// Patient
	PatientReader interface {
		GetByID(id string) (entity.Patient, error)
		GetAll() ([]entity.Patient, error)
	}

	PatientWriter interface {
		Create(patient entity.Patient) error
		Update(patient entity.Patient) error
		Delete(id string) error
	}

	Patient interface {
		PatientReader
		PatientWriter
	}
)
