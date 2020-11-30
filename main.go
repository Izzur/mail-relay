package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

var (
	user = os.Getenv("MAIL_RELAY_SMTP_USER")
	pass = os.Getenv("MAIL_RELAY_SMTP_PASS")
	host = os.Getenv("MAIL_RELAY_SMTP_HOST")
	port = os.Getenv("MAIL_RELAY_SMTP_PORT")
	from = os.Getenv("MAIL_RELAY_FROM_EMAIL")
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
	// sendMail()
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
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to[0])
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", message)

	port, _ := strconv.Atoi(port)
	d := gomail.NewDialer(host, port, user, pass)
	err := d.DialAndSend(m)
	return err
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
