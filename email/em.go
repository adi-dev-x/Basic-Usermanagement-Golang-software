package email

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"net/smtp"
	"os"
)

// GenerateOTP generates a random OTP of given length
func GenerateOTP(length int) int {
	rand.Seed(time.Now().UnixNano())

	// Generate a random number between 10000 and 99999
	randomNum := rand.Intn(90000) + 10000

	fmt.Println("Random 5-digit number:", randomNum)
	return randomNum
}

// SendEmailWithOTP sends an email with the OTP
func SendEmailWithOTP(email string) error {
	// Generate OTP
	otp := strconv.Itoa(GenerateOTP(6))

	// Construct email message
	message := fmt.Sprintf("Subject: OTP for Verification\n\nYour OTP is: %s", otp)

	SMTPemail := os.Getenv("EMAIL")
	log.Println("my email is-------", SMTPemail)
	SMTPpass := os.Getenv("PASSWORD")

	// Authenticate with SMTP server
	auth := smtp.PlainAuth("", SMTPemail, SMTPpass, "smtp.gmail.com")

	// Send email using SMTP server
	err := smtp.SendMail("smtp.gmail.com:587", auth, SMTPemail, []string{email}, []byte(message))
	if err != nil {
		log.Println("Error sending email:", err)
		return err
	}

	return nil
}
