package follow

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/makeworld-the-better-one/go-gemini"
	"github.com/spf13/cobra"
	"github.com/tkanos/twx/cmd/context"
)

var f follow

// followCmd represents the follow command
var followCmd = &cobra.Command{
	Use:   "follow <NICK> <URL>",
	Short: "Follow another user of an existing twtxt.txt feed",
	Long: `Follow another user of an existing twtxt.txt feed
	
	Example: follow hacker-news https://feeds.twtxt.net/hacker-news/twtxt.txt`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 2 {
			log.Fatal("Wrong Arguments passed")
		}

		f.client = getHTTPClient(context.Config.TimeoutDuration())

		if err := f.Run(args[0], args[1]); err != nil {
			log.Fatal(err)
		}
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		//Save config and Twtxfile
		// write the configuration to the selected config file
		if err := context.Config.Save(); err != nil {
			log.Fatalf("can't save config: %s", err)
		}

		//Write TwtxtFile Metadata
		context.TwtFile.SaveTwtxtFileWithMetadata(context.Config.Twtxt.DiscloseIdentity)

	},
}

var replace bool

func Init(rootCmd *cobra.Command) {

	f = follow{}

	followCmd.Flags().BoolVarP(&replace, "replace", "r", false, "if the nick exist it will be replaced")

	rootCmd.AddCommand(followCmd)
}

type follow struct {
	client *http.Client
}

func (f *follow) Run(nick, url string) error {

	if context.Config.Twtxt.CheckFollowing {
		err := f.validateFeed(nick, url)
		if err != nil {
			return fmt.Errorf("could not validate, %w", err)
		}
	}

	//Follow user in Config
	context.Config.Follow(nick, url, replace)

	if context.Config.Twtxt.DiscloseIdentity {
		//Follow user in twtxt.txt
		context.TwtFile.Follow(nick, url, replace)
	}

	return nil
}

func (f *follow) validateFeed(nick, url string) error {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return f.validateHttp(nick, url)
	}

	if strings.HasPrefix(url, "gemini://") {
		return f.validateGemini(nick, url)
	}

	return nil
}

func (f *follow) validateHttp(nick, url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("could not create request to %s, %w", url, err)
	}

	if context.Config.Twtxt.DiscloseIdentity {
		req.Header.Set("User-Agent", fmt.Sprintf("twx/{%s (+%s; @%s)", context.Version, url, nick))
	} else {
		req.Header.Set("User-Agent", fmt.Sprintf("twx/%s", context.Version))
	}

	res, err := f.client.Do(req)
	if err != nil {
		return fmt.Errorf("could not request %s, %w", url, err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("not OK request, received %v", res.StatusCode)
	}

	return nil
}

func (f *follow) validateGemini(nick, url string) error {
	res, err := gemini.Fetch(url)
	if err != nil {
		return fmt.Errorf("could not request %s, %w", url, err)
	}
	defer res.Body.Close()

	if res.Status > 200 {
		return fmt.Errorf("gemini status %d failed", res.Status)
	}

	return nil

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
