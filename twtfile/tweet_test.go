package twtfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*func TestScan(t *testing.T) {

	Scan("nick", `2022-09-11T12:32:34Z (#kjsrlnq) @<mckinley https://twtxt.net/user/mckinley/twtxt.txt>  Haha no Twts ingested that are found to be in the future are just silently dropped on the floor 🤣`)
	Scan("nick", `2022-09-11 12:32:34Z (#kjsrlnq) @<mckinley https://twtxt.net/user/mckinley/twtxt.txt>  Haha no Twts ingested that are found to be in the future are just silently dropped on the floor 🤣`)

	Scan("nick", `2022-09-11T12:32:34Z @<mckinley https://twtxt.net/user/mckinley/twtxt.txt>  Haha no Twts ingested that are found to be in the future are just silently dropped on the floor 🤣`)
	Scan("nick", `2022-09-18T10:58:14-04:00 Haha no Twts ingested that are found to be in the future are just silently dropped on the floor 🤣`)
}
*/

func TestParseTweet(t *testing.T) {

	// Arrange
	test_data := []struct {
		line     string
		expected *Tweet
		message  string
	}{
		{
			line:     "",
			expected: nil,
			message:  "1. Shoudln't parse empty lines",
		},
		{
			line:     "it's a commented line",
			expected: nil,
			message:  "2. Shoudln't parse commented lines",
		},
		{
			line:     "Hello World",
			expected: nil,
			message:  "3. Shoudln't parse non well formated tweet",
		},
		{
			line:     "2022-09-11T12:32:34Z Hello World",
			expected: &Tweet{Nick: "Nick", URL: "url", Hash: "", Text: "Hello World"},
			message:  "4. Should parse tweet",
		},
		{
			line:     "2022-09-11T12:32:34Z (#mn7yqma) Hello Nick",
			expected: &Tweet{Nick: "Nick", URL: "url", Hash: "(#mn7yqma)", Text: "Hello Nick"},
			message:  "5. Should parse reply",
		},
	}

	// Act
	for _, tt := range test_data {

		actual := parseTweet("Nick", "url", tt.line)

		if tt.expected == nil {
			assert.Nil(t, actual)
		} else {
			assert.Equal(t, tt.expected.Nick, actual.Nick, tt.message)
			assert.Equal(t, tt.expected.URL, actual.URL, tt.message)
			assert.Equal(t, tt.expected.Hash, actual.Hash, tt.message)
			assert.Equal(t, tt.expected.Text, actual.Text, tt.message)
			assert.False(t, actual.tweeting, tt.message)
		}
	}
}
