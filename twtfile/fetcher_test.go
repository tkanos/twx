package twtfile

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var templateTest string = `# Twtxt is an open, distributed microblogging platform that
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
# following   = 2
#
# link = Linkedin http://www.linkedin.com/My/Project
#
# follow = hacker-news https://feeds.twtxt.net/hacker-news/twtxt.txt
#
`

func TestFetchFile(t *testing.T) {

	// Arrange
	file, _ := ioutil.TempFile("", "twtxt*.txt")
	os.WriteFile(file.Name(), []byte(templateTest), 0644)
	defer os.Remove(file.Name())

	//Act
	_, meta, err := fetchFile("nick", file.Name(), "url")

	//Assert
	assert.NoError(t, err)
	assert.Equal(t, "Nick", meta.Nick)
	assert.Equal(t, "http://my.social/twtxt.txt", meta.URL)
	assert.Equal(t, "http://my.social/myavatar.png", meta.Avatar)
	assert.Equal(t, "That is just a unit test bro", meta.Description)
	assert.Equal(t, 100, meta.Followers)
	assert.Equal(t, 2, meta.Following)
	assert.Equal(t, 1, len(meta.Follow))
	assert.Equal(t, 1, len(meta.Link))

}

//Mock Http And Fetch

//Mock Gemini and Fetch
