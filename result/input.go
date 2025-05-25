package result

import "encoding/json"

type CreateResultsInput struct {
	UUID       string          `json:"uuid" binding:"required"`
	Data       json.RawMessage `json:"data" binding:"required"`
	Score      float64         `json:"score" binding:"required"`
	Keterangan json.RawMessage `json:"keterangan"`
}

type KeteranganResultsInput struct {
	Keterangan json.RawMessage `json:"keterangan"`
}

type GetResultByIDInput struct {
	ID int `uri:"id"`
}

type UpdateResultsInput struct {
	Lulus bool `json:"lulus" binding:"required"`
}
