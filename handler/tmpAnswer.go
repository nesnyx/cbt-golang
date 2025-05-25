package handler

import (
	"cbt/helper"
	"cbt/peserta"
	"cbt/tmpAnswer"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TmpAnswerHandler struct {
	service tmpAnswer.ServiceInterface
}

func NewTmpAnswerHandler(service tmpAnswer.ServiceInterface) *TmpAnswerHandler {
	return &TmpAnswerHandler{service: service}
}

func (h *TmpAnswerHandler) FindByQuestionID(c *gin.Context) {
	var input tmpAnswer.GetTmpSoalInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		response := helper.APIResponse("Failed Get Tmp Answer By Question ID", 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	getTmpAnswer, err := h.service.FindByQuestionID(input)
	if err != nil {
		response := helper.APIResponse("Failed Get Tmp Answer By Question ID", 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Successfully Get Tmp Answer By Question ID", 200, "success", getTmpAnswer)
	c.JSON(http.StatusOK, response)
}

func (h *TmpAnswerHandler) Create(c *gin.Context) {
	var input tmpAnswer.CreateTmpAnswerInput
	err := c.ShouldBindJSON(&input)
	getCookieUUID := c.MustGet("currentUser").(peserta.GetPesertaTokenInput)
	if err != nil {
		response := helper.APIResponse("Failed Create New Tmp Answer", 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	newTmpAnswer, err := h.service.Save(input, getCookieUUID.UUID)
	if err != nil {
		response := helper.APIResponse("Failed Create New Tmp Answer", 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Successfully Create New Tmp Answer", 201, "success", newTmpAnswer)
	c.JSON(http.StatusOK, response)
}

func (h *TmpAnswerHandler) GetAll(c *gin.Context) {
	getAllTmpAnswer, err := h.service.GetAll()
	if err != nil {
		response := helper.APIResponse("Failed Get All Tmp Answer", 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Successfully Get All Tmp Answer", 200, "success", getAllTmpAnswer)
	c.JSON(http.StatusOK, response)

}

func (h *TmpAnswerHandler) GetAllAnswerByUUID(c *gin.Context) {
	getByUUID := c.MustGet("currentUser").(peserta.GetPesertaTokenInput)
	getAnswers, err := h.service.GetAnswerByUUID(getByUUID.UUID)
	if err != nil {
		response := helper.APIResponse("Failed Get All Tmp Answer", 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Successfully Get All Tmp Answer", 200, "success", getAnswers)
	c.JSON(http.StatusOK, response)
}
