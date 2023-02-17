package handler

import (
	"crowdfundingweb/helper"
	"crowdfundingweb/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

// mapping input dari user untuk menjadi struct input
type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) RegisterUser(ctx *gin.Context) {
	// tangkap input dari user
	var input user.RegisterUserInput

	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Register failed", http.StatusUnprocessableEntity, "error", errorMessage)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// map input dari user ke struct RegisterUserInput
	// struct di atas passing sebagai parameter service
	newUser, err := h.userService.RegisterUser(input)

	// newUser diformat agar output response body sesuai dengan API spec
	formatter := user.FormatUserResponse(newUser, "token")
	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", formatter)
	if err != nil {
		response := helper.APIResponse("Register failed", http.StatusBadRequest, "error", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(ctx *gin.Context) {
	// user memasukkan email dan password
	// input ditangkap handler
	var input user.LoginInput

	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// mapping dari input user ke input struct
	// input struct pass ke service
	// dalam service mencari email yang dimasukkan
	// jika ketemu maka password dicocokkan
	userLogin, err := h.userService.Login(input)

	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	formatter := user.FormatUserResponse(userLogin, "token")
	response := helper.APIResponse("Login success", http.StatusOK, "success", formatter)

	ctx.JSON(http.StatusOK, response)
}
