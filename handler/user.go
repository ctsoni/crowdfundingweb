package handler

import (
	"crowdfundingweb/helper"
	"crowdfundingweb/user"
	"fmt"
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
	response := helper.APIResponse("Successfully loggedin", http.StatusOK, "success", formatter)

	ctx.JSON(http.StatusOK, response)
}

func (h *userHandler) CheckEmailAvailabilty(ctx *gin.Context) {
	// ada input email dari user
	var input user.CheckEmailInput
	// input email di mapping ke struct input
	// struct input pass ke service
	// service akan memanggil repository untuk menentukan email sudah ada atau belum
	// repository query ke db
	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Email checking failed", http.StatusUnprocessableEntity, "error", errorMessage)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// mapping dari input user ke input struct
	// input struct pass ke service
	// dalam service mencari email yang dimasukkan
	// jika ketemu maka password dicocokkan
	isEmailAvailable, err := h.userService.IsEmailAvailable(input)

	if err != nil {
		errorMessage := gin.H{"errors": "Server error"}
		response := helper.APIResponse("Email checking failed", http.StatusUnprocessableEntity, "error", errorMessage)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data := gin.H{
		"is_available": isEmailAvailable,
	}

	metaMessage := "Email has been registered"
	if isEmailAvailable {
		metaMessage = "Email is available"
	}

	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)

	ctx.JSON(http.StatusOK, response)

}

func (h *userHandler) UploadAvatar(ctx *gin.Context) {
	// tangkap input dari user (bukan json tapi form body)
	// simpan gambar di folder "images/"
	// di service panggil repository
	// JWT hardcoded (user yang login id == 1)
	// Repo ambil data user yang id == 1
	// Repo update data user simpan lokasi file
	file, err := ctx.FormFile("avatar")
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	// harusnya dapat dari JWT
	userID := 12

	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)

	err = ctx.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{
		"is_uploaded": true,
	}
	response := helper.APIResponse("Upload success", http.StatusOK, "success", data)
	ctx.JSON(http.StatusOK, response)
}
