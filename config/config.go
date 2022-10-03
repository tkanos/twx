// Package config all the configurable variables within TJ-DSP.
//
// Configs can be loaded from an external .toml file as well as via environment variables.
// Defaults are set within this package.
//
// When both environment variables and config file variables are loaded, the order of precedence is:
//
// Environment > Config file > Default
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/naoina/toml"
	"github.com/tkanos/twx/utils"
)

// Configuration is top-level and represents all configuration settings.
type Configuration struct {
	path string
	// Twtxt
	Twtxt TwtxtConfig

	// Following
	Following map[string]string

	// Hook
	Hook map[string]string
}

// TwtxtConfig represents the twtxt config [https://twtxt.readthedocs.io/en/latest/user/configuration.html].
type TwtxtConfig struct {
	//nick : your nick, will be displayed in your timeline
	Nick string
	//twtfile: path to your local twtxt file
	TwtFile string
	//twturl: URL to your public twtxt file
	TwtURL string
	//check_following: try to resolve URLs when listing followings
	CheckFollowing bool
	//use_pager: use a pager (less) to display your timeline
	UsePager bool
	//use_cache : cache remote twtxt files locally
	UseCache bool
	//porcelain style output in an easy-to-parse format
	Porcelain bool
	//disclose_identity: include nick and twturl in twtxt’s user-agent
	DiscloseIdentity bool
	//incluse_yourself: to include yourself in the timeline or not.
	IncludeYourself bool
	//character_warning: warn when composed tweet has more characters
	CharacterWarning int
	//limit_timeline: limit amount of tweets shown in your timeline
	LimitTimeline int
	//timeline_update_interval: time in seconds cache is considered up-to-date
	TimelineUpdateInterval int
	//timeout: maximal time a http request is allowed to take
	Timeout         string
	timeoutDuration time.Duration
	//sorting: sort timeline either descending or ascending
	Sorting string
	//use_abs_time: use absolute datetimes in your timeline
	UseAbsTime bool
	//pre_tweet_hook: command to be executed before tweeting
	PreTweetHook string
	//post_tweet_hook: command to be executed after tweeting
	PostTweetHook string
	//timeline_ascii_images: Show images on ascii in console
	TimelineAsciiImages bool
}

func (c *Configuration) TimeoutDuration() time.Duration {
	return c.Twtxt.timeoutDuration
}

func (c *Configuration) ConfigPath() string {
	return c.path
}

func HomeDirectory() string {
	v, _ := os.UserHomeDir()
	return v
}

var config *Configuration

// GetConfig initializes and return the top-level Configuration struct.
func NewConfig(path string) (*Configuration, error) {
	if config == nil {
		var err error
		if path, err = setup(path); err != nil {
			return config, err
		}
	}

	config.path = path

	if config.Following == nil {
		config.Following = map[string]string{}
	}

	return config, nil
}

func (c *Configuration) Follow(nick, url string, replace bool) {
	if _, ok := c.Following[nick]; !ok || replace {
		c.Following[nick] = url
	}
}

func (c *Configuration) Unfollow(nick string) {
	delete(c.Following, nick)
}

func (c *Configuration) Save() error {
	//viper.Set("twtxt.nick", c.Twtxt.Nick)
	f, err := os.Create(c.path)
	if err != nil {
		// failed to create/open the file
		return err
	}
	if err := toml.NewEncoder(f).Encode(config); err != nil {
		// failed to encode
		return err
	}
	if err := f.Close(); err != nil {
		// failed to close the file
		return err
	}

	return nil
}

func setup(path string) (string, error) {
	config = &Configuration{}

	bindDefault()
	var err error
	if path, err = bindFile(path); err != nil {
		return path, err
	}

	//
	if config.Twtxt.Timeout != "5s" {
		if durationValue, err := time.ParseDuration(config.Twtxt.Timeout); err == nil {
			config.Twtxt.timeoutDuration = durationValue
		} else {
			config.Twtxt.Timeout = "5s"
		}
	}

	bindEnv()

	if len(config.Twtxt.TwtFile) > 0 && config.Twtxt.TwtFile[0] == '~' {
		config.Twtxt.TwtFile = utils.ExpandTilde(config.Twtxt.TwtFile)
	}

	return path, nil

}

