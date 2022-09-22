package quickstart

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetNick(t *testing.T) {
	test_data := []struct {
		f       func() io.Reader
		def     string
		changed bool
		message string
	}{
		{
			f: func() io.Reader {
				var stdin bytes.Buffer
				stdin.Write([]byte("\n"))

				return &stdin
			},
			def:     "nick",
			changed: false,
			message: "1. Nick shouldn't change if no input",
		},
		{
			f: func() io.Reader {
				var stdin bytes.Buffer
				stdin.Write([]byte("test\n"))

				return &stdin
			},
			def:     "nick",
			changed: true,
			message: "2. Nick should change if new input",
		},
	}

	for _, tt := range test_data {
		// Arrange
		q := quickstart{
			reader: tt.f(),
			nick:   tt.def,
		}

		// Act
		q.setNick()

		// Assert
		assert.Equal(t, tt.def == q.nick, !tt.changed, tt.message)
	}
}

func TestSetConfigFile(t *testing.T) {
	test_data := []struct {
		f       func() io.Reader
		def     string
		changed bool
		message string
	}{
		{
			f: func() io.Reader {
				var stdin bytes.Buffer
				stdin.Write([]byte("\n"))

				return &stdin
			},
			def:     "config",
			changed: false,
			message: "1. Config shouldn't change if no input",
		},
		{
			f: func() io.Reader {
				var stdin bytes.Buffer
				stdin.Write([]byte("config2\n"))

				return &stdin
			},
			def:     "config",
			changed: true,
			message: "2. Config should change if new input",
		},
	}

	for _, tt := range test_data {
		// Arrange
		q := quickstart{
			reader:         tt.f(),
			configFilePath: tt.def,
		}

		// Act
		q.setConfigFile()

		// Assert
		assert.Equal(t, tt.def == q.configFilePath, !tt.changed, tt.message)
	}
}

func TestSetTwtxtFile(t *testing.T) {
	test_data := []struct {
		f       func() io.Reader
		def     string
		changed bool
		message string
	}{
		{
			f: func() io.Reader {
				var stdin bytes.Buffer
				stdin.Write([]byte("\n"))

				return &stdin
			},
			def:     "twtxt.txt",
			changed: false,
			message: "1. twtxtfile shouldn't change if no input",
		},
		{
			f: func() io.Reader {
				var stdin bytes.Buffer
				stdin.Write([]byte("mine-twtxt.txt\n"))

				return &stdin
			},
			def:     "twtxt.txt",
			changed: true,
			message: "2. twtxtfile should change if new input",
		},
	}

	for _, tt := range test_data {
		// Arrange
		q := quickstart{
			reader:        tt.f(),
			twtxtFilePath: tt.def,
		}

		// Act
		q.setTwtxtFile()

		// Assert
		assert.Equal(t, tt.def == q.twtxtFilePath, !tt.changed, tt.message)
	}
}

func TestSetUrl(t *testing.T) {
	test_data := []struct {
		f       func() io.Reader
		def     string
		changed bool
		message string
	}{
		{
			f: func() io.Reader {
				var stdin bytes.Buffer
				stdin.Write([]byte("\n"))

				return &stdin
			},
			def:     "http://my.social/twtxt.txt",
			changed: false,
			message: "1. Url shouldn't change if no input",
		},
		{
			f: func() io.Reader {
				var stdin bytes.Buffer
				stdin.Write([]byte("http://my.social/mine-twtxt.txt\n"))

				return &stdin
			},
			def:     "http://my.social/twtxt.txt",
			changed: true,
			message: "2. Url should change if new input",
		},
	}

	for _, tt := range test_data {
		// Arrange
		q := quickstart{
			reader: tt.f(),
			url:    tt.def,
		}

		// Act
		q.setUrl()

		// Assert
		assert.Equal(t, tt.def == q.url, !tt.changed, tt.message)
	}
}

