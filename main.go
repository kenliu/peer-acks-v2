package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nlopes/slack"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

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

	// set up request handlers
	router := gin.Default()
	bindStaticRoutes(router)

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", nil)
	})

	router.GET("/acks", func(c *gin.Context) {
		messages := fetchAcks(db, "")
		c.JSON(http.StatusOK, gin.H{"acks": messages})
	})

	router.POST("/acks", func(ctx *gin.Context) {
		message := ctx.PostForm("message")
		log.Println(message)

		senderEmail := getUserEmail(ctx)
		createAck(db, message, senderEmail)

		ctx.HTML(http.StatusOK, "ack_submitted.tmpl", nil)
	})

	//TODO implement delete
	//router.DELETE("/acks/:id", func(c *gin.Context) {
	//})

	// my acks page
	router.GET("/myacks", func(c *gin.Context) {
		messages := fetchAcks(db, getUserEmail(c))
		c.HTML(http.StatusOK, "myacks.tmpl", gin.H{"acks": messages})
	})

	// report page
	router.GET("/report", func(c *gin.Context) {
		messages := fetchAcks(db, "")
		c.HTML(http.StatusOK, "report.tmpl", gin.H{"acks": messages})
	})

	// liveness/readiness probe
	router.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// slack slash command
	router.POST("/slack/slashcommand", func(c *gin.Context) {
		userName := c.PostForm("user_name")
		message := c.PostForm("text")
		userId := c.PostForm("user_id")

		log.Println("received slack slash command from user:" + userName + " with message text: \"" + message + "\"")

		//for key, value := range c.Request.PostForm {
		//	log.Println(key, value)
		//}

		slackApi := slack.New(os.Getenv("SLACK_OAUTH_TOKEN"))
		if message == "" || message == "help" {
			showSlashHelpText(c)
		} else {
			user, err := slackApi.GetUserInfo(userId)
			if err != nil {
				log.Printf("%s\n", err)
				//FIXME handle this better
				return
			}
			log.Printf("ID: %s, Fullname: %s, Email: %s\n", user.ID, user.Profile.RealName, user.Profile.Email)
			err = createAck(db, message, user.Profile.Email)
			postAckToSlack(os.Getenv("SLACK_ACKS_CHANNELID"), message)

			if err == nil {
				c.String(http.StatusOK, "%s", "_thanks for recognizing your peer!_")
			} else {
				c.Status(http.StatusInternalServerError)
			}
		}
	})

	// run server on configured port
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = ":8080" //default
	}
	router.Run(serverPort)
}

func showSlashHelpText(c *gin.Context) {
	const helpMessage = "_Use the_ `/ack` _command like this:_ `/ack shout out to somebody for doing something good` \n_You can do this from any channel in Slack._"
	c.String(http.StatusOK, "%s", helpMessage)
}

func createAck(db *sql.DB, message string, senderEmail string) error {
	//TODO handle case where msg is empty
	//TODO trim spaces in message
	query := "INSERT INTO acks (msg, sender_email, updated_at) values ($1, $2, current_timestamp)"
	_, err := db.Exec(query, message, senderEmail)
	if err != nil {
		log.Println(err) //FIXME handle this better
	}
	return err
}

func bindStaticRoutes(router *gin.Engine) {
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
		log.Println("detected logged in email header: " + email)
		email = strings.ReplaceAll(email, "accounts.google.com:", "")
		log.Println("detected logged in email: " + email)
	}
	return email
}

// empty senderEmail string queries for all acks
func fetchAcks(db *sql.DB, senderEmail string) []string {
	log.Println("called fetchAcks()")
	messages := make([]string, 0)
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
		log.Println(err)
	}
	defer rows.Close()

	var message string
	log.Println("fetching acks:")
	for rows.Next() {
		err := rows.Scan(&message)
		if err != nil {
			log.Println(err)
		}

		log.Println(message)
		messages = append(messages, message)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
	}

	return messages
}

func postAckToSlack(channelID string, message string) error {
	api := slack.New(os.Getenv("SLACK_OAUTH_TOKEN"))
	channelID, timestamp, err := api.PostMessage(channelID, slack.MsgOptionText(message, false))
	if err != nil {
		log.Println(err)
	}
	log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	return err
}

func filterMultilineAcksForWeb(messages []string) []string {
	filtered := make([]string, len(messages))
	for i, v := range messages {
		filtered[i] = strings.ReplaceAll(v, "\n", "<br>")
	}
	return filtered
}
