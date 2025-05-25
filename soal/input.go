package soal

import "encoding/json"

type CreateSoalInput struct {
	KodeSoal  string          `json:"kode_soal" binding:"required"`
	Waktu     int             `json:"waktu" binding:"required"`
	Questions json.RawMessage `json:"questions" binding:"required"`
}

type CreateTmpSoalInputan struct {
	Pertanyaan    json.RawMessage `json:"pertanyaan" binding:"required"`
	Options       json.RawMessage `json:"options" binding:"required"`
	CorrectAnswer string          `json:"correct_answer"`
	TimeDuration  int             `json:"time_duration"`
	Type          string          `json:"type" binding:"required"`
}

type GetSoalByKodeSoalInput struct {
	KodeSoal string `uri:"kode_soal" binding:"required"`
}

type GetSoalByIDInput struct {
	ID int `uri:"id" binding:"required"`
}

type UpdateTmpSoalInputan struct {
	ID              int             `json:"id" binding:"required" `
	Updates         json.RawMessage `json:"updates" binding:"required"`
	QuestionText    string          `json:"question_text" binding:"required"`
	QuestionPicture string          `json:"question_picture" binding:"required"`
	CorrectAnswer   string          `json:"correct_answer" binding:"required"`
	TimeDuration    int             `json:"time_duration"`
}

type UpdateTmpSoalOptionsInputan struct {
	ID              int             `json:"id" binding:"required" `
	Updates         json.RawMessage `json:"updates" binding:"required"`
	QuestionText    string          `json:"question_text" binding:"required"`
	QuestionPicture string          `json:"question_picture" binding:"required"`
	CorrectAnswer   string          `json:"correct_answer" binding:"required"`
	TimeDuration    int             `json:"time_duration"`
}
