package postgres

import (
	"context"
	"errors"

	"zahne/dentist"
)

// Interface guard - ensures Repository implements dentist.DentistRepository at compile time.
var _ dentist.DentistRepository = (*Repository)(nil)

// NewRepository creates a new Repository for the dentist domain.
func NewRepository() *Repository {
	return &Repository{}
}

// Repository is the PostgreSQL stub implementation of dentist.DentistRepository.
type Repository struct{}

func (r *Repository) GetByID(ctx context.Context, id int) (dentist.Dentist, error) {
	return dentist.Dentist{}, errors.New("not implemented")
}

func (r *Repository) Store(ctx context.Context, d dentist.Dentist) (dentist.Dentist, error) {
	return dentist.Dentist{}, errors.New("not implemented")
}
