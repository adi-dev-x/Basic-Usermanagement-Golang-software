package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"projectgo/authentication"
	"projectgo/database"
	"projectgo/email"
	"projectgo/model"
	"regexp"

	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func UserLogin(c *gin.Context) {
	var loginReq struct {
		Phone    string `json:"phone" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.BindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !userisValidPhoneNumber(loginReq.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number format"})
		return
	}
	var existinguser model.UserModel
	fmt.Println("this is the phone ", loginReq.Phone)
	if err := database.DB.Where("phone = ?", loginReq.Phone).First(&existinguser).Error; err != nil {

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid phone number or phone number is not present"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existinguser.Password), []byte(loginReq.Password)); err != nil {
		// Incorrect password
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid phone number or password"})
		return
	}

	// //Generate JWT token for the patient
	token, err := authentication.GenerateUserToken(loginReq.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return the token
	c.JSON(http.StatusOK, gin.H{
		"Status":  "Success",
		"message": "Login sucessful",
		"token":   token,
	})
}

func UserSignup(c *gin.Context) {
	var user model.UserModel
	// Binding JSON data to user struct
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !isValidEmail(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// Validate phone number format
	if !userisValidPhoneNumber(user.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number format"})
		return
	}
	fmt.Println(user)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	var existinguser model.UserModel
	if err := database.DB.Where("phone = ?", user.Email).First(&existinguser).Error; err == nil {
		// user already exists, return error
		c.JSON(http.StatusConflict, gin.H{"message": "user already exists"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "database error"})
		return
	}

	// Send OTP to the email
	err1 := email.SendEmailWithOTP(user.Email)
	if err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send OTP", "data": err1.Error()})
		return
	}

	userData, err := json.Marshal(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal user", "data": err.Error()})
		return
	}
	// Store phone number in Redis for OTP verification
	key := fmt.Sprintf("user:%s", user.Phone)
	err = database.SetRedis(key, userData, time.Minute*5)
	if err != nil {
		fmt.Println("Error setting user in Redis:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Status": false, "Data": nil, "Message": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Otp generated successfully. Proceed to verification page>>>"})

}
func userisValidPhoneNumber(phone string) bool {
	// Simple regex pattern for basic phone number validation
	fmt.Println(" check pfone validity")
	const phoneRegex = `^\+?[1-9]\d{1,14}$` // E.164 international phone number format
	re := regexp.MustCompile(phoneRegex)
	return re.MatchString(phone)
}
func UserOtpVerify(c *gin.Context) {

	// Bind OTP verification request data
	var OTPverify model.VerifyOTP
	if err := c.BindJSON(&OTPverify); err != nil {
		// fmt.Println("i'm here")
		fmt.Println("Error parsing JSON:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"Status": false, "Data": nil, "Message": "Failed to parse JSON data"})
		return
	}

	// Check if OTP is empty
	if OTPverify.Otp == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Status": false, "Message": "OTP is required"})
	}

	// Retrieve patient data from Redis
	key := fmt.Sprintf("user:%s", OTPverify.Phone)
	value, err := database.GetRedis(key)
	if err != nil {
		fmt.Println("Error retrieving OTP from Redis:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "Data": nil, "Message": "Internal server error"})
		return
	}

	// Bind user data from request body
	var userData model.UserModel

	err = json.Unmarshal([]byte(value), &userData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unmarshal patient", "data": err.Error()})
		return
	}

	err = database.DB.Create(&userData).Error
	if err != nil {
		fmt.Println("Error creating Patient:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Status": false, "Data": nil, "Message": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Status": true, "Message": "OTP verified successfully and user has been created. Login to continue..."})
}
