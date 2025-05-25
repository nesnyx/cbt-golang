package institution

import (
	"errors"
	"gorm.io/gorm"
)

type ServiceInterface interface {
	CreateNewInstitution(institution CreateNewInstitutionInput) (Institution, error)
	UpdateInstitution(input UpdateInstitutionInput) (InstitutionResponse, error)
	DeleteInstitution(input DeleteInstitutionInput) (bool, error)
	GetAllInsitutions() ([]Institution, error)
}

type Service struct {
	repository RepositoryInterface
}

func NewService(repository RepositoryInterface) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetAllInsitutions() ([]Institution, error) {
	institutions, err := s.repository.FindAll()
	if err != nil {
		return institutions, err
	}
	return institutions, nil
}

func (s *Service) DeleteInstitution(input DeleteInstitutionInput) (bool, error) {
	_, errGetInstitution := s.repository.FindByID(input.ID)
	if errors.Is(errGetInstitution, gorm.ErrRecordNotFound) {
		return false, errors.New("institution not found")
	}
	deleteInstitution, err := s.repository.Delete(input.ID)
	if err != nil {
		return false, err
	}
	return deleteInstitution, nil
}

func (s *Service) CreateNewInstitution(input CreateNewInstitutionInput) (Institution, error) {
	var institution Institution

	institution.Name = input.Name
	institution.Address = input.Address
	newInstitution, err := s.repository.Save(institution)
	if err != nil {
		return newInstitution, err
	}
	return newInstitution, nil
}

func (s *Service) UpdateInstitution(input UpdateInstitutionInput) (InstitutionResponse, error) {
	var institution Institution
	//validate exist or not
	getInstitution, errGetInstitution := s.repository.FindByID(input.ID)
	if errors.Is(errGetInstitution, gorm.ErrRecordNotFound) {
		return getInstitution, errors.New("institution not found")
	}
	institution.Name = input.Name
	institution.Address = input.Address
	updateInstitution, err := s.repository.Update(input.ID, institution)
	if err != nil {
		return updateInstitution, err
	}
	return updateInstitution, nil
}
