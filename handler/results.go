package handler

import (
	"cbt/helper"
	"cbt/peserta"
	"cbt/result"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ResultsHandler struct {
	service result.ServiceInterface
}

func NewResultsHandler(service result.ServiceInterface) *ResultsHandler {
	return &ResultsHandler{service}
}

func (h *ResultsHandler) SaveResults(c *gin.Context) {
	var input result.CreateResultsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response := helper.APIResponse("Save Result Failed", 400, "error", gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}
	newResult, err := h.service.CreateNewResult(input)
	if err != nil {
		response := helper.APIResponse("Save Result Failed", 400, "error", gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Save Result Successfully", 201, "success", newResult)
	c.JSON(http.StatusBadRequest, response)

}

func (h *ResultsHandler) GetAllResults(c *gin.Context) {
	getUserRole, _ := c.MustGet("currentUser").(peserta.GetPesertaTokenInput)
	splitRole := strings.Split(getUserRole.Role, "|")
	getAll, err := h.service.GetAllResults(splitRole[1])
	if err != nil {
		response := helper.APIResponse("Get All Results Failed", 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Get All Results Successfully", 200, "success", getAll)
	c.JSON(http.StatusOK, response)
}

func (h *ResultsHandler) DeleteResult(c *gin.Context) {
	var input result.GetResultByIDInput
	err := c.ShouldBindUri(&input)
	if err != nil {
		errorMessage := gin.H{"errors": err}
		response := helper.APIResponse("Failed Delete Results", 400, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	deleteResult, err := h.service.DeleteResults(input.ID)
	if err != nil {
		errorMessage := gin.H{"errors": err}
		response := helper.APIResponse("Failed Delete Results", 400, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Success Delete Results", 200, "success", deleteResult)
	c.JSON(http.StatusBadRequest, response)
}

func (h *ResultsHandler) UpdateResult(c *gin.Context) {
	var input result.UpdateResultsInput
	var inputId result.GetResultByIDInput
	errUri := c.ShouldBindUri(&inputId)
	err := c.ShouldBindJSON(&input)

	if err != nil {
		errorMessage := gin.H{"errors": err}
		response := helper.APIResponse("Failed Update Results", 400, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if errUri != nil {
		errorMessage := gin.H{"errors": errUri}
		response := helper.APIResponse("Failed Update Results", 400, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	updateResult, err := h.service.UpdateResutls(inputId.ID, input.Lulus)
	if err != nil {
		errorMessage := gin.H{"errors": err}
		response := helper.APIResponse("Failed Update Results", 400, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Failed Update Results", 400, "error", updateResult)
	c.JSON(http.StatusBadRequest, response)

}
