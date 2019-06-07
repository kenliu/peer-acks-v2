package main

import (
	"database/sql"
	"debug/elf"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	// initialize the DB
	log.Println("attempting to connect to DB")
	datasource := os.Getenv("DATASOURCE")
	db, err := sql.Open("postgres", datasource)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	defer db.Close()
	log.Println("DB connection successful")

	// set up request handlers
	router := gin.Default()
	createStaticRoutes(router)

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", fetchAcks(db, getUserEmail(c)))
	})

	router.GET("/acks", func(c *gin.Context) {
		c.JSON(http.StatusOK, fetchAcks(db, getUserEmail(c)))
	})

	router.POST("/acks", func(ctx *gin.Context) {
		message := ctx.PostForm("message")
		log.Println(message)

		//TODO handle case where msg is empty
		//TODO trim spaces in message
		query := "INSERT INTO acks (msg, sender_email, updated_at) values ($1, $2, current_timestamp)"
		_, err := db.Exec(query, message, getUserEmail(ctx))
		if err != nil {
			log.Fatal(err) //TODO handle this better
		}

		ctx.HTML(http.StatusOK, "ack_submitted.tmpl", fetchAcks(db, getUserEmail(ctx)))
	})

	//TODO implement delete
	//router.DELETE("/acks/:id", func(c *gin.Context) {
	//})

	// my acks page
	router.GET("/myacks", func(c *gin.Context) {
		c.HTML(http.StatusOK, "myacks.tmpl", fetchAcks(db, getUserEmail(c)))
	})

	// report page
	router.GET("/report", func(c *gin.Context) {
		c.HTML(http.StatusOK, "report.tmpl", fetchAcks(db, ""))
	})

	// run server on configured port
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = ":8080" //default
	}
	router.Run(serverPort)
}

func createStaticRoutes(router *gin.Engine) {
	// create static routes
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")
	router.StaticFile("/radiator", "./templates/radiator.html")
	router.Static("/resources", "./resources")

	router.LoadHTMLGlob("templates/*")
}

// GCP IAP sets this header for logged in users
const GoogleIapUserHeader = "x-goog-authenticated-user-email"

// email address only used for local testing (not running in IAP environment)
const DevEmailAddress = "dev.email@cockroachlabs.com"

func getUserEmail(c *gin.Context) string {
	log.Println("called getUserEmail()")
	var email string

	//check to see if we're running in a local environment and set a dummy user email
	if os.Getenv("ENVIRONMENT") == "development" {
		email = DevEmailAddress
	} else if c.GetHeader(GoogleIapUserHeader) != "" {
		email = c.GetHeader(GoogleIapUserHeader)
		log.Println("detected logged in email: " + email)
	}
	return email
}

// empty senderEmail string queries for all acks
func fetchAcks(db *sql.DB, senderEmail string) gin.H {
	log.Println("called fetchAcks()")
	var messages []string
	var rows *sql.Rows
	var err error
	var query string

	//FIXME figure out how to properly pass timestamps to db.Query instead of concatenating a string here
	curDateString := fmt.Sprintf(time.Now().Format("2006-01-02"))
	if senderEmail == "" {
		log.Println("fetching all acks, last 7 days")
		query = "select msg from acks where created_at > ('" + curDateString + "' - 7) order by updated_at desc"
		rows, err = db.Query(query)
	} else {
		log.Println("fetching acks by user email: " + senderEmail)
		query = "select msg from acks where sender_email = $1 and created_at > ('" + curDateString + "' - 7) order by updated_at desc"
		rows, err = db.Query(query, senderEmail)
	}

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var message string
	log.Println("fetching acks:")
	for rows.Next() {
		err := rows.Scan(&message)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(message)
		messages = append(messages, message)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return gin.H{"acks": messages}
}
