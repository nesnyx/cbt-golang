package tmpAnswer

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type ServiceInterface interface {
	Save(input CreateTmpAnswerInput, uuid string) (TmpAnswer, error)
	GetAll() ([]TmpAnswer, error)
	FindByQuestionID(input GetTmpSoalInput) (TmpAnswer, error)
	GetAnswerByUUID(uuid string) ([]TmpAnswer, error)
}

type Service struct {
	repository RepositoryInterface
}

func NewService(repository RepositoryInterface) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetAnswerByUUID(uuid string) ([]TmpAnswer, error) {
	answers, err := s.repository.FindAnswersByUUID(uuid)
	if err != nil {
		return answers, err
	}
	return answers, nil
}

func (s *Service) FindByQuestionID(input GetTmpSoalInput) (TmpAnswer, error) {

	tmpAnswer, err := s.repository.FindByQuestionID(input.QuestionsID, input.PesertaUUID)
	fmt.Println(tmpAnswer.Answers)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("Data tidak ditemukan")
	} else {
		fmt.Println("Data ditemukan")
	}
	//if err != nil {
	//	return tmpAnswer, err
	//}
	return tmpAnswer, nil
}

func (s *Service) Save(input CreateTmpAnswerInput, uuid string) (TmpAnswer, error) {
	var sem = make(chan struct{}, 10)
	sem <- struct{}{}
	defer func() { <-sem }()

	var tmpAnswer TmpAnswer

	resultChan := make(chan TmpAnswer, 1)
	errChan := make(chan error, 1)

	go func() {
		getTmpAnswer, err := s.repository.FindByQuestionID(input.QuestionsID, uuid)
		fmt.Println(getTmpAnswer)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("Data tidak ditemukan")
			tmpAnswer.PesertaUUID = uuid

			tmpAnswer.Answers = input.Answers
			tmpAnswer.QuestionsID = input.QuestionsID
			tmpAnswer.Time = int(time.Now().Unix())
			newTmpAnswer, err := s.repository.Save(tmpAnswer)
			if err != nil {
				errChan <- err
				return
			}
			resultChan <- newTmpAnswer
			return
		}
		fmt.Println("Data ditemukan")
		updateAnswer, err := s.repository.Update(input.QuestionsID, uuid, input.Answers, time.Now().Unix())
		tmpAnswer.PesertaUUID = uuid

		tmpAnswer.Answers = input.Answers
		tmpAnswer.QuestionsID = input.QuestionsID
		fmt.Println(updateAnswer)
		resultChan <- tmpAnswer
		return
	}()
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errChan:
		return TmpAnswer{}, err
	}
}

func (s *Service) GetAll() ([]TmpAnswer, error) {
	tmpAnswer, err := s.repository.FindAll()
	fmt.Println(len(tmpAnswer))
	if err != nil {
		return tmpAnswer, err
	}
	return tmpAnswer, nil
}
