package result

import "fmt"

type ServiceInterface interface {
	CreateNewResult(input CreateResultsInput) (ResultsEntity, error)
	GetAllResults(jenjang string) ([]Results, error)
	DeleteResults(id int) (bool, error)
	UpdateResutls(id int, lulus bool) (ResultsEntity, error)
}

type Service struct {
	repository RepositoryInterface
}

func NewService(repository RepositoryInterface) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) DeleteResults(id int) (bool, error) {
	deleteResult, err := s.repository.Delete(id)
	if err != nil {
		return deleteResult, err
	}
	return deleteResult, nil
}

func (s *Service) GetAllResults(jenjang string) ([]Results, error) {

	getAllResults, err := s.repository.FindAll(jenjang)
	fmt.Print(err)
	if err != nil {
		return getAllResults, err
	}
	return getAllResults, nil
}

func (s *Service) CreateNewResult(input CreateResultsInput) (ResultsEntity, error) {

	var results ResultsEntity

	resultChan := make(chan ResultsEntity)
	errorChan := make(chan error)

	go func() {
		results.UUID = input.UUID
		results.Data = input.Data
		results.Keterangan = input.Keterangan
		results.Score = input.Score
		newResult, err := s.repository.Save(results)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- newResult

	}()
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return ResultsEntity{}, err
	}
}

func (s *Service) UpdateResutls(id int, lulus bool) (ResultsEntity, error) {
	updateResults, err := s.repository.Update(id, lulus)
	if err != nil {
		return updateResults, err
	}
	return updateResults, nil
}
