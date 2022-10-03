package twtfile

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

type TwtFile struct {
	Path   string
	Tweets Tweets
	Meta   TweetMetadata
}

type Tweets []Tweet

func (tweets Tweets) Len() int {
	return len(tweets)
}
func (tweets Tweets) Less(i, j int) bool {
	return tweets[i].Created.Before(tweets[j].Created)
}
func (tweets Tweets) Swap(i, j int) {
	tweets[i], tweets[j] = tweets[j], tweets[i]
}

func NewTwtFile(nick, path, url string) (*TwtFile, error) {
	t := &TwtFile{Path: path}

	var err error
	t.Tweets, t.Meta, err = fetchFile(nick, path, url)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t *TwtFile) Follow(nick, url string, replace bool) {
	if _, ok := t.Meta.Follow[nick]; !ok || replace {
		t.Meta.Follow[nick] = url
		if !replace {
			t.Meta.Following = t.Meta.Following + 1
		}
	}
}

func (t *TwtFile) Unfollow(nick string) {
	if _, ok := t.Meta.Follow[nick]; ok {
		t.Meta.Following = t.Meta.Following - 1
	}

	delete(t.Meta.Follow, nick)
}

func (t *TwtFile) Tweet(nick, url, text, replyHash string, created time.Time) (Tweet, error) {
	tweet := Tweet{
		Nick:     nick,
		URL:      url,
		Created:  created,
		Text:     text,
		tweeting: true,
	}

	if replyHash != "" {
		tweet.Hash = getReplyHash(replyHash)
	}

	return tweet, t.appendTweet(tweet.String())
}

func (t *TwtFile) appendTweet(payload string) error {
	if _, err := os.Stat(t.Path); errors.Is(err, os.ErrNotExist) {
		// file does not exist
	}

	f, err := os.OpenFile(t.Path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.WriteString(payload); err != nil {
		return err
	}

	return nil
}

func getReplyHash(replyHash string) string {
	if strings.HasPrefix(replyHash, "#") {
		return fmt.Sprintf("(%s)", replyHash)
	}
	return fmt.Sprintf("(#%s)", replyHash)
}

func (t *TwtFile) Save(discloseIdentity bool) error {
	// Read twtxt.txt File
	file, err := os.Open(t.Path)
	if err != nil {
		return fmt.Errorf("could not open file %s, %w", t.Path, err)
	}
	defer file.Close()

	// Create new twtxt.txt file
	newfilePath := path.Join(path.Dir(t.Path), "newtwtxt.txt")
	outfile, err := os.Create(newfilePath)
	if err != nil {
		return fmt.Errorf("could not create write file  %s, %w", newfilePath, err)
	}
	defer outfile.Close()

	var metaWriteError error

	// append at the start
	if discloseIdentity {
		tpl, metaWriteError := t.Meta.CreateTwtxtMetaTemplate()
		if metaWriteError == nil {
			_, err = outfile.WriteString(string(tpl))
			if err != nil {
				return fmt.Errorf("could not write to file  %s, %w", newfilePath, err)
			}
		}
	}

	scanner := bufio.NewScanner(file)
	// read the file to be appended to and output all of it
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " ")
		if len(line) == 0 || (len(line) > 0 && (line[0] != '#' || metaWriteError != nil)) {
			outfile.WriteString(line)
			outfile.WriteString("\n")
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("could not read twtxt.txt file  %s, %w", newfilePath, err)
	}
	// ensure all lines are written
	outfile.Sync()
	// over write the old file with the new one
	err = os.Rename(newfilePath, t.Path)
	if err != nil {
		return fmt.Errorf("could not rename %s to %s, %w", newfilePath, t.Path, err)
	}

	return nil
}
