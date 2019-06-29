package slack

import (
	"regexp"
	"strings"
)

type mention struct {
	matchedText string
	id          string
	name        string
}

func FindEscapedUserAndChannelMentions(msg string) map[string][]string {
	mentions := make(map[string][]string)

	userMatches := findUserMentionsWithMatchedText(msg)
	//if len(userMatches) > 0 {
	//	log.Println("found matched user mentions: ", userMatches)
	//}
	userMentions := make([]string, len(userMatches))
	for i, e := range userMatches {
		userMentions[i] = e.matchedText
	}
	mentions["users"] = userMentions

	channelMatches := findChannelMentionsWithMatchedText(msg)
	//if len(channelMatches) > 0 {
	//	log.Println("found matched channel mentions: ", userMatches)
	//}
	channelMentions := make([]string, len(channelMatches))
	for i, e := range channelMatches {
		channelMentions[i] = e.matchedText
	}
	mentions["channels"] = channelMentions
	return mentions
}

func UnescapeUserAndChannelMentions(msg string, userAndChannelNames map[string]string) string {
	unescapedString := msg
	channels := findChannelMentionsWithMatchedText(msg)
	users := findUserMentionsWithMatchedText(msg)

	for _, escapedUser := range users {
		unescaped := userAndChannelNames[escapedUser.id]
		//TODO error if nil
		unescapedString = strings.ReplaceAll(unescapedString, escapedUser.matchedText, "@"+unescaped)
	}

	for _, escapedChannel := range channels {
		unescaped := userAndChannelNames[escapedChannel.id]
		//TODO error if nil
		unescapedString = strings.ReplaceAll(unescapedString, escapedChannel.matchedText, "#"+unescaped)
	}
	return unescapedString
}

func FindUserMentions(msg string) []string {
	matches := findUserMentionsWithMatchedText(msg)
	users := make([]string, len(matches))
	for i, match := range matches {
		users[i] = match.id
	}
	return users
}

const escapedChannelRegex = `<#([\w\d]+)(\|([a-z0-9][a-z0-9._-]*))?>`
const escapedUserRegex = `<@([\w\d]+)(\|([a-z0-9][a-z0-9._-]*))?>`

func findUserMentionsWithMatchedText(msg string) []mention {
	//log.Println("trying to find user mentions: ", msg)
	return findMentionsWithMatchedText(escapedUserRegex, msg)
}

func FindChannelMentions(msg string) []string {
	matches := findChannelMentionsWithMatchedText(msg)
	channels := make([]string, len(matches))
	for i, match := range matches {
		channels[i] = match.id
	}
	return channels
}

func findChannelMentionsWithMatchedText(msg string) []mention {
	//log.Println("trying to find channel mentions: ", msg)
	return findMentionsWithMatchedText(escapedChannelRegex, msg)
}

// generalized function for finding escaped user or channel with a regex
func findMentionsWithMatchedText(regexWithGroup string, msg string) []mention {
	matches := regexp.MustCompile(regexWithGroup).FindAllStringSubmatch(msg, -1)
	userIds := make([]mention, len(matches))
	for i, match := range matches {
		//log.Println("match: ", match)
		m := mention{matchedText: match[0], id: match[1]}

		// if the 4th element is not empty then there were 3 matched groups
		// meaning that the name was part of the escaped text
		// e.g. <@U4FNDDT8T|bobby> instead of <@U4FNDDT8T>
		if match[3] != "" {
			m.name = match[3]
		}
		userIds[i] = m
	}
	return userIds
}

func UnescapeUserId(escapedUser string) string {
	return findUserMentionsWithMatchedText(escapedUser)[0].id
}

func UnescapeChannelId(escapedChannel string) string {
	return findChannelMentionsWithMatchedText(escapedChannel)[0].id
}