func bindDefault() {
	// Twtxt
	config.Twtxt.CheckFollowing = true
	config.Twtxt.CheckFollowing = true
	config.Twtxt.UseCache = true
	config.Twtxt.CharacterWarning = -1
	config.Twtxt.LimitTimeline = 20
	config.Twtxt.TimelineUpdateInterval = 10
	config.Twtxt.Timeout = "5s"
	config.Twtxt.timeoutDuration = 5 * time.Second
	config.Twtxt.Sorting = "descending"

	config.Twtxt.UsePager = false
	config.Twtxt.Porcelain = false
	config.Twtxt.DiscloseIdentity = false
	config.Twtxt.UseAbsTime = false
}

func bindEnv() {

	if v := os.Getenv("TWTXT_NICK"); v != "" {
		config.Twtxt.Nick = v
	}
	if v := os.Getenv("TWTXT_TWTFILE"); v != "" {
		config.Twtxt.TwtFile = v
	}
	if v := os.Getenv("TWTXT_TWTURL"); v != "" {
		config.Twtxt.TwtURL = v
	}
	if v := os.Getenv("TWTXT_SORTING"); v != "" {
		if v == "descensing" || v == "ascending" {
			config.Twtxt.Sorting = v
		}
	}
	if v := os.Getenv("TWTXT_CHECK_FOLLOWING"); v != "" {
		if boolValue, err := strconv.ParseBool(v); err == nil {
			config.Twtxt.CheckFollowing = boolValue
		}
	}
	if v := os.Getenv("TWTXT_USE_PAGER"); v != "" {
		if boolValue, err := strconv.ParseBool(v); err == nil {
			config.Twtxt.UsePager = boolValue
		}
	}
	if v := os.Getenv("TWTXT_USE_CACHE"); v != "" {
		if boolValue, err := strconv.ParseBool(v); err == nil {
			config.Twtxt.UseCache = boolValue
		}
	}
	if v := os.Getenv("TWTXT_PORCELAIN"); v != "" {
		if boolValue, err := strconv.ParseBool(v); err == nil {
			config.Twtxt.Porcelain = boolValue
		}
	}
	if v := os.Getenv("TWTXT_DISCLOSE_IDENTITY"); v != "" {
		if boolValue, err := strconv.ParseBool(v); err == nil {
			config.Twtxt.DiscloseIdentity = boolValue
		}
	}
	if v := os.Getenv("TWTXT_USE_ABS_TIME"); v != "" {
		if boolValue, err := strconv.ParseBool(v); err == nil {
			config.Twtxt.UseAbsTime = boolValue
		}
	}
	if v := os.Getenv("TWTXT_CHARACTER_WARNING"); v != "" {
		if intValue, err := strconv.Atoi(v); err == nil {
			config.Twtxt.CharacterWarning = intValue
		}
	}
	if v := os.Getenv("TWTXT_LIMIT_TIMELINE"); v != "" {
		if intValue, err := strconv.Atoi(v); err == nil {
			config.Twtxt.LimitTimeline = intValue
		}
	}
	if v := os.Getenv("TWTXT_TIMELINE_UPDATE_INTERVAL"); v != "" {
		if intValue, err := strconv.Atoi(v); err == nil {
			config.Twtxt.TimelineUpdateInterval = intValue
		}
	}

	if v := os.Getenv("TWTXT_TIMEOUT"); v != "" {
		if durationValue, err := time.ParseDuration(v); err == nil {
			config.Twtxt.Timeout = v
			config.Twtxt.timeoutDuration = durationValue
		}
	}

}

func bindFile(filePath string) (string, error) {
	var f *os.File
	var path string
	filepaths := []string{
		filePath,
		fmt.Sprintf("%s/.config/twx/twtxt.toml", HomeDirectory()),
		fmt.Sprintf("%s/.twtxt.toml", HomeDirectory()),
		fmt.Sprintf("%s/Library/Application Support/twx", HomeDirectory()),
		"/etc/twx/conf.d/twtxt.toml",
	}
	for _, fp := range filepaths {
		if _, err := os.Stat(fp); err == nil {
			path = fp
			break
		}
	}

	if path == "" {
		return "", fmt.Errorf("config not found")
	}

	f, err := os.Open(path)
	if err != nil {
		return path, err
	}
	defer f.Close()

	if err := toml.NewDecoder(f).Decode(&config); err != nil {
		return path, err
	}

	return path, nil
}
