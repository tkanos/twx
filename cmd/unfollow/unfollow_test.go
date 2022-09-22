package unfollow

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tkanos/twx/cmd/context"
	"github.com/tkanos/twx/config"
	"github.com/tkanos/twx/twtfile"
)

func TestRemoveFollowingInConfig(t *testing.T) {

	// Arrange
	test_data := []struct {
		f       func() string
		message string
	}{
		{
			f: func() string {
				return tmpFileWithContent(t, []byte(`[TWTXT]
Nick="nick"

[Following]
Nick = "http://example.com/Nick/twtxt.txt"`))
			},
			message: "1. Config should remove followed user",
		},
		{
			f: func() string {
				return tmpFileWithContent(t, []byte(`[TWTXT]
Nick="nick"`))
			},
			message: "2. Config should not have user, if the user was not followed anyway",
		},
	}

	for _, tt := range test_data {
		// Arrange
		file := tt.f()
		defer os.Remove(file)

		f := unfollow{}
		context.Config, _ = config.NewConfig(file)

		// Act
		err := f.removeFollowingInConfig("Nick")

		// Assert
		assert.NoErrorf(t, err, tt.message)

		body, _ := ioutil.ReadFile(file)
		assert.NotContainsf(t, string(body), `nick = "http://example.com/Nick/twtxt.txt"`, tt.message)

	}
}

func tmpFileWithContent(t *testing.T, content []byte) string {
	file, err := ioutil.TempFile("", "config_test_data_")
	if err != nil {
		t.Error("Error creating file with test data", err)
		t.FailNow()
	}

	_, err = io.Copy(file, bytes.NewReader(content))
	if err != nil {
		t.Error("Error writing test data", err)
		t.FailNow()
	}

	newName := file.Name() + ".toml"
	err = os.Rename(file.Name(), newName)
	if err != nil {
		t.Error("Error renaming test data file", err)
		t.FailNow()
	}

	return newName
}

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

func TestRemoveFollowingInMetadataOfTwtxtFile(t *testing.T) {

	// Arrange
	test_data := []struct {
		f         func() string
		err       bool
		following string
		message   string
	}{
		{
			f: func() string {
				file, _ := ioutil.TempFile("", "twtxt*.txt")
				return file.Name()
			},
			following: "following   = 0",
			message:   "1. when file is empty we should get metadata",
		},
		{
			f: func() string {
				file, _ := ioutil.TempFile("", "twtxt*.txt")
				os.WriteFile(file.Name(), []byte(twtxtDefaultTest), 0644)
				return file.Name()
			},
			following: "following   = 1",
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
			following: "following   = 0",
			message:   "3. when filedoesn't have yet metadata it should have after",
		},
	}

	for _, tt := range test_data {
		// Arrange
		file := tt.f()
		defer os.Remove(file)

		f := unfollow{}
		context.TwtFile, _ = twtfile.NewTwtFile("Nick", file, "url")
		context.Config = &config.Configuration{
			Twtxt: config.TwtxtConfig{
				DiscloseIdentity: true,
			},
		}

		// Act
		err := f.removeFollowingInMetadataOfTwtxtFile("Nick")

		// Assert
		assert.Equal(t, err != nil, tt.err, tt.message)

		body, _ := ioutil.ReadFile(file)
		assert.NotContainsf(t, string(body), "Nick http://example.com/Nick/twtxt.txt", tt.message)
		assert.Containsf(t, string(body), tt.following, tt.message)

	}
}
