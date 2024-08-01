package handler

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"notifications/database"
	"notifications/models"
	"os"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

type InputUser struct {
	Email 		string 		`json:"email" binding:"required"`
	Password 	string 		`json:"password" binding:"required"`
}

type InputEmail struct {
	Email		string 		`json:"email" binding:"required"`	
}

type InputOtpPassword struct {
	NewPassword string 		`json:"newPassword" binding:"required"`
	Otp			string		`json:"otp" binding:"required"`
}

func GenerateOTP(length int) (string, error) {
	charset := "0123456789"
	otp := make([]byte, length)
	_, err := rand.Read(otp)

	if err != nil {
		return "", err
	}

	for i := range otp {
		otp[i] = charset[int(otp[i])%len(charset)]
	}

	return string(otp), nil
}

func CreateUser(c *gin.Context) {
	var input InputUser
	var ExistingUser models.Gouser

	err := c.ShouldBindJSON(&input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"status": "fail",
			"message": "failed to load json",
		})
		return
	}

	result := database.DB.Where("email = ?", input.Email).Find(&ExistingUser).Error

	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"status": "fail",
			"message": "failed to create user, silahkan gunakan email yang lain",
		})
		return
	}

	newUser := models.Gouser {
		Email: input.Email,
		Password: input.Password,
	}

	database.DB.Create(&newUser)

	c.JSON(http.StatusCreated, gin.H {
		"status": "success",
		"message": "user created",
		"data": input.Email,
	})
}

func InputEmailChangePassword(c *gin.Context) {
	err := godotenv.Load()

	if err != nil {
		panic("Failed to load .env file")
	}

	var input InputEmail
	var ExistingUser models.Gouser

	err = c.ShouldBindJSON(&input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"status": "fail",
			"message": "failed to load json",
		})
		return
	}

	result := database.DB.Where("email = ?", input.Email).Find(&ExistingUser).Error

	if result != nil {
		c.JSON(http.StatusNotFound, gin.H {
			"status": "fail",
			"message": "email tidak ditemukan",
		})
		return
	}

	otp, err := GenerateOTP(8)

	if err != nil {
		log.Fatal(err)
	}

	mail := gomail.NewMessage()
	
	mail.SetHeader("From", os.Getenv("EMAIL"))
	mail.SetHeader("To", input.Email)
	mail.SetHeader("Subject", "Your OTP Verification")

	mail.SetBody("text/plain", fmt.Sprintf("This is your OTP Code %s, please don't share it with anyone", otp))

	setup := gomail.NewDialer(os.Getenv("EMAIL_SMTP"), 587, os.Getenv("EMAIL"), os.Getenv("EMAIL_PASS"))
	err = setup.DialAndSend(mail)

	if err != nil {
		log.Fatal(err)
	}

	newOtp := models.Gootp {
		Email: input.Email,
		Otp: otp,
	}

	database.DB.Create(&newOtp)

	c.JSON(http.StatusOK, gin.H {
		"status": "success",
		"message": "OTP sent to your email",
	})
}

func InputOtpChangePassword(c *gin.Context) {
	var input InputOtpPassword
	var ExistingOtp models.Gootp
	var ChangePassword models.Gouser

	err := c.ShouldBindJSON(&input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"status": "fail",
			"message": "failed to load json",
		})
		return
	}

	result := database.DB.Where("otp = ?", input.Otp).First(&ExistingOtp).Error

	if result != nil {
		c.JSON(http.StatusNotFound, gin.H {
			"status": "fail",
			"message": "OTP salah, coba lagi",
		})
		return
	}

	database.DB.Where("email = ?", ExistingOtp.Email).First(&ChangePassword)
	ChangePassword.Password = input.NewPassword
	database.DB.Save(&ChangePassword)
	database.DB.Delete(&ExistingOtp)

	c.JSON(http.StatusOK, gin.H {
		"status": "success",
		"message": "password changed",
	})
}