package twtfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var expectedTemplate string = `# Twtxt is an open, distributed microblogging platform that
# uses human-readable text files, common transport protocols,
# and free software.
#
# Learn more about twtxt at  https://github.com/buckket/twtxt
#
# nick        = Nick
# url         = http://my.social/twtxt.txt
# avatar      = http://my.social/myavatar.png
# description = That is just a unit test bro
#
# followers   = 100
# following   = 1
#
# link = Linkedin http://www.linkedin.com/My/Project
#
# follow = hacker-news https://feeds.twtxt.net/hacker-news/twtxt.txt
# 
`

func TestCreateTwtxtTemplate(t *testing.T) {
	// Arrange
	meta := TweetMetadata{
		Nick:        "Nick",
		URL:         "http://my.social/twtxt.txt",
		Avatar:      "http://my.social/myavatar.png",
		Description: "That is just a unit test bro",
		Followers:   100,
		Following:   2,
		Link:        map[string]string{"Linkedin": "http://www.linkedin.com/My/Project"},
		Follow:      map[string]string{"hacker-news": "https://feeds.twtxt.net/hacker-news/twtxt.txt"},
	}

	// Act
	b, err := meta.CreateTwtxtMetaTemplate()

	// Assert
	if err != nil {
		assert.FailNow(t, "could not generate template twtxt", err)
	}

	assert.Equal(t, expectedTemplate, string(b), "twtxt file template should be generated")
}
