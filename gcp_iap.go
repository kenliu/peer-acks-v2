package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strings"
)

// GCP IAP sets this header for logged in users
const GoogleIapUserHeader = "x-goog-authenticated-user-email"

// email address only used for local testing (not running in IAP environment)
const DevEmailAddress = "dev.email@cockroachlabs.com"

func getUserEmail(c *gin.Context) string {
	//log.Println("called getUserEmail()")
	var email string

	//check to see if we're running in a local environment and set a dummy user email
	if os.Getenv("ENVIRONMENT") == "development" {
		email = DevEmailAddress
	} else if c.GetHeader(GoogleIapUserHeader) != "" {
		email = c.GetHeader(GoogleIapUserHeader)
		log.Println("detected logged in email header: " + email)
		email = strings.ReplaceAll(email, "accounts.google.com:", "")
		log.Println("detected logged in email: " + email)
	} else {
		log.Fatal("user email address not detected. app configuration problem?")
	}
	return email
}
