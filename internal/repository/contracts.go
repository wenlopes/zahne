package repository

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
		GetByID(id int)
	}
)
