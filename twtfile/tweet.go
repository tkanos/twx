package twtfile

import (
	"encoding/base32"
	"fmt"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/blake2b"
)

type Tweet struct {
	Nick     string
	URL      string
	Created  time.Time
	Hash     string
	Text     string
	tweeting bool
}

func (t Tweet) String() string {
	hash := t.Hash
	if !t.tweeting {
		if hash == "" {
			hash = t.GenerateHash(t.URL, t.Created, t.Text)
		}
		return fmt.Sprintf("%s\t%s\t%s %s\n", t.Nick, t.Created.UTC().Format(time.RFC3339), hash, t.Text)
	}
	if hash == "" {
		return fmt.Sprintf("%s\t%s\n", t.Created.UTC().Format(time.RFC3339), t.Text)
	} else {
		return fmt.Sprintf("%s\t%s %s\n", t.Created.UTC().Format(time.RFC3339), hash, t.Text)
	}
}

func (t Tweet) GenerateHash(url string, created time.Time, text string) string {
	payload := url + "\n" + created.Format(time.RFC3339) + "\n" + text
	sum := blake2b.Sum256([]byte(payload))
	encoding := base32.StdEncoding.WithPadding(base32.NoPadding)
	hash := strings.ToLower(encoding.EncodeToString(sum[:]))
	hash = hash[len(hash)-7:]

	return fmt.Sprintf("(#%s)", hash)
}

var retweet = regexp.MustCompile(`(^(.+)\s+(\(\#.+\))\s+(.+)$)|(^(.+?)\s+(.+)$)`)

func parseTweet(nick, url, line string) *Tweet {
	if len(line) == 0 {
		return nil
	}

	var t *Tweet

	groups := retweet.FindStringSubmatch(line)
	if groups == nil || groups[0] == "" {
		return nil
	}

	if groups[1] != "" {

		created, err := parseTime(groups[2])
		if err != nil {
			return nil
		}

		t = &Tweet{
			Nick:    nick,
			URL:     url,
			Created: created,
			Hash:    groups[3],
			Text:    groups[4],
		}

	} else if groups[5] != "" {
		created, err := parseTime(groups[6])
		if err != nil {
			return nil
		}
		t = &Tweet{
			Nick:    nick,
			URL:     url,
			Created: created,
			Text:    groups[7],
		}
	}

	return t
}

func parseTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}
