package result

import "encoding/json"

type Results struct {
	ID         int             `json:"id"`
	UUID       string          `json:"uuid"`
	Data       json.RawMessage `json:"data"`
	Score      float64         `json:"score"`
	Lulus      bool            `json:"lulus"`
	Keterangan json.RawMessage `json:"keterangan"`
	Name       string          `json:"name"`
	NoTelp     string          `json:"no_telp"`
	Email      string          `json:"email"`
}

type ResultsEntity struct {
	ID         int             `json:"id"`
	UUID       string          `json:"uuid"`
	Data       json.RawMessage `json:"data"`
	Score      float64         `json:"score"`
	Lulus      bool            `json:"lulus"`
	Keterangan json.RawMessage `json:"keterangan"`
}

type ResultsJoinPeserta struct {
	Name  string `json:"name"`
	Telp  string `json:"no_telp"`
	Email string `json:"Email"`
}
