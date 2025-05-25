package peserta

import "encoding/json"

type Peserta struct {
	UUID          string
	Name          string
	NoPendaftaran string
	Password      string
	Hash          string
	Salt          string
	AllowedToken  string
	StartTime     int
	LastLogin     int
	Answers       json.RawMessage `gorm:"type:json"`
}

type PesertaEntity struct {
	UUID          string
	Name          json.RawMessage `gorm:"type:json"`
	NoPendaftaran string
	Password      string
	Hash          string
	Salt          string
	AllowedToken  string
	StartTime     int
	Answers       json.RawMessage `gorm:"type:json"`
	Role          string
}
