package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	//"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"projectgo/authentication"
	"projectgo/database"
	"projectgo/email"

	"projectgo/model"
)

var Err string
var Verify model.AdminModel
var UserTable []model.UserModel

func PostAdmin(c *gin.Context) {
	var admin model.AdminModel
	// Binding JSON data to user struct
	if err := c.BindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Validate that all required fields have values
	if admin.Name == "" || admin.Email == "" || admin.Password == "" ||
		admin.Phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields must have values"})
		return
	}
	if !isValidEmail(admin.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// Validate phone number format
	if !isValidPhoneNumber(admin.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number format"})
		return
	}

	fmt.Println(admin)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	admin.Password = string(hashedPassword)

	var existingadmin model.AdminModel
	if err := database.DB.Where("phone = ?", admin.Email).First(&existingadmin).Error; err == nil {
		// admin already exists, return error
		c.JSON(http.StatusConflict, gin.H{"message": "admin already exists"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "database error"})
		return
	}

	// Send OTP to the email
	err1 := email.SendEmailWithOTP(admin.Email)
	if err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send OTP", "data": err1.Error()})
		return
	}

	adminData, err := json.Marshal(&admin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal admin", "data": err.Error()})
		return
	}
	// Store phone number in Redis for OTP verification
	key := fmt.Sprintf("admin:%s", admin.Phone)
	err = database.SetRedis(key, adminData, time.Minute*5)
	if err != nil {
		fmt.Println("Error setting admin in Redis:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Status": false, "Data": nil, "Message": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Otp generated successfully. Proceed to verification page>>>"})
}
func AdminLogin(c *gin.Context) {
	var loginReq struct {
		Phone    string `json:"phone" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.BindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !isValidPhoneNumber(loginReq.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number format"})
		return
	}
	var existinguser model.AdminModel

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
	token, err := authentication.GenerateAdminToken(loginReq.Phone)
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
func isValidEmail(email string) bool {
	// Simple regex pattern for basic email validation
	fmt.Println(" check email validity")
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func isValidPhoneNumber(phone string) bool {
	// Simple regex pattern for basic phone number validation
	fmt.Println(" check pfone validity")
	const phoneRegex = `^\+?[1-9]\d{1,14}$` // E.164 international phone number format
	re := regexp.MustCompile(phoneRegex)
	return re.MatchString(phone)
}
func AdminOtpVerify(c *gin.Context) {

	// Bind OTP verification request data
	var OTPverify model.VerifyOTPAdmin
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
	key := fmt.Sprintf("admin:%s", OTPverify.Phone)
	value, err := database.GetRedis(key)
	if err != nil {
		fmt.Println("Error retrieving OTP from Redis:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "Data": nil, "Message": "Internal server error"})
		return
	}

	// Bind user data from request body
	var userData model.AdminModel

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

func AdminLogout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
