package config

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig_Default(t *testing.T) {
	// arrange
	config = nil
	content := []byte(`
[twtxt] 
	nick = "Nick"`)
	filePath := tmpFileWithContent(t, content)
	defer os.Remove(filePath)

	// act
	conf, _ := NewConfig(filePath)

	// assert
	assert.NotNil(t, conf)
	assert.Equal(t, 20, conf.Twtxt.LimitTimeline)
}

func TestGetConfig_EnvVariable(t *testing.T) {
	// arrange
	config = nil
	content := []byte(`
[twtxt] 
	nick = "Nick"`)
	filePath := tmpFileWithContent(t, content)
	defer os.Remove(filePath)
	os.Setenv("TWTXT_NICK", "Nicolas")

	// act
	conf, _ := NewConfig(filePath)

	// assert
	assert.NotNil(t, conf)
	if config != nil {
		assert.Equal(t, "Nicolas", conf.Twtxt.Nick)
	}
}

func TestGetConfig_LocalFile(t *testing.T) {
	// arrange
	config = nil
	content := []byte(`
[twtxt] 
	nick = "Nick"
	limit_timeline = 40
	timeout = "10s"

[following] 
	A = "http://a.com/twtxt.txt"
	B = "http://b.com/twtxt.txt"`)

	filePath := tmpFileWithContent(t, content)
	defer os.Remove(filePath)

	// act
	conf, _ := NewConfig(filePath)

	// assert
	assert.NotNil(t, conf)
	if config != nil {
		assert.Equal(t, 40, config.Twtxt.LimitTimeline)
		assert.Equal(t, 10*time.Second, config.TimeoutDuration())
		assert.Len(t, conf.Following, 2)
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

func TestUnfollow(t *testing.T) {

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

		c, _ := NewConfig(file)

		// Act
		c.Unfollow("Nick")

		// Assert
		_, ok := c.Following["Nick"]
		assert.Equal(t, false, ok, tt.message)

	}
}

func TestFollow(t *testing.T) {

	// Arrange
	test_data := []struct {
		f       func() string
		err     bool
		message string
	}{
		{
			f: func() string {
				return tmpFileWithContent(t, []byte(`[TWTXT]
				Nick="nick"`))
			},
			message: "1. Config should have the followed user",
		},
		{
			f: func() string {
				return tmpFileWithContent(t, []byte(`[TWTXT]
Nick="nick"

[Following]
Nick = "http://example.com/Nick/twtxt.txt"`))
			},
			err:     false,
			message: "2. Config should not rewrite user",
		},
	}

	for _, tt := range test_data {
		// Arrange
		file := tt.f()
		defer os.Remove(file)

		c, _ := NewConfig(file)

		// Act
		c.Follow("Nick", "http://example.com/Nick/twtxt.txt", false)

		// Assert
		_, ok := c.Following["Nick"]
		assert.Equal(t, true, ok, tt.message)

	}
}
