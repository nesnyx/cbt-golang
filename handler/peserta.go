package handler

import (
	"cbt/auth"
	"cbt/helper"
	"cbt/peserta"
	"cbt/result"
	"cbt/utils"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type pesertaHandler struct {
	service     peserta.Service
	authService auth.Service
}

func NewPesertaHandler(service peserta.Service, authService auth.Service) *pesertaHandler {
	return &pesertaHandler{service, authService}
}

func (h *pesertaHandler) Login(c *gin.Context) {
	var input peserta.LoginInput
	err := c.ShouldBindJSON(&input)
	expirationTime := time.Now().Add(86400 * time.Second)
	if err != nil {
		response := helper.APIResponse("Failed Login", 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	login, err := h.service.Login(input, expirationTime.Unix())
	if err != nil {
		response := helper.APIResponse("Failed Login", 400, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := h.authService.GenerateToken(login.UUID)
	if err != nil {
		response := helper.APIResponse("Failed Login", 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	encodedCookies, _ := utils.SetEncryptCookies(token)

	cookieValue := fmt.Sprintf("%d|%s|%d", expirationTime.Unix(), encodedCookies, 86400)
	c.SetCookie("session_token", cookieValue, 86400, "/", "cbt.edunex.id", true, true)
	c.Header("Session-UUID", login.UUID)
	c.JSON(http.StatusOK, gin.H{"token": token, "encodedCookies": encodedCookies, "cookieValue": cookieValue})
}

func (h *pesertaHandler) Logout(c *gin.Context) {
	_, err := h.service.Logout(c)
	if err != nil {
		response := helper.APIResponse("Failed Logout", 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	c.SetCookie("session_token", "", -1, "/", "cbt.edunex.id", true, true) // Menghapus cookie

	c.JSON(http.StatusOK, gin.H{"message": "Logout berhasil!"})
}

func (h *pesertaHandler) StartQuiz(c *gin.Context) {
	getCookieUUID := c.MustGet("currentUser").(peserta.GetPesertaTokenInput)
	updateStartTime, err := h.service.UpdateStartTime(getCookieUUID.UUID, true, false)
	if err != nil {
		response := helper.APIResponse("Failed UpdateStartTime", 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Success", 200, "success", updateStartTime)
	c.JSON(http.StatusOK, response)

}

func (h *pesertaHandler) FinishQuiz(c *gin.Context) {
	getCookieUUID := c.MustGet("currentUser").(peserta.GetPesertaTokenInput)
	var input result.KeteranganResultsInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		response := helper.APIResponse("Failed UpdateStartTime", 400, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}
	updateStartTime, err := h.service.UpdateStartTimeFinish(getCookieUUID.UUID, false, true, input.Keterangan)
	if err != nil {
		response := helper.APIResponse("Failed UpdateStartTime", 400, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Success", 200, "success", updateStartTime)
	c.JSON(http.StatusOK, response)

}

func (h *pesertaHandler) GetAll(c *gin.Context) {
	getUserRole, _ := c.MustGet("currentUser").(peserta.GetPesertaTokenInput)
	splitRole := strings.Split(getUserRole.Role, "|")
	getAllPeserta, err := h.service.GetAll(splitRole[1])
	if err != nil {
		response := helper.APIResponse("Failed Get All Peserta", 400, strconv.Itoa(http.StatusBadRequest), nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Successfully Get All Peserta", 200, strconv.Itoa(http.StatusCreated), getAllPeserta)
	c.JSON(http.StatusOK, response)
}

func (h *pesertaHandler) GetByUUID(c *gin.Context) {
	var inputUUID peserta.GetUUIDPesertaInput
	err := c.ShouldBindUri(&inputUUID)
	if err != nil {
		response := helper.APIResponse("Failed Get Peserta By UUID", 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	getPeserta, err := h.service.GetByUUID(inputUUID.UUID)
	if err != nil {
		response := helper.APIResponse("Failed Get Peserta By UUID", 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Successfully Get Peserta By UUID", 200, "success", getPeserta)
	c.JSON(http.StatusOK, response)
}

func (h *pesertaHandler) SavePeserta(c *gin.Context) {
	var input peserta.CreatePesertaInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errorMessage := gin.H{"errors": err}
		response := helper.APIResponse("Failed Create New Peserta", 400, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	createPeserta, err := h.service.Create(input)
	if err != nil {
		response := helper.APIResponse("Failed Create New Peserta", 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Successfully Create New Peserta", 201, "success", createPeserta)
	c.JSON(http.StatusCreated, response)
}

func (h *pesertaHandler) SaveNewPeserta(c *gin.Context) {
	var input peserta.CreateNewPesertaInput
	getUserRole, _ := c.MustGet("currentUser").(peserta.GetPesertaTokenInput)
	splitRole := strings.Split(getUserRole.Role, "|")
	fmt.Println(splitRole[1])
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errorMessage := gin.H{"errors": err}
		response := helper.APIResponse("Failed Create New Peserta", 400, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	createPeserta, err := h.service.CreateNewPeserta(input, splitRole[1])
	if err != nil {
		response := helper.APIResponse("Failed Create New Peserta", 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Successfully Create New Peserta", 201, "success", createPeserta)
	c.JSON(http.StatusCreated, response)
}

func (h *pesertaHandler) SaveNewPeserta2(c *gin.Context) {
	var input peserta.CreateNewPesertaInput2
	getUserRole, _ := c.MustGet("currentUser").(peserta.GetPesertaTokenInput)
	splitRole := strings.Split(getUserRole.Role, "|")
	fmt.Println(splitRole[1])
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errorMessage := gin.H{"errors": err}
		response := helper.APIResponse("Failed Create New Peserta", 400, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	createPeserta, err := h.service.CreateNewPeserta2(input, splitRole[1])
	if err != nil {
		response := helper.APIResponse("Failed Create New Peserta", 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Successfully Create New Peserta", 201, "success", createPeserta)
	c.JSON(http.StatusCreated, response)
}

func (h *pesertaHandler) SaveNewAdmin(c *gin.Context) {
	var input peserta.CreateNewAdminInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errorMessage := gin.H{"errors": err}
		response := helper.APIResponse("Failed Create New Peserta", 400, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	createAdmin, err := h.service.CreateNewAdmin(input)
	if err != nil {
		response := helper.APIResponse("Failed Create New Peserta", 400, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Successfully Create New Peserta", 201, "success", createAdmin)
	c.JSON(http.StatusCreated, response)
}

func (h *pesertaHandler) UpdatePeserta(c *gin.Context) {
	var input peserta.UpdatePesertaInput
	var uuid peserta.GetUUIDPesertaInput
	err := c.ShouldBindJSON(&input)
	errUri := c.ShouldBindUri(&uuid)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Failed Update Peserta", 400, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	if errUri != nil {
		errorMessage := gin.H{"errors": errUri.Error()}
		response := helper.APIResponse("Failed Update Peserta", 400, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	updatePeserta, err := h.service.UpdatePeserta(uuid.UUID, input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Failed Update Peserta", 400, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Successfully Update Peserta", 200, "success", updatePeserta)
	c.JSON(http.StatusOK, response)
}

func (h *pesertaHandler) DeletePeserta(c *gin.Context) {
	var input peserta.GetUUIDPesertaInput
	err := c.ShouldBindUri(&input)
	if err != nil {
		errorMessage := gin.H{"errors": err}
		response := helper.APIResponse("Failed Delete Peserta", 400, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	deletePeserta, err := h.service.DeletePeserta(input.UUID)
	if err != nil {
		errorMessage := gin.H{"errors": err}
		response := helper.APIResponse("Failed Delete Peserta", 400, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Successfully Delete Peserta", 200, "success", deletePeserta)
	c.JSON(http.StatusOK, response)

}
