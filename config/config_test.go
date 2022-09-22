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
