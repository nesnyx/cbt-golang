package soal

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceInterface interface {
	CreateSoal(input CreateSoalInput) (Soal, error)
	GetAllSoal() ([]Soal, error)
	GetByKodeSoal(input GetSoalByKodeSoalInput) (Soal, error)
	DeleteSoal(input GetSoalByIDInput) (bool, error)
	CreateTmpSoal(input CreateTmpSoalInputan) (TmpSoal, error)
	GetAllTmpSoal() ([]TmpSoal, error)
	UpdateTmpSoal(input UpdateTmpSoalOptionsInputan) (bool, error)
	GetCorrectAnswerByQuestionID(questionId string) (*TmpSoal, error)
	GetAllTmpSoalByType(typeSoal string) ([]TmpSoal, error)
	GetAllTmpSoalByType2(typeSoal string) ([]TmpSoal, error)
}

type Service struct {
	repository RepositoryInterface
}

func NewService(repository RepositoryInterface) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetAllTmpSoalByType(typeSoal string) ([]TmpSoal, error) {
	soal, err := s.repository.FindAllTmpSoalByType(typeSoal)
	if err != nil {
		return soal, err
	}
	return soal, nil
}
func (s *Service) GetAllTmpSoalByType2(typeSoal string) ([]TmpSoal, error) {
	soal, err := s.repository.FindAllTmpSoalByTypeStudent(typeSoal)
	if err != nil {
		return soal, err
	}
	return soal, nil
}

func (s *Service) GetAllTmpSoal() ([]TmpSoal, error) {
	soal, err := s.repository.FindAllTmpAnswer()
	if err != nil {
		return soal, err
	}
	return soal, nil
}

func (s *Service) GetCorrectAnswerByQuestionID(questionId string) (*TmpSoal, error) {
	correctAnswer, err := s.repository.FindCorrectAnswerByQuestionID(questionId)
	if err != nil {
		return correctAnswer, err
	}
	fmt.Println(correctAnswer)
	return correctAnswer, nil
}

func (s *Service) UpdateTmpSoal(input UpdateTmpSoalOptionsInputan) (bool, error) {
	var soal TmpSoal
	_, err := s.repository.FindByID(uint(input.ID))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Printf(fmt.Sprintf("id With %v not found", input.ID))
		return false, errors.New(fmt.Sprintf("id With %v not found", input.ID))
	}

	optionsJSON, _ := json.Marshal(input.Updates)

	soal.Options = optionsJSON
	updateTmpSoal, err := s.repository.Update(uint(input.ID), input.QuestionText, input.QuestionPicture, input.CorrectAnswer, soal)
	if err != nil {
		return false, err
	}
	fmt.Println(updateTmpSoal)

	return true, nil
}

func (s *Service) CreateTmpSoal(input CreateTmpSoalInputan) (TmpSoal, error) {
	var soal TmpSoal

	questionId := uuid.New().String()[:10]
	questionJson, _ := json.Marshal(input.Pertanyaan)
	optionsJson, _ := json.Marshal(input.Options)
	soal.QuestionID = questionId
	soal.CorrectAnswer = input.CorrectAnswer
	soal.Pertanyaan = questionJson
	soal.Options = optionsJson
	soal.Type = input.Type
	createTmpSoal, err := s.repository.SaveTmpSoal(soal)
	if err != nil {
		return createTmpSoal, err
	}
	return createTmpSoal, nil
}

func (s *Service) GetAllSoal() ([]Soal, error) {
	soal, err := s.repository.FindAll()
	if err != nil {
		return soal, err
	}
	return soal, nil
}

func (s *Service) DeleteSoal(input GetSoalByIDInput) (bool, error) {
	soal, err := s.repository.Delete(input.ID)
	if err != nil {
		return false, err
	}
	return soal, err
}

func (s *Service) CreateSoal(input CreateSoalInput) (Soal, error) {
	var soal Soal
	questionsJson, err := json.Marshal(input.Questions)
	if err != nil {
		return soal, err
	}
	soal.KodeSoal = input.KodeSoal
	soal.Waktu = input.Waktu
	soal.Questions = questionsJson
	newSoal, err := s.repository.Save(soal)
	if err != nil {
		return newSoal, err
	}
	return newSoal, nil
}

func (s *Service) GetByKodeSoal(input GetSoalByKodeSoalInput) (Soal, error) {
	soal, err := s.repository.FindByKodeSoal(input.KodeSoal)
	if err != nil {
		return soal, err
	}
	return soal, nil
}
