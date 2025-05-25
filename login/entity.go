package login

import "encoding/json"

type Login struct {
	UUID          string
	Name          string
	NoPendaftaran string
	AllowedToken  string
	StartTime     int
	Answers       json.RawMessage `gorm:"type:json"`
}
