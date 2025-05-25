package peserta

import (
	"cbt/auth"
	"cbt/result"
	"cbt/soal"
	"cbt/tmpAnswer"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	bcrypt2 "github.com/aerospike/aerospike-client-go/pkg/bcrypt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Service interface {
	Create(input CreatePesertaInput) (Peserta, error)
	CreateNewPeserta(input CreateNewPesertaInput, jenjang string) (PesertaEntity, error)
	CreateNewPeserta2(input CreateNewPesertaInput2, jenjang string) (PesertaEntity, error)
	CreateNewAdmin(input CreateNewAdminInput) (CreateNewAdminResponse, error)
	GetAll(jenjang string) ([]GetAllPesertaInput, error)
	UpdateByUUID(inputUUID GetUUIDPesertaInput, input UpdatePesertaAnswers) (Peserta, error)
	GetByUUID(uuid string) (GetPesertaTokenInput, error)
	Login(input LoginInput, lastLogin int64) (Peserta, error)
	Logout(c *gin.Context) (bool, error)
	UpdateStartTime(uuid string, isStarted bool, isFinish bool) (Peserta, error)
	UpdateStartTimeFinish(uuid string, isStarted bool, isFinish bool, keterangan json.RawMessage) (Peserta, error)
	UpdatePeserta(uuid string, input UpdatePesertaInput) (UpdatePesertaInput, error)
	DeletePeserta(uuid string) (bool, error)
}

type service struct {
	repository       Repository
	authService      auth.Service
	tmpAnswerService tmpAnswer.ServiceInterface
	tmpSoalService   soal.ServiceInterface
	resultService    result.ServiceInterface
}

func NewService(repository Repository, authService auth.Service, tmpAnswerService tmpAnswer.ServiceInterface, tmpSoalService soal.ServiceInterface, resultService result.ServiceInterface) *service {
	return &service{repository, authService, tmpAnswerService, tmpSoalService, resultService}
}

func (s *service) UpdateStartTime(uuid string, isStarted bool, isFinish bool) (Peserta, error) {
	timeNow := time.Now().Unix()
	updateStartTime, err := s.repository.UpdateStartTimeByUUID(uuid, timeNow, isStarted, isFinish)

	if err != nil {
		return Peserta{}, err
	}
	return Peserta{
		UUID:      updateStartTime.UUID,
		StartTime: updateStartTime.StartTime,
	}, nil

}

func (s *service) UpdateStartTimeFinish(uuid string, isStarted bool, isFinish bool, keterangan json.RawMessage) (Peserta, error) {

	getPeserta, err := s.repository.FindByUUID(uuid)
	if err != nil {
		return Peserta{}, err
	}
	updateStartTime, err := s.repository.UpdateStartTimeByUUID(uuid, 0, isStarted, isFinish)
	if err != nil {
		return Peserta{}, err
	}
	tmpAnswers, err := s.tmpAnswerService.GetAnswerByUUID(uuid)
	if err != nil {
		return Peserta{}, errors.New("tmp Answer by UUID Error")
	}

	var score float64 = 0
	for _, answer := range tmpAnswers {
		tmpSoal, err := s.tmpSoalService.GetCorrectAnswerByQuestionID(answer.QuestionsID)
		if err != nil {
			return Peserta{}, err
		}
		if tmpSoal == nil {
			fmt.Println("Warning: tmpSoal is nil for QuestionID:", answer.QuestionsID)
			continue
		}
		if tmpSoal.CorrectAnswer == answer.Answers {
			score += 2.5
		} else if tmpSoal.CorrectAnswer != answer.Answers {
			score += 0
		}
	}

	data := map[string]interface{}{
		"start_time":  getPeserta.StartTime,
		"finish_time": time.Now().Unix(),
	}

	fmt.Println("Score : ", score)
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("JSON Encoding Error:", err)
		return Peserta{}, err
	}
	var input result.CreateResultsInput
	input.UUID = uuid
	input.Score = score
	input.Data = jsonData
	input.Keterangan = keterangan
	_, errResult := s.resultService.CreateNewResult(input)
	if errResult != nil {
		return Peserta{}, errors.New("apakah disini errornya?")
	}

	return Peserta{
		UUID:      updateStartTime.UUID,
		StartTime: updateStartTime.StartTime,
	}, nil

}

