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

func TestParseMeta(t *testing.T) {

	// Arrange
	test_data := []struct {
		line     string
		expected TweetMetadata
		message  string
	}{
		{
			line:     "This is a line not strating by #",
			expected: TweetMetadata{},
			message:  "1. Shoudln't parse not strating by # lines",
		},
		{
			line:     "",
			expected: TweetMetadata{},
			message:  "2. Shoudln't parse empty lines",
		},
		{
			line:     "#",
			expected: TweetMetadata{},
			message:  "3. Shoudln't parse empty lines, starting by #",
		},
		{
			line:     "#  This is a comment",
			expected: TweetMetadata{},
			message:  "4. Shoudln't parse commented lines",
		},
		{
			line:     `#  nick = Nick`,
			expected: TweetMetadata{Nick: "Nick"},
			message:  "5. Shoud parse nick",
		},
		{
			line:     `#  Nick  =   Nick`,
			expected: TweetMetadata{Nick: "Nick"},
			message:  "6. Shoud parse Nick",
		},
		{
			line:     `#  url = http://myurl.com`,
			expected: TweetMetadata{URL: "http://myurl.com"},
			message:  "7. Shoud parse url",
		},
		{
			line:     `#  description = my description`,
			expected: TweetMetadata{Description: "my description"},
			message:  "8. Shoud parse desciption",
		},
		{
			line:     `#  avatar = http://myurl.com`,
			expected: TweetMetadata{Avatar: "http://myurl.com"},
			message:  "9. Shoud parse avatar",
		},
		{
			line:     `#  followers    =   1`,
			expected: TweetMetadata{Followers: 1},
			message:  "10. Shoud parse followers",
		},
		{
			line:     `#  following    =   1`,
			expected: TweetMetadata{Following: 1},
			message:  "11. Shoud parse following",
		},
		{
			line:     `#  link    =  me http://myurl.com`,
			expected: TweetMetadata{Link: map[string]string{"me": "http://myurl.com"}},
			message:  "12. Shoud parse link",
		},
		{
			line:     `#  follow    =  me http://myurl.com`,
			expected: TweetMetadata{Follow: map[string]string{"me": "http://myurl.com"}},
			message:  "12. Shoud parse follow",
		},
	}

	// Act
	for _, tt := range test_data {

		actual := TweetMetadata{}
		actual = parseMeta(tt.line, actual)

		assert.Equal(t, tt.expected, actual, tt.message)

	}

}
