package follow

import (
	"fmt"
	"net"
	"net/http"
	"os"
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
			fmt.Println("Wrong Arguments passed")
			os.Exit(1)
		}

		if err := f.Run(args[0], args[1]); err != nil {
			panic(err)
		}
	},
}

var replace bool

func Init(rootCmd *cobra.Command) {

	f = follow{
		client: getHTTPClient(10 * time.Second),
	}

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

	//Trying to add following user to Config file
	err := f.addFollowingInConfig(nick, url)
	if err != nil {
		return err
	}

	if context.Config.Twtxt.DiscloseIdentity {
		//Trying to add following user to twtxt.txt file
		err = f.addFollowingInMetadataOfTwtxtFile(nick, url)
		if err != nil {
			return err
		}
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

func (f *follow) addFollowingInConfig(nick, url string) error {
	// check for duplicate and whether duplicate is allowed
	if _, ok := context.Config.Following[nick]; ok && !replace {
		return nil // fmt.Errorf("already following @%s %s", nick, url)
	}

	// update the following section
	context.Config.Following[nick] = url

	// write the configuration to the selected config file
	if err := context.Config.Save(); err != nil {
		return err
	}

	return nil
}

func (f *follow) addFollowingInMetadataOfTwtxtFile(nick, url string) error {
	if _, ok := context.TwtFile.Meta.Follow[nick]; !ok || replace {
		context.TwtFile.Meta.Follow[nick] = url
		if !replace {
			context.TwtFile.Meta.Following = context.TwtFile.Meta.Following + 1
		}
	}

	//Save
	context.TwtFile.SaveTwtxtFileWithMetadata(context.Config.Twtxt.DiscloseIdentity)
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
