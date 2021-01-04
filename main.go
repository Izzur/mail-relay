package main

import (
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

var (
	user string
	pass string
	host string
	port string
	from string
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	userName, exists := os.LookupEnv("USERNAMEEMAIL")
	if exists {
		user = userName
		log.Println(user)
	}
	passUser, exists := os.LookupEnv("PASSWORD")
	if exists {
		pass = passUser
		log.Println(pass)
	}
	hostUser, exists := os.LookupEnv("HOST")
	if exists {
		host = hostUser
		log.Println(host)
	}
	portUser, exists := os.LookupEnv("PORTEMAIL")
	if exists {
		port = portUser
		log.Println(port)
	}
	formUser, exists := os.LookupEnv("FROM")
	if exists {
		from = formUser
		log.Println(from)
	}

}

func main() {
	r := gin.Default()
	r.GET("/", hello)
	r.GET("/health", health)
	r.POST("/api.sendgrid.com/v3/mail/send", sendgrid)
	r.POST("/api.sendinblue.com/v3/smtp/email", sendinblue)
	r.Run()
}

func hello(c *gin.Context) {
	c.JSON(200, gin.H{"status": "Hallo"})
}

func health(c *gin.Context) {
	c.JSON(200, gin.H{"status": "OK"})
}

func sendgrid(c *gin.Context) {
	var body Sendgrid
	if c.Request.Header["Authorization"] == nil {
		c.JSON(http.StatusBadRequest, gin.H{})
	}
	token := strings.Split(c.Request.Header["Authorization"][0], " ")[1]
	log.Println(token)
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	to := mapPersonToEmail(body.Personalizations[0].To)
	cc := []string{}
	sendMailv2(to, cc, body.Personalizations[0].Subject, body.Content[0].Value)
	c.JSON(http.StatusAccepted, gin.H{})
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
func sendMailv2(to []string, cc []string, subject, message string) {
	auth := smtp.PlainAuth("", user, pass, host)

	to = []string{to[0]}
	msg := []byte("From:" + from + "\r\n" +
		"To: " + to[0] + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"utf-8\"\r\n\r\n" +
		"\r\n" + message + "\r\n")
	err := smtp.SendMail(host+":"+port, auth, from, to, msg)
	if err != nil {
		log.Fatal(err)
	}
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

// Sendgrid request body
type Sendgrid struct {
	Personalizations []Personalizations `json:"personalizations"`
	From             Person             `json:"from"`
	Subject          string             `json:"subject"`
	Content          []SendgridContent  `json:"content"`
}

// Personalizations for Sendgrid
type Personalizations struct {
	To      []Person `json:"to"`
	Subject string   `json:subject`
}

// SendgridContent for email body
type SendgridContent struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
