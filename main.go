package main

import (
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
	c.JSON(202, c.Request.Body)
}