package peserta

import "encoding/json"

type CreatePesertaInput struct {
	Name          string          `json:"name" binding:"required"`
	NoPendaftaran string          `json:"no_pendaftaran" binding:"required"`
	AllowedToken  string          `json:"allowed_token"`
	StartTime     int             `json:"start_time"`
	Answers       json.RawMessage `json:"answers"`
}

type UpdatePesertaStartTimeInput struct {
	StartTime int `json:"start_time" binding:"required"`
}

type CreateNewPesertaInput struct {
	Name          json.RawMessage `json:"identity" binding:"required"`
	NoPendaftaran string          `json:"no_pendaftaran" binding:"required"`
	Answers       json.RawMessage `json:"answers"`
}

type CreateNewPesertaInput2 struct {
	Name          json.RawMessage `json:"identity" binding:"required"`
	Password      string          `json:"password" binding:"required"`
	NoPendaftaran string          `json:"no_pendaftaran" binding:"required"`
	Answers       json.RawMessage `json:"answers"`
}

type CreateNewAdminInput struct {
	Name json.RawMessage `json:"identity" binding:"required"`
	Type string          `json:"type" binding:"required"`
}

type CreateNewAdminResponse struct {
	Name json.RawMessage `json:"name"`
}

type UpdatePesertaAnswers struct {
	Answers json.RawMessage `json:"answers" binding:"required"`
}

type GetAllPesertaInput struct {
	UUID          string          `json:"uuid" binding:"required"`
	Name          json.RawMessage `json:"identity" binding:"required"`
	NoPendaftaran string          `json:"no_pendaftaran"`
	Role          string          `json:"role"`
}

type GetPesertaTokenInput struct {
	Name         json.RawMessage `json:"identity" binding:"required"`
	UUID         string          `json:"UUID" binding:"required"`
	Role         string          `json:"role"`
	StartTime    int             `json:"start_time"`
	IsStarted    bool            `json:"is_started"`
	IsFinish     bool            `json:"is_finished"`
	AllowedToken string          `json:"allowed_token"`
}

type GetUUIDPesertaInput struct {
	UUID string `uri:"uuid" binding:"required"`
}

type UpdatePesertaAnswerInput struct {
	Answer     string `json:"answer" binding:"required"`
	UUID       string `json:"uuid" binding:"required"`
	QuestionID string `json:"question_id" binding:"required"`
}

type GetByNameInput struct {
	Name string `json:"name" binding:"required"`
}

type LoginInput struct {
	Name          string `json:"name"`
	NoPendaftaran string `json:"no_pendaftaran" binding:"required"`
	Password      string `json:"password" binding:"required"`
	Token         string
}

type UpdatePesertaInput struct {
	Name          string `json:"name" binding:"required"`
	NoPendaftaran string `json:"no_pendaftaran" binding:"required"`
	Institusi     string `json:"institusi" binding:"required"`
	NoTelp        string `json:"no_telp" binding:"required"`
	Email         string `json:"email" binding:"required"`
}
