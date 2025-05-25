package soal

import "encoding/json"

type Soal struct {
	ID        int             `json:"id"`
	KodeSoal  string          `json:"kode_soal"`
	Waktu     int             `json:"waktu"`
	Questions json.RawMessage `json:"questions"`
}

type TmpSoal struct {
	ID            int             `json:"id"`
	QuestionID    string          `json:"question_id"`
	Pertanyaan    json.RawMessage `json:"pertanyaan"`
	Options       json.RawMessage `json:"options"`
	CorrectAnswer string          `json:"correct_answer"`
	Type          string          `json:"type"`
}
