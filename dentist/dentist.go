package dentist

//go:generate mockgen -source=dentist.go -destination=mock/mock_dentist.go -package=mock

import "context"

// Dentist represents the dentist domain entity.
type Dentist struct {
	Name string
	CRO  string
}

// DentistRepository defines the persistence contract for the Dentist domain.
type DentistRepository interface {
	GetByID(ctx context.Context, id int) (Dentist, error)
	Store(ctx context.Context, dentist Dentist) (Dentist, error)
}
