package main

import (
	"notifications/database"
	"notifications/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	database.DatabaseConnect()

	router.POST("/createuser", handler.CreateUser)
	router.POST("/sendemail", handler.InputEmailChangePassword)
	router.POST("/sendotp", handler.InputOtpChangePassword)

	router.Run()
}