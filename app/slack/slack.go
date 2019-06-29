package slack

import (
	"database/sql"
	"encoding/json"
	"github.com/kenliu/peer-acks-v2/app/dataaccess"
	"github.com/nlopes/slack"
	"log"
	"net/http"
	"os"
)

const slashHelpText = "_Use the_ `/ack` _command like this:_ `/ack shout out to somebody for doing something good` \n_You can do this from any channel in Slack._"
const slashAckSentText = "_thanks for recognizing your fellow roacher!_"

func PostAckToSlack(channelID string, message string) error {
	channelID, timestamp, err := slackApi().PostMessage(channelID, slack.MsgOptionText(message, false))
	if err != nil {
		log.Println(err)
	}
	log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	return err
}

func GetChannelInfo(channelId string) (*slack.Channel, error) {
	log.Println("getting channel info for channelId:", channelId)
	channel, err := slackApi().GetChannelInfo(channelId)
	if err != nil {
		log.Println(err)
	}
	return channel, err
}

func GetUserInfo(userId string) (*slack.User, error) {
	log.Println("getting user info for userId:", userId)
	user, err := slackApi().GetUserInfo(userId)
	if err != nil {
		log.Println(err)
	}
	return user, err
}

func LookupUserAndChannelNames(escapedUsers []string, escapedChannels []string) (map[string]string, error) {
	names := make(map[string]string) //map of user/channel ids to names

	for _, escapedUser := range escapedUsers {
		userId := UnescapeUserId(escapedUser)
		user, err := GetUserInfo(userId)
		if err != nil {
			return nil, err
		}
		names[userId] = user.Name
	}

	for _, escapedChannel := range escapedChannels {
		channelId := UnescapeChannelId(escapedChannel)
		channel, err := GetChannelInfo(channelId)
		if err != nil {
			return nil, err
		}
		names[channelId] = channel.Name
	}
	return names, nil
}

func HandleChallengeEvent(body []byte) (string, error) {
	var request map[string]string
	err := json.Unmarshal(body, &request)
	return request["challenge"], err
}

func HandleSlashCommand(message string, userId string, db *sql.DB) (string, error) {
	var err error
	slackApi := slackApi()
	var responseMessage string
	if message == "" || message == "help" {
		responseMessage = slashHelpText
	} else {
		var user *slack.User
		user, err = slackApi.GetUserInfo(userId)
		if err != nil {
			log.Println(err)
		}

		// Slack messages can have escaped user ids/channel names, so we have to unescape them
		// so they can be displayed properly on the web.
		// see https://api.slack.com/slash-commands#creating_commands

		escaped := FindEscapedUserAndChannelMentions(message)
		mappings, err := LookupUserAndChannelNames(escaped["users"], escaped["channels"])
		//TODO handle error here
		unescaped := UnescapeUserAndChannelMentions(message, mappings)
		log.Println("unescaped slack peerack:", unescaped)

		log.Printf("ID: %s, Fullname: %s, Email: %s\n", user.ID, user.Profile.RealName, user.Profile.Email)

		//TODO write unescaped ack to DB
		err = dataaccess.CreateAck(db, message, unescaped, user.Profile.Email, dataaccess.SOURCE_SLACK)
		if err != nil {
			log.Println(err)
		}
		err = PostAckToSlack(os.Getenv("SLACK_ACKS_CHANNELID"), message)
		//TODO handle error here
		responseMessage = slashAckSentText
	}
	return responseMessage, err
}

func ValidateRequestSignature(headers http.Header, body []byte, secret string) error {
	sv, _ := slack.NewSecretsVerifier(headers, secret)
	sv.Write(body)
	return sv.Ensure()
}

// support for verification tokens is deprecated in Slack, but it's a quick way to add authorization
func ValidateVerificationToken(requestToken string) bool {
	secretToken := os.Getenv("SLACK_VERIFICATION_TOKEN")
	if secretToken == "" {
		log.Fatal("SLACK_VERIFICATION_TOKEN environment variable not set")
	}
	return secretToken == requestToken
}

func slackApi() *slack.Client {
	return slack.New(os.Getenv("SLACK_OAUTH_TOKEN"))
}
