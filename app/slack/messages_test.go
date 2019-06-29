package slack

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestFindChannelMentions(t *testing.T) {
	expected := []string{}
	actual := FindChannelMentions("foo")
	assert.Equal(t, expected, actual)

	expected = []string{"C012ABCDE"}
	actual = FindChannelMentions("ask <@U012ABCDEF> to bake a birthday cake for <@U345GHIJKL> in <#C012ABCDE>")
	assert.Equal(t, expected, actual)
}

func TestChannelRegex(t *testing.T) {
	actual := regexp.MustCompile(escapedChannelRegex).FindStringSubmatch("<#C012ABCDE|dsocial>")
	expected := []string{"<#C012ABCDE|dsocial>", "C012ABCDE", "|dsocial", "dsocial"}
	assert.Equal(t, expected, actual)

	actual = regexp.MustCompile(escapedChannelRegex).FindStringSubmatch("<#C012ABCDE>")
	expected = []string{"<#C012ABCDE>", "C012ABCDE", "", ""}
	assert.Equal(t, expected, actual)

	actual = regexp.MustCompile(escapedChannelRegex).FindStringSubmatch("<#C012ABCDE|d-social>")
	expected = []string{"<#C012ABCDE|d-social>", "C012ABCDE", "|d-social", "d-social"}
	assert.Equal(t, expected, actual)
}

func TestFindChannelMentionsWithMatchedText(t *testing.T) {
	expected := []mention{}
	actual := findChannelMentionsWithMatchedText("foo")
	assert.Equal(t, expected, expected)

	expected = []mention{{id: "C012ABCDE", matchedText: "<#C012ABCDE>"}}
	actual = findChannelMentionsWithMatchedText("ask <@U012ABCDEF> to bake a birthday cake for <@U345GHIJKL> in <#C012ABCDE>")
	assert.Equal(t, expected, actual)

	expected = []mention{{id: "C012ABCDE", matchedText: "<#C012ABCDE|d-social>", name: "d-social"}}
	actual = findChannelMentionsWithMatchedText("ask <@U012ABCDEF> to bake a birthday cake for <@U345GHIJKL> in <#C012ABCDE|d-social>")
	assert.Equal(t, expected, actual)
}

func TestFindUserMentions(t *testing.T) {
	expected := []string{}
	actual := FindUserMentions("foo")
	assert.Equal(t, expected, actual)

	expected = []string{"U012ABCDEF", "U345GHIJKL"}
	actual = FindUserMentions("ask <@U012ABCDEF> to bake a birthday cake for <@U345GHIJKL> in <#C012ABCDE>")
	assert.Equal(t, expected, actual)

	expected = []string{"U012ABCDEF", "U345GHIJKL"}
	actual = FindUserMentions("ask <@U012ABCDEF|someuser> to bake a birthday cake for <@U345GHIJKL> in <#C012ABCDE>")
	assert.Equal(t, expected, actual)
}

func TestFindEscapedUserAndChannelMentions(t *testing.T) {
	expected := map[string][]string{
		"users":    {},
		"channels": {},
	}
	actual := FindEscapedUserAndChannelMentions("foo")
	assert.Equal(t, actual, expected)

	expected = map[string][]string{
		"users":    {"<@U012ABCDEF>", "<@U345GHIJKL>"},
		"channels": {"<#C012ABCDE>"},
	}
	actual = FindEscapedUserAndChannelMentions("ask <@U012ABCDEF> to bake a birthday cake for <@U345GHIJKL> in <#C012ABCDE>")
	assert.Equal(t, expected, actual)
}

func TestUnescapeUserAndChannelMentions(t *testing.T) {
	expected := "ask @crushermd to bake a birthday cake for @worf in #d-social"
	mapping := map[string]string{
		"U012ABCDEF": "crushermd",
		"U345GHIJKL": "worf",
		"C012ABCDE":  "d-social",
	}
	actual := UnescapeUserAndChannelMentions("ask <@U012ABCDEF> to bake a birthday cake for <@U345GHIJKL> in <#C012ABCDE>", mapping)
	assert.Equal(t, expected, actual)

	actual = UnescapeUserAndChannelMentions("ask <@U012ABCDEF|crushermd> to bake a birthday cake for <@U345GHIJKL|worf> in <#C012ABCDE|d-social>", mapping)
	assert.Equal(t, expected, actual)

}
