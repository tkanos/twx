package twtfile

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var twtxtDefaultTest = `# Twtxt is an open, distributed microblogging platform that
# uses human-readable text files, common transport protocols,
# and free software.
#
# Learn more about twtxt at  https://github.com/buckket/twtxt
#
# nick        = Nick
# url         = https://example.com/user/Nick/twtxt.txt
# avatar      = https://example.com/user/Nick/avatar#ope5cuoscdenqdddnuqexq2mfpjkty3en3agialamcebtvybfvqa
# description = My Life
#
# followers   = 1
# following   = 2
#
# link = linkedin https://www.linkedin.com/in/Nick-2b0842542/
#
# follow = hacker-news https://feeds.twtxt.net/hacker-news/twtxt.txt
# follow = Nick http://example.com/Nick/twtxt.txt
#

2022-03-09T16:09:43Z	hello world
2022-03-09T16:09:44Z	hello everybody
2022-03-09T16:09:45Z	hello all
2022-03-09T16:09:46Z	hello gang
`

func TestUnfollow(t *testing.T) {

	// Arrange
	test_data := []struct {
		f         func() string
		err       bool
		following int
		message   string
	}{
		{
			f: func() string {
				file, _ := ioutil.TempFile("", "twtxt*.txt")
				return file.Name()
			},
			following: 0,
			message:   "1. when file is empty we should get metadata",
		},
		{
			f: func() string {
				file, _ := ioutil.TempFile("", "twtxt*.txt")
				os.WriteFile(file.Name(), []byte(twtxtDefaultTest), 0644)
				return file.Name()
			},
			following: 1,
			message:   "2. when file is valid we should get one less follower",
		},
		{
			f: func() string {
				file, _ := ioutil.TempFile("", "twtxt*.txt")
				os.WriteFile(file.Name(), []byte(`2022-03-09T16:09:43Z	hello world
						2022-03-09T16:09:44Z	hello everybody
						2022-03-09T16:09:45Z	hello all
						2022-03-09T16:09:46Z	hello gang`), 0644)
				return file.Name()
			},
			following: 0,
			message:   "3. when filedoesn't have yet metadata it should have after",
		},
	}

	for _, tt := range test_data {
		// Arrange
		file := tt.f()
		defer os.Remove(file)

		twt, _ := NewTwtFile("Nick", file, "url")

		// Act
		twt.Unfollow("Nick")

		// Assert
		_, ok := twt.Meta.Follow["Nick"]
		assert.Equal(t, false, ok, tt.message)
		assert.Equal(t, tt.following, twt.Meta.Following, tt.message)

	}
}

func TestFollow(t *testing.T) {

	// Arrange
	test_data := []struct {
		f         func() string
		err       bool
		following int
		message   string
	}{
		{
			f: func() string {
				file, _ := ioutil.TempFile("", "twtxt*.txt")
				return file.Name()
			},
			following: 1,
			message:   "1. when file is empty we should get metadata",
		},
		{
			f: func() string {
				file, _ := ioutil.TempFile("", "twtxt*.txt")
				os.WriteFile(file.Name(), []byte(twtxtDefaultTest), 0644)
				return file.Name()
			},
			following: 3,
			message:   "2. when file is valid we should get one more follower",
		},
		{
			f: func() string {
				file, _ := ioutil.TempFile("", "twtxt*.txt")
				os.WriteFile(file.Name(), []byte(`2022-03-09T16:09:43Z	hello world
						2022-03-09T16:09:44Z	hello everybody
						2022-03-09T16:09:45Z	hello all
						2022-03-09T16:09:46Z	hello gang`), 0644)
				return file.Name()
			},
			following: 1,
			message:   "3. when file doesn't have yet metadata it should have after",
		},
	}

	for _, tt := range test_data {
		// Arrange
		file := tt.f()
		defer os.Remove(file)

		twt, _ := NewTwtFile("Nick", file, "url")

		// Act
		twt.Follow("Nick2", "http://example.com/Nick2/twtxt.txt", false)

		// Assert
		_, ok := twt.Meta.Follow["Nick2"]
		assert.Equal(t, true, ok, tt.message)
		assert.Equal(t, tt.following, twt.Meta.Following, tt.message)
	}
}