func TestSetDiscloseIdentity(t *testing.T) {
	test_data := []struct {
		f       func() io.Reader
		def     bool
		changed bool
		message string
	}{
		{
			f: func() io.Reader {
				var stdin bytes.Buffer
				stdin.Write([]byte("\n"))

				return &stdin
			},
			def:     true,
			changed: false,
			message: "1. DiscloseIdentity shouldn't change if no input",
		},
		{
			f: func() io.Reader {
				var stdin bytes.Buffer
				stdin.Write([]byte("N\n"))

				return &stdin
			},
			def:     true,
			changed: true,
			message: "2. DiscloseIdentity should change if new input",
		},
	}

	for _, tt := range test_data {
		// Arrange
		q := quickstart{
			reader:           tt.f(),
			discloseIdentity: tt.def,
		}

		// Act
		q.setDiscloseIdentity()

		// Assert
		assert.Equal(t, tt.def == q.discloseIdentity, !tt.changed, tt.message)
	}
}

func TestSetFollowNews(t *testing.T) {
	test_data := []struct {
		f       func() io.Reader
		def     bool
		changed bool
		sizeMap int
		message string
	}{
		{
			f: func() io.Reader {
				var stdin bytes.Buffer
				stdin.Write([]byte("\n"))

				return &stdin
			},
			def:     true,
			changed: false,
			sizeMap: 1,
			message: "1. news shouldn't change if no input",
		},
		{
			f: func() io.Reader {
				var stdin bytes.Buffer
				stdin.Write([]byte("N\n"))

				return &stdin
			},
			def:     true,
			changed: true,
			sizeMap: 0,
			message: "2. news should change if new input",
		},
	}

	for _, tt := range test_data {
		// Arrange
		q := quickstart{
			reader: tt.f(),
			news:   tt.def,
		}

		// Act
		q.setFollowNews()

		// Assert
		assert.Equal(t, tt.def == q.news, !tt.changed, tt.message)
		assert.Equal(t, tt.sizeMap, len(q.follow), tt.message)
	}
}

var expectedConfig = `[twtxt]
nick              = Nick
twtfile           = %s/twtxt.txt
twturl            = http://my.social
disclose_identity = true

[following]
twtxt_news = https://buckket.org/twtxt_news.txt
`

func TestCreateIniConfig(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tmpDir) // clean up

	ec := fmt.Sprintf(expectedConfig, tmpDir)

	q := quickstart{
		nick:             "Nick",
		twtxtFilePath:    filepath.Join(tmpDir, "twtxt.txt"),
		url:              "http://my.social",
		discloseIdentity: true,
		news:             true,
		configFilePath:   filepath.Join(tmpDir, "config"),
		follow: map[string]string{
			"twtxt_news": "https://buckket.org/twtxt_news.txt",
		},
	}

	q.createIniConfig()

	body, err := ioutil.ReadFile(q.configFilePath)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}

	assert.Equal(t, ec, string(body), "1- Ini should be created")
}

func TestCreateIniConfigWhenWrongPath(t *testing.T) {
	q := quickstart{
		nick:             "Nick",
		twtxtFilePath:    "twtxt.txt",
		url:              "http://my.social",
		discloseIdentity: true,
		news:             true,
		configFilePath:   "Not/Existing/path/config",
		follow: map[string]string{
			"twtxt_news": "https://buckket.org/twtxt_news.txt",
		},
	}

	err := q.createIniConfig()

	assert.Error(t, err, "1- Ini should send an error, if it could not create a config file")
}

var expectedTwtxtFile string = `# Twtxt is an open, distributed microblogging platform that
# uses human-readable text files, common transport protocols,
# and free software.
#
# Learn more about twtxt at  https://github.com/buckket/twtxt
#
# nick        = Nick
# url         = http://my.social
# avatar      = 
# description = 
#
# followers   = 0
# following   = 1
#
# follow = twtxt_news https://buckket.org/twtxt_news.txt
# 
`

func TestCreateTwtxtFile(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tmpDir) // clean up

	q := quickstart{
		nick:             "Nick",
		twtxtFilePath:    filepath.Join(tmpDir, "twtxt.txt"),
		url:              "http://my.social",
		discloseIdentity: true,
		news:             true,
		configFilePath:   filepath.Join(tmpDir, "config"),
		follow: map[string]string{
			"twtxt_news": "https://buckket.org/twtxt_news.txt",
		},
	}

	err := q.createTwtxtFile()
	if err != nil {
		assert.Fail(t, "could not create twtxt file", err)
	}

	body, err := ioutil.ReadFile(q.twtxtFilePath)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}

	assert.Equal(t, expectedTwtxtFile, string(body), "1- twtxt file should be created")

}
