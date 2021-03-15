package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

var (
	user      string
	pass      string
	host      string
	port      string
	from      string
	relayPort string
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err == nil {
		userName, exists := os.LookupEnv("RELAY_SMTP_USER")
		if !exists {
			log.Fatalln("Environment variable should contain RELAY_SMTP_USER")
		}
		user = userName

		passUser, exists := os.LookupEnv("RELAY_SMTP_PASS")
		if !exists {
			log.Fatalln("Environment variable should contain RELAY_SMTP_PASS")
		}
		pass = passUser

		hostUser, exists := os.LookupEnv("RELAY_SMTP_HOST")
		if !exists {
			log.Fatalln("Environment variable should contain RELAY_SMTP_HOST")
		}
		host = hostUser

		portUser, exists := os.LookupEnv("RELAY_SMTP_PORT")
		if !exists {
			log.Fatalln("Environment variable should contain RELAY_SMTP_PORT")
		}
		port = portUser

		fromUser, exists := os.LookupEnv("RELAY_FROM_EMAIL")
		if !exists {
			log.Fatalln("Environment variable should contain RELAY_FROM_EMAIL")
		}
		from = fromUser

		portService, exists := os.LookupEnv("RELAY_SERVICE_PORT")
		if !exists {
			log.Fatalln("Environment variable should contain RELAY_SERVICE_PORT")
		}
		relayPort = portService
	}
}

func main() {
	r := gin.Default()
	r.GET("/", hello)
	r.GET("/health", health)
	r.POST("/api.sendgrid.com/v3/mail/send", sendgrid)
	r.POST("/api.sendinblue.com/v3/smtp/email", sendinblue)
	r.Run(":" + relayPort)
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

	sendMailv2(to, cc, body.Personalizations[0].Subject, body.Content[0].Value, body.Attanctment)
	c.JSON(200, gin.H{"status": "OK"})
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

func sendMailv2(to []string, cc []string, subject, message string, base []ChildAttactment) {
	to = []string{to[0]}
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to[0])
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", message)

	if len(base) > 0 {
		dec, err1 := base64.StdEncoding.DecodeString(base[0].Content)
		if err1 != nil {
			panic(err1)
		}

		f, err1 := os.Create(fmt.Sprintf("%s.%s", base[0].Filename, base[0].Type))
		if err1 != nil {
			panic(err1)
		}
		defer f.Close()

		if _, err := f.Write(dec); err != nil {
			panic(err)
		}
		if err := f.Sync(); err != nil {
			panic(err)
		}
		m.Attach(fmt.Sprintf("./%s.%s", base[0].Filename, base[0].Type))
	}

	port, _ := strconv.Atoi(port)
	d := gomail.NewDialer(host, port, user, pass)

	if err := d.DialAndSend(m); err != nil {
		log.Println(err)
	}
	if len(base) > 0 {
		e := os.Remove(fmt.Sprintf("./%s.%s", base[0].Filename, base[0].Type))
		if e != nil {
			log.Fatal(e)
		}
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
	Attanctment      []ChildAttactment  `json:"attachments"`
}

// Personalizations for Sendgrid
type Personalizations struct {
	To      []Person `json:"to"`
	Subject string   `json:"subject"`
}

// SendgridContent for email body
type SendgridContent struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type ChildAttactment struct {
	Content  string `json:"content"`
	Type     string `json:"type"`
	Filename string `json:"filename"`
}
