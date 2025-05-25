package handler

import (
	"cbt/helper"
	"cbt/institution"
	"net/http"

	"github.com/gin-gonic/gin"
)

type institutionHandler struct {
	service institution.ServiceInterface
}

func NewInstitutionHandler(service institution.ServiceInterface) *institutionHandler {
	return &institutionHandler{
		service,
	}
}

func (h *institutionHandler) Create(c *gin.Context) {
	var input institution.CreateNewInstitutionInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		message := err.Error()
		response := helper.APIResponse(message, 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	newInstitution, err := h.service.CreateNewInstitution(input)
	if err != nil {
		message := err.Error()
		response := helper.APIResponse(message, 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	c.JSON(http.StatusOK, newInstitution)
}

func (h *institutionHandler) Update(c *gin.Context) {
	var input institution.UpdateInstitutionInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		response := helper.APIResponse("Update Institution Failed", 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	updateInstitution, err := h.service.UpdateInstitution(input)
	if err != nil {
		response := helper.APIResponse(err.Error(), 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	c.JSON(http.StatusOK, updateInstitution)
}

func (h *institutionHandler) Delete(c *gin.Context) {
	var input institution.DeleteInstitutionInput
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Delete Institution Failed", 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
}
