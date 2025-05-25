package institution

type CreateNewInstitutionInput struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
}

type DeleteInstitutionInput struct {
	ID int `uri:"id" binding:"required"`
}

type UpdateInstitutionInput struct {
	ID      int    `json:"id" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
}
