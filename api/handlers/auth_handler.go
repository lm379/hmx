package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lm379/hmx/models"
	"github.com/lm379/hmx/services"
	"github.com/lm379/hmx/utils"
)

// HandleSendCode (POST /api/v1/auth/send-code)
func HandleSendCode(c *gin.Context) {
	var input models.SendCodeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	response, status := services.SendVerificationCode(input.Email)
	if status != http.StatusOK {
		utils.Error(c, status, response["error"].(string))
		return
	}
	utils.Success(c, response)
}

// HandleRegister (POST /api/v1/auth/register)
func HandleRegister(c *gin.Context) {
	var input models.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	response, status := services.RegisterUser(input)
	if status != http.StatusCreated {
		utils.Error(c, status, response["error"].(string))
		return
	}
	utils.Created(c, response)
}

// HandleLogin (POST /api/v1/auth/login)
func HandleLogin(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	response, status := services.LoginUser(input)
	if status != http.StatusOK {
		utils.Error(c, status, response["error"].(string))
		return
	}
	utils.Success(c, response)
}

// HandleForgetPassword (POST /api/v1/auth/forget)
func HandleForgetPassword(c *gin.Context) {
	var input models.ForgetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	response, status := services.ForgetPassword(input)
	if status != http.StatusOK {
		utils.Error(c, status, response["error"].(string))
		return
	}
	utils.Success(c, response)
}
