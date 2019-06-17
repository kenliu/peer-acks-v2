package slack

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/kenliu/peer-acks-v2/app/dataaccess"
	"github.com/nlopes/slack"
	"log"
	"net/http"
	"os"
)

func PostAckToSlack(channelID string, message string) error {
	api := slack.New(os.Getenv("SLACK_OAUTH_TOKEN"))
	channelID, timestamp, err := api.PostMessage(channelID, slack.MsgOptionText(message, false))
	if err != nil {
		log.Println(err)
	}
	log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	return err
}

func HandleSlashCommand(message string, c *gin.Context, err error, userId string, db *sql.DB) error {
	slackApi := slack.New(os.Getenv("SLACK_OAUTH_TOKEN"))
	if message == "" || message == "help" {
		showSlashHelpText(c)
	} else {
		var user *slack.User
		user, err = slackApi.GetUserInfo(userId)
		if err != nil {
			log.Println(err)
		}

		log.Printf("ID: %s, Fullname: %s, Email: %s\n", user.ID, user.Profile.RealName, user.Profile.Email)
		err = dataaccess.CreateAck(db, message, user.Profile.Email, dataaccess.SOURCE_SLACK)
		if err != nil {
			log.Println(err)
		}
		err = PostAckToSlack(os.Getenv("SLACK_ACKS_CHANNELID"), message)
	}
	return err
}

func showSlashHelpText(c *gin.Context) {
	const helpMessage = "_Use the_ `/ack` _command like this:_ `/ack shout out to somebody for doing something good` \n_You can do this from any channel in Slack._"
	c.String(http.StatusOK, "%s", helpMessage)
}
