package main

import (
	"crowdfundingweb/auth"
	"crowdfundingweb/handler"
	"crowdfundingweb/helper"
	"crowdfundingweb/user"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// making connection to mysql local database
	dsn := "root:root@tcp(127.0.0.1:3306)/crowdfundingweb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Connection to database success")

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	authService := auth.NewJWTService()
	userHandler := handler.NewUserHandler(userService, authService)

	r := gin.Default()
	api := r.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailabilty)
	api.POST("/avatars", AuthMiddleware(authService, userService), userHandler.UploadAvatar)
	r.Run("127.0.0.1:8080")
}

func AuthMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			resonse := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, resonse)
			return
		}

		var tokenString string
		tokenArray := strings.Split(authHeader, " ")
		if len(tokenArray) == 2 {
			tokenString = tokenArray[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			resonse := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, resonse)
			return
		}

		payload, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			resonse := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, resonse)
			return
		}

		userID := int(payload["user_id"].(float64))

		user, err := userService.GetUserById(userID)
		if err != nil {
			resonse := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, resonse)
			return
		}

		ctx.Set("currentUser", user)
	}
}
