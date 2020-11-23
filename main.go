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
	s, err := smtp.Dial("localhost:2525")
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{})
	}
	fmt.Println(s)
	err = s.Quit()
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{})
	}
	c.JSON(202, c.Request.Body)
}