func (s *service) Logout(c *gin.Context) (bool, error) {

	getUUID := c.MustGet("currentUser").(GetPesertaTokenInput).UUID
	user, err := s.repository.FindByUUID(getUUID)
	if err != nil {
		return false, err
	}
	fmt.Println(user)
	_, errDeleted := s.repository.DeleteAllowedToken(user.UUID)
	_, errDeletedLastLogin := s.repository.DeleteLastLogin(user.UUID)
	if errDeletedLastLogin == nil {
		return true, nil
	}
	if errDeleted != nil {
		return false, err
	}
	return true, nil
}

func (s *service) Login(input LoginInput, lastLogin int64) (Peserta, error) {
	noPendaftaran := input.NoPendaftaran
	password := input.Password
	user, err := s.repository.FindByNoPendaftaran(noPendaftaran)

	if err != nil {
		return user, errors.New("no Pendaftaran does not match")
	}
	if user.Password != password {
		return user, errors.New("password does not match")
	}

	compareHash := bcrypt2.Match(password+user.Salt, user.Hash)
	fmt.Println(compareHash)
	if !compareHash {
		return user, errors.New("password does not match")
	}
	//newToken, _ := s.authService.GenerateToken(user.UUID)
	//if user.AllowedToken == "" {
	//	_, err := s.repository.SaveAllowedToken(newToken, user.UUID)
	//	_, errLogin := s.repository.SaveLastLogin(lastLogin, user.UUID)
	//	if errLogin != nil {
	//
	//		return user, errLogin
	//	}
	//	if err != nil {
	//
	//		return user, err
	//	}
	//	return user, nil
	//}

	//if len(user.AllowedToken) > 0 {
	//	expiredToken, err := s.authService.ValidateToken(user.AllowedToken)
	//	fmt.Println("expiredToken", expiredToken)
	//	if err != nil {
	//		_, errToken := s.repository.DeleteAllowedToken(user.UUID)
	//		_, errLastLogin := s.repository.DeleteLastLogin(user.UUID)
	//		_, errAllowedToken := s.repository.SaveAllowedToken(newToken, user.UUID)
	//		_, errLogin := s.repository.SaveLastLogin(lastLogin, user.UUID)
	//		if errToken != nil || errLastLogin != nil {
	//			return user, errors.New("something wrong with deleted Allowed Token and Last Login User")
	//		}
	//		if errLogin != nil || errAllowedToken != nil {
	//			return user, errLogin
	//		}
	//		return user, err
	//	}
	//	fmt.Println(errors.New("user Masih Login"))
	//	return user, errors.New("user Masih Login")
	//}

	return user, err
}

func (s *service) GetByUUID(uuid string) (GetPesertaTokenInput, error) {
	peserta, err := s.repository.FindByUUID(uuid)
	if err != nil {
		return GetPesertaTokenInput{}, err
	}
	return peserta, nil
}

func (s *service) UpdateByUUID(inputUUID GetUUIDPesertaInput, input UpdatePesertaAnswers) (Peserta, error) {
	peserta := Peserta{}
	answersJson, err := json.Marshal(input.Answers)
	if err != nil {
		return Peserta{}, err
	}
	peserta.Answers = answersJson
	updatedPeserta, err := s.repository.UpdateByUUID(inputUUID.UUID, peserta)
	if err != nil {
		return updatedPeserta, err
	}
	return updatedPeserta, nil
}

func (s *service) Create(input CreatePesertaInput) (Peserta, error) {
	uuidRandom := uuid.New()
	//answersJson, err := json.Marshal(input.Answers)
	randomPassword, _ := generateSecureRandomString(8)
	fmt.Println(randomPassword)
	salt, err := generateSecureRandomString(64)
	concat := randomPassword + salt
	fmt.Println(salt)

	hash, errHash := bcrypt2.Hash(concat)
	if errHash != nil {
		return Peserta{}, errHash
	}
	fmt.Println(hash)
	if err != nil {
		return Peserta{}, err
	}
	peserta := Peserta{}
	peserta.Name = input.Name
	peserta.NoPendaftaran = input.NoPendaftaran
	peserta.UUID = uuidRandom.String()
	peserta.Hash = hash
	peserta.Salt = salt
	peserta.Password = randomPassword
	peserta.AllowedToken = input.AllowedToken
	//peserta.Answers = answersJson
	peserta.StartTime = input.StartTime

	newPeserta, err := s.repository.Save(peserta)
	if err != nil {
		return newPeserta, err
	}
	return newPeserta, nil

}

