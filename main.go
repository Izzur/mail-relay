package main

import (
	"fmt"
	"net/smtp"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/health", health)
	r.POST("/api.sendgrid.com/v3/mail/send", sendSMTP)
	r.POST("/api.sendinblue.com/v3/smtp/email", sendSMTP)
	r.Run()
}

func health(c *gin.Context) {
	c.JSON(200, gin.H{"status": "OK"})
}

func sendSMTP(c *gin.Context) {
	user := ""
	pass := ""
	host := "smtp.mailtrap.io"
	auth := smtp.PlainAuth("", user, pass, host)
	to := []string{"recipient@example.net"}
	msg := []byte("To: recipient@example.net\r\n" +
		"Subject: discount Gophers!\r\n" +
		"\r\n" +
		"This is the email body.\r\n")
	err := smtp.SendMail(host+":2525", auth, "sender@example.org", to, msg)
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(202, c.Request.Body)
}
