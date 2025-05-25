package handler

import (
	"cbt/helper"
	"cbt/peserta"
	"cbt/soal"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type SoalHandler struct {
	service soal.ServiceInterface
}

func NewSoalHandler(service soal.ServiceInterface) *SoalHandler {
	return &SoalHandler{service}
}

func (h *SoalHandler) GetAllTmpSoalByType(c *gin.Context) {
	getUserRole, _ := c.MustGet("currentUser").(peserta.GetPesertaTokenInput)
	splitRole := strings.Split(getUserRole.Role, "|")
	if splitRole[1] == "SMA" {
		getTmpSoal, err := h.service.GetAllTmpSoalByType(splitRole[1])
		if err != nil {
			response := helper.APIResponse("Failed Get All Tmp Soal By Type", 400, "error", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}
		response := helper.APIResponse("Success Get All Tmp Soal By Type", 200, "success", getTmpSoal)
		c.JSON(http.StatusOK, response)
		return
	} else if splitRole[1] == "SMP" {
		getTmpSoal, err := h.service.GetAllTmpSoalByType(splitRole[1])
		if err != nil {
			response := helper.APIResponse("Failed Get All Tmp Soal By Type", 400, "error", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}
		response := helper.APIResponse("Success Get All Tmp Soal By Type", 200, "success", getTmpSoal)
		c.JSON(http.StatusOK, response)
		return
	}
}

func (h *SoalHandler) GetAllTmpSoalByTypeStudent(c *gin.Context) {
	getUserRole, _ := c.MustGet("currentUser").(peserta.GetPesertaTokenInput)
	var getUserRoleStruct map[string]interface{}
	err := json.Unmarshal([]byte(getUserRole.Name), &getUserRoleStruct)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	// Mengambil nilai dari key "Name"
	getJenjang := getUserRoleStruct["jenjang"]
	fmt.Println(getJenjang)
	if getJenjang == "SMA" {
		getTmpSoal, err := h.service.GetAllTmpSoalByType("SMA")
		if err != nil {
			response := helper.APIResponse("Failed Get All Tmp Soal By Type", 400, "error", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}
		response := helper.APIResponse("Success Get All Tmp Soal By Type", 200, "success", getTmpSoal)
		c.JSON(http.StatusOK, response)
		return
	} else if getJenjang == "SMP" {
		getTmpSoal, err := h.service.GetAllTmpSoalByType("SMP")
		if err != nil {
			response := helper.APIResponse("Failed Get All Tmp Soal By Type", 400, "error", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}
		response := helper.APIResponse("Success Get All Tmp Soal By Type", 200, "success", getTmpSoal)
		c.JSON(http.StatusOK, response)
		return
	}
}

func (h *SoalHandler) CreateSoal(c *gin.Context) {
	var input soal.CreateSoalInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errorMessage := gin.H{"errors": err}
		response := helper.APIResponse("Failed Create New Soal", 400, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	newSoal, err := h.service.CreateSoal(input)
	if err != nil {
		response := helper.APIResponse("Failed Create New Soal", 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Successfully Create New Soal", 201, "success", newSoal)
	c.JSON(http.StatusOK, response)
}

func (h *SoalHandler) GetSoal(c *gin.Context) {
	getAll, err := h.service.GetAllSoal()
	if err != nil {
		response := helper.APIResponse("Failed Get All Soal", 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Successfully Get All Soal", 200, "success", getAll)
	c.JSON(http.StatusOK, response)
}

func (h *SoalHandler) GetSoalByKodeSoal(c *gin.Context) {
	var input soal.GetSoalByKodeSoalInput
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed Get Soal By Kode Soal", 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	getSoal, err := h.service.GetByKodeSoal(input)
	if err != nil {
		response := helper.APIResponse("Failed Get Soal By Kode Soal", 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	c.JSON(http.StatusOK, getSoal)

}

func (h *SoalHandler) DeleteSoal(c *gin.Context) {
	var input soal.GetSoalByIDInput
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed Delete Soal By ID", 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	deleteSoal, err := h.service.DeleteSoal(input)
	if err != nil {
		response := helper.APIResponse("Failed Delete Soal By ID", 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Successfully Delete Soal By ID", 200, "success", deleteSoal)
	c.JSON(http.StatusOK, response)
}

func (h *SoalHandler) CreateTmpSoal(c *gin.Context) {
	var input soal.CreateTmpSoalInputan
	err := c.ShouldBindJSON(&input)
	if err != nil {
		response := helper.APIResponse("Failed Create Tmp Soal", 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	newSoal, err := h.service.CreateTmpSoal(input)
	if err != nil {
		response := helper.APIResponse("Failed Create Tmp Soal", 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Success Create Tmp Soal", 201, "success", newSoal)
	c.JSON(http.StatusOK, response)
}

func (h *SoalHandler) GetAllTmpSoal(c *gin.Context) {
	getTmpSoal, err := h.service.GetAllTmpSoal()
	if err != nil {
		response := helper.APIResponse("Failed Get All Tmp Soal", 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Successfully Get All Tmp Soal", 200, "success", getTmpSoal)
	c.JSON(http.StatusOK, response)

}

func (h *SoalHandler) UpdateTmpSoal(c *gin.Context) {
	var input soal.UpdateTmpSoalOptionsInputan
	err := c.ShouldBindJSON(&input)
	if err != nil {
		message := err.Error()
		response := helper.APIResponse(message, 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	updateTmpSoal, err := h.service.UpdateTmpSoal(input)

	if err != nil {
		message := err.Error()
		response := helper.APIResponse(message, 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Update Tmp Soal", 200, "success", updateTmpSoal)
	c.JSON(http.StatusOK, response)

}