func (s *service) CreateNewPeserta(input CreateNewPesertaInput, jenjang string) (PesertaEntity, error) {
	uuidRandom := uuid.New()
	randomPassword, _ := generateSecureRandomString(8)
	salt, err := generateSecureRandomString(64)
	concat := randomPassword + salt
	hash, errHash := bcrypt2.Hash(concat)
	if errHash != nil {
		return PesertaEntity{}, errHash
	}

	if err != nil {
		return PesertaEntity{}, err
	}
	identityPeserta, _ := json.Marshal(input.Name)

	peserta := PesertaEntity{}
	peserta.Name = identityPeserta
	peserta.NoPendaftaran = input.NoPendaftaran
	peserta.UUID = uuidRandom.String()
	peserta.Hash = hash
	peserta.Salt = salt
	peserta.Password = randomPassword
	peserta.Role = "peserta|" + jenjang

	newPeserta, err := s.repository.SavePeserta(peserta)
	if err != nil {
		return newPeserta, err
	}
	return newPeserta, nil

}

func (s *service) CreateNewPeserta2(input CreateNewPesertaInput2, jenjang string) (PesertaEntity, error) {
	uuidRandom := uuid.New()
	randomPassword := input.Password
	salt, err := generateSecureRandomString(64)
	concat := randomPassword + salt
	hash, errHash := bcrypt2.Hash(concat)
	if errHash != nil {
		return PesertaEntity{}, errHash
	}

	if err != nil {
		return PesertaEntity{}, err
	}
	identityPeserta, _ := json.Marshal(input.Name)

	peserta := PesertaEntity{}
	peserta.Name = identityPeserta
	peserta.NoPendaftaran = input.NoPendaftaran
	peserta.UUID = uuidRandom.String()
	peserta.Hash = hash
	peserta.Salt = salt
	peserta.Password = randomPassword
	peserta.Role = "peserta|" + jenjang

	newPeserta, err := s.repository.SavePeserta(peserta)
	if err != nil {
		return newPeserta, err
	}
	return newPeserta, nil

}

func (s *service) CreateNewAdmin(input CreateNewAdminInput) (CreateNewAdminResponse, error) {
	uuidRandom := uuid.New()
	//answersJson, err := json.Marshal(input.Answers)
	randomPassword, _ := generateSecureRandomString(8)
	randomNoPendaftaran, _ := generateSecureRandomString(10)
	fmt.Println(randomPassword)
	salt, err := generateSecureRandomString(64)
	concat := randomPassword + salt
	fmt.Println(salt)

	hash, errHash := bcrypt2.Hash(concat)
	if errHash != nil {
		return CreateNewAdminResponse{}, errHash
	}
	fmt.Println(hash)
	if err != nil {
		return CreateNewAdminResponse{}, err
	}
	identityPeserta, _ := json.Marshal(input.Name)

	peserta := PesertaEntity{}
	peserta.Name = identityPeserta
	peserta.UUID = uuidRandom.String()
	peserta.NoPendaftaran = randomNoPendaftaran
	peserta.Hash = hash
	peserta.Salt = salt
	peserta.Password = randomPassword
	peserta.Role = "admin" + "|" + input.Type
	//peserta.Answers = answersJson
	createAdmin := CreateNewAdminResponse{}
	createAdmin.Name = identityPeserta
	_, err = s.repository.SavePeserta(peserta)
	if err != nil {
		return CreateNewAdminResponse{}, err
	}
	return createAdmin, nil

}

func (s *service) GetAll(jenjang string) ([]GetAllPesertaInput, error) {
	peserta, err := s.repository.FindAll(jenjang)
	if err != nil {
		return nil, err
	}
	return peserta, nil
}

func (s *service) UpdatePeserta(uuid string, input UpdatePesertaInput) (UpdatePesertaInput, error) {

	updatePeserta, err := s.repository.UpdatePeserta(uuid, input)
	if err != nil {
		return updatePeserta, err
	}
	return updatePeserta, nil
}

func (s *service) DeletePeserta(uuid string) (bool, error) {
	deletePeserta, err := s.repository.DeletePeserta(uuid)
	if err != nil {
		return deletePeserta, err
	}
	return deletePeserta, nil
}

func generateSecureRandomString(n int) (string, error) {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, n)
	for i := 0; i < n; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterBytes))))
		if err != nil {
			return "", err
		}
		result[i] = letterBytes[randomIndex.Int64()]
	}
	return string(result), nil
}
