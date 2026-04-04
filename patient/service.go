package patient

import "context"

// Interface guard — ensures patientService implements PatientUseCase at compile time.
var _ PatientUseCase = (*patientService)(nil)

// NewService creates a new patientService with the given PatientRepository.
func NewService(repo PatientRepository) *patientService {
	return &patientService{repo: repo}
}

type patientService struct {
	repo PatientRepository
}

func (s *patientService) CreatePatient(ctx context.Context, input CreatePatientInput) (Patient, error) {
	return s.repo.Create(ctx, input)
}

func (s *patientService) GetPatient(ctx context.Context, id int) (Patient, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *patientService) GetAllPatients(ctx context.Context) ([]Patient, error) {
	return s.repo.GetAll(ctx)
}

func (s *patientService) UpdatePatient(ctx context.Context, input UpdatePatientInput) (Patient, error) {
	return s.repo.Update(ctx, input)
}

func (s *patientService) DeletePatient(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
