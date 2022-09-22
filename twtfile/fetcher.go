package twtfile

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/makeworld-the-better-one/go-gemini"
)

type Fetcher struct {
	version          string
	discloseIdentity bool
	client           *http.Client
}

func NewFetcher(version string, discloseIdentity bool, timeout time.Duration) *Fetcher {
	return &Fetcher{version: version,
		discloseIdentity: discloseIdentity,
		client:           getHTTPClient(timeout),
	}
}

func (f *Fetcher) Fetch(nick, url string) (Tweets, TweetMetadata, error) {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return fetchHttp(nick, url, f.version, f.discloseIdentity)
	}

	if strings.HasPrefix(url, "gemini://") {
		return fetchGemini(nick, url)
	}

	return Tweets{}, TweetMetadata{}, nil
}

func fetchHttp(nick, url, version string, discloseIdentity bool) (Tweets, TweetMetadata, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Tweets{}, TweetMetadata{}, fmt.Errorf("could not create request to %s, %w", url, err)
	}

	if discloseIdentity {
		req.Header.Set("User-Agent", fmt.Sprintf("twx/{%s (+%s; @%s)", version, url, nick))
	} else {
		req.Header.Set("User-Agent", fmt.Sprintf("twx/%s", version))
	}

	res, err := client.Do(req)
	if err != nil {
		return Tweets{}, TweetMetadata{}, fmt.Errorf("could not request %s, %w", url, err)
	}

	if res.StatusCode != 200 {
		return Tweets{}, TweetMetadata{}, fmt.Errorf("not OK request, received %v", res.StatusCode)
	}
	body := res.Body
	defer body.Close()

	b, err := io.ReadAll(body)
	if err != nil {
		return Tweets{}, TweetMetadata{}, fmt.Errorf("cannot Read body received %s, %w", url, err)
	}

	ts, meta := parseBody(nick, url, b)
	return ts, meta, nil
}

func fetchGemini(nick, url string) (Tweets, TweetMetadata, error) {
	res, err := gemini.Fetch(url)
	if err != nil {
		return nil, TweetMetadata{}, fmt.Errorf("could not request %s, %w", url, err)
	}
	defer res.Body.Close()

	if res.Status > 200 {
		return nil, TweetMetadata{}, fmt.Errorf("gemini status %d failed", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Tweets{}, TweetMetadata{}, fmt.Errorf("cannot Read body received %s, %w", url, err)
	}

	ts, meta := parseBody(nick, url, body)
	return ts, meta, nil
}

func fetchFile(nick, path, url string) (Tweets, TweetMetadata, error) {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, TweetMetadata{}, fmt.Errorf("can't open file %s", path)
	}

	ts, meta := parseBody(nick, url, body)
	return ts, meta, nil

}

func parseBody(nick string, url string, b []byte) (Tweets, TweetMetadata) {
	ts := Tweets{}
	meta := TweetMetadata{
		Follow: map[string]string{},
		Link:   map[string]string{},
	}
	scanner := bufio.NewScanner(bytes.NewReader(b))
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		if line[0] == '#' {
			meta = parseMeta(line, meta)
		}
		t := parseTweet(nick, url, line)
		if t != nil {
			ts = append(ts, *t)
		}
	}
	return ts, meta
}

var client *http.Client
var once sync.Once

func getHTTPClient(timeout time.Duration) *http.Client {
	if client == nil {
		once.Do(func() {
			netTransport := &http.Transport{
				Dial: (&net.Dialer{
					Timeout:   time.Second,
					KeepAlive: 0,
				}).Dial,
				TLSHandshakeTimeout: 5 * time.Second,
				IdleConnTimeout:     0,
				MaxIdleConnsPerHost: 50000,
				MaxIdleConns:        50000,
			}

			client = &http.Client{
				Timeout:   timeout,
				Transport: netTransport,
			}
		})
	}

	return client

}
