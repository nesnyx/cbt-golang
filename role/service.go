package role

type Service interface {
	GetByID(id int) (Roles, error)
}

type service struct {
	repository RepositoryInterface
}

func NewService(repository RepositoryInterface) *service {
	return &service{repository}
}

func (s *service) GetByID(id int) (Roles, error) {
	role, err := s.repository.FindByID(id)
	if err != nil {
		return role, err
	}
	return role, nil
}
