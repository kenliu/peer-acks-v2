package slack

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/kenliu/peer-acks-v2/app/dataaccess"
	"github.com/nlopes/slack"
	"log"
	"os"
)

const SLASH_HELP_TEXT = "_Use the_ `/ack` _command like this:_ `/ack shout out to somebody for doing something good` \n_You can do this from any channel in Slack._"

func PostAckToSlack(channelID string, message string) error {
	api := slack.New(os.Getenv("SLACK_OAUTH_TOKEN"))
	channelID, timestamp, err := api.PostMessage(channelID, slack.MsgOptionText(message, false))
	if err != nil {
		log.Println(err)
	}
	log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	return err
}

func HandleSlashCommand(message string, c *gin.Context, err error, userId string, db *sql.DB) (string, error) {
	slackApi := slack.New(os.Getenv("SLACK_OAUTH_TOKEN"))
	var responseMessage string
	if message == "" || message == "help" {
		responseMessage = SLASH_HELP_TEXT
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
		responseMessage = "_thanks for recognizing your fellow roacher!_"
	}
	return responseMessage, err
}
