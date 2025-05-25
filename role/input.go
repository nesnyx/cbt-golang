package role

type GetByIDInput struct {
	ID int `json:"id" binding:"required"`
}
