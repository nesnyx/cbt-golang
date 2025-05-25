package tmpAnswer

type CreateTmpAnswerInput struct {
	UUID        string `json:"uuid"`
	Answers     string `json:"answers" binding:"required"`
	QuestionsID string `json:"questions_id" binding:"required"`
}

type QuestionIDInput struct {
	QuestionsID string `json:"questions_id" binding:"required"`
}

type GetTmpSoalInput struct {
	SoalID      int    `json:"soal_id"`
	QuestionsID string `json:"questions_id" binding:"required"`
	PesertaUUID string `json:"peserta_uuid" binding:"required"`
}
