package services

import (
	"log"
	"os"

	"gopkg.in/gomail.v2"
)

// SendEmail is the Golang equivalent of sending an email (we can't use Javascript here!)
func SendEmail(to string, subject string, htmlBody string) error {
	m := gomail.NewMessage()
	
	// Setup headers
	m.SetHeader("From", os.Getenv("GOOGLE_USER")) // usually your Gmail address
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	// Note: If you want to use OAuth2 in Go, it requires a bit more setup with golang.org/x/oauth2 
	// because Go doesn't have "nodemailer". 
	// The easiest way to start is using an App Password instead of OAuth2.


	
    d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("GOOGLE_USER"), os.Getenv("GOOGLE_CLIENT_SECRET"))
	if err := d.DialAndSend(m); err != nil {
		log.Println("Error sending email:", err)
		return err
	}
	
	log.Println("Email sent successfully to:", to)
	return nil
}
