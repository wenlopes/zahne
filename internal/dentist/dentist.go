package dentist

type Dentist struct {
	Name string
	CRO  string
}

type Reader interface {
	GetByCRO(cro string) (Dentist, error)
}

type Writer interface {
	Save(dentist Dentist) error
}

type Repository struct {
	Reader
	Writer
}

type UseCase interface {
	RegisterDentist(dentist Dentist) error
	FindDentistByCRO(cro string) (Dentist, error)
}
