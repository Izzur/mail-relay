package main

import (
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/health", health)
	r.POST("/api.sendgrid.com/v3/mail/send", sendgrid)
	r.POST("/api.sendinblue.com/v3/smtp/email", sendinblue)
	r.Run()
}

func health(c *gin.Context) {
	c.JSON(200, gin.H{"status": "OK"})
}

func sendgrid(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{})
}

func sendinblue(c *gin.Context) {
	var body Sendinblue
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	to := mapPersonToEmail(body.To)
	cc := []string{}
	sendMail(to, cc, body.Subject, body.HTMLContent)
	c.JSON(http.StatusAccepted, gin.H{})
}

func sendMail(to []string, cc []string, subject, message string) error {
	user := os.Getenv("MAIL_RELAY_SMTP_USER")
	pass := os.Getenv("MAIL_RELAY_SMTP_PASS")
	host := os.Getenv("MAIL_RELAY_SMTP_HOST")
	port := os.Getenv("MAIL_RELAY_SMTP_PORT")
	fromName := os.Getenv("MAIL_RELAY_FROM_NAME")
	fromEmail := os.Getenv("MAIL_RELAY_FROM_EMAIL")

	body := "From: " + fromName + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Cc: " + strings.Join(cc, ",") + "\n" +
		"Subject: " + subject + "\n\n" +
		message

	auth := smtp.PlainAuth("", user, pass, host)
	smtpAddr := fmt.Sprintf("%s:%s", host, port)

	err := smtp.SendMail(smtpAddr, auth, fromEmail, append(to, cc...), []byte(body))
	if err != nil {
		return err
	}

	return nil
}

func mapPersonToEmail(person []Person) []string {
	res := make([]string, len(person))
	for i, val := range person {
		res[i] = val.Email
	}
	return res
}

// Person name and email
type Person struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Sendinblue request body
type Sendinblue struct {
	To          []Person `json:"to"`
	Subject     string   `json:"subject"`
	HTMLContent string   `json:"htmlContent"`
}
