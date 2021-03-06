package main

import (
	"bytes"
	"database/sql"
	//	"errors"
	"github.com/gin-gonic/gin"
	"github.com/kenliu/peer-acks-v2/app/dataaccess"
	"github.com/kenliu/peer-acks-v2/app/slack"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type AppConfig struct {
	datasource         string
	environment        string
	slackOAuthToken    string
	slackAcksChannelID string
	serverPort         string
}

func main() {
	// initialize the DB connection
	log.Println("attempting to connect to DB")
	datasource := os.Getenv("DATASOURCE")
	db, err := sql.Open("postgres", datasource)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	defer db.Close()
	log.Println("DB connection successful")

	//sentry.Init(sentry.ClientOptions{
	//	Environment: "staging",
	//	Dsn: "https://93c7a505a8a14c659d087caf91f165f9@sentry.io/1483004",
	//})
	//
	//errors.New()
	//
	// set up request handlers
	router := gin.Default()
	bindStaticRoutes(router)

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", nil)
	})

	router.GET("/acks", func(c *gin.Context) {
		messages := dataaccess.FetchAcks(db, "")
		c.JSON(http.StatusOK, gin.H{"acks": messages})
	})

	router.POST("/acks", func(ctx *gin.Context) {
		var err error
		message := ctx.PostForm("message")
		senderEmail := getUserEmail(ctx)

		err = dataaccess.CreateAck(db, message, "", senderEmail, dataaccess.SOURCE_WEB)
		//FIXME handle error
		err = slack.PostAckToSlack(os.Getenv("SLACK_ACKS_CHANNELID"), message)

		if err != nil {
			ctx.Status(http.StatusInternalServerError)
		} else {
			ctx.HTML(http.StatusOK, "ack_submitted.tmpl", nil)
		}
	})

	//TODO implement delete
	//router.DELETE("/acks/:id", func(c *gin.Context) {
	//})

	// my acks page
	router.GET("/myacks", func(c *gin.Context) {
		messages := dataaccess.FetchAcks(db, getUserEmail(c))
		c.HTML(http.StatusOK, "myacks.tmpl", gin.H{"acks": messages})
	})

	// report page
	router.GET("/report", func(c *gin.Context) {
		messages := dataaccess.FetchAcks(db, "")
		c.HTML(http.StatusOK, "report.tmpl", gin.H{"acks": messages})
	})

	// liveness/readiness probe
	router.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	router.POST("/slack/events", func(c *gin.Context) {
		//TODO detect event type and discard unknown events, as a safety measure
		body, err := c.GetRawData()
		challenge, err := slack.HandleChallengeEvent(body)
		if err != nil {
			c.Status(http.StatusInternalServerError)
		} else {
			c.String(http.StatusOK, "%s", challenge)
		}
	})

	// slack slash command
	router.POST("/slack/slashcommand", func(c *gin.Context) {
		verr := slack.ValidateRequestSignature(c.Request.Header, repeatableReadBody(c), os.Getenv("SLACK_SIGNING_SECRET"))
		if verr != nil {
			log.Println(verr)
			log.Println("request signature verification failed")
			c.Status(http.StatusForbidden)
			return
		}

		var err error
		userName := c.PostForm("user_name")
		message := c.PostForm("text")
		userId := c.PostForm("user_id")
		log.Println("received slack slash command from user:" + userName + " with message text: \"" + message + "\"")

		//for key, value := range c.Request.PostForm {
		//	log.Println(key, value)
		//}

		var responseMessage string
		responseMessage, err = slack.HandleSlashCommand(message, userId, db)

		if err != nil {
			c.Status(http.StatusInternalServerError)
		} else {
			// return response to Slack API to display message to user
			c.String(http.StatusOK, "%s", responseMessage)
		}
	})

	// run server on configured port
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = ":8080" //default
	}
	router.Run(serverPort)
}

func repeatableReadBody(c *gin.Context) []byte {
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	}
	// Restore the io.ReadCloser to its original state
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes
}

func bindStaticRoutes(router *gin.Engine) {
	// create static routes
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")
	router.StaticFile("/radiator", "./templates/radiator.html")
	router.Static("/resources", "./resources")

	router.LoadHTMLGlob("templates/*")
}
