package tweet

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/tkanos/twx/cmd/context"
	"github.com/tkanos/twx/config"
	"github.com/tkanos/twx/twtfile"
)

var t tweet

// followCmd represents the follow command
var tweetCmd = &cobra.Command{
	Use:   "tweet [-r hash] <tweet>",
	Short: "Tweet or reply a  text",
	Long: `tweet or reply a text
	
	Example: 
	tweet Hello world
	tweet -r ab123c Nice to meet you.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if context.Config.Twtxt.PreTweetHook != "" {
			cmd := exec.Command("/bin/sh", "-c", context.Config.Twtxt.PreTweetHook)
			cmd.Dir = config.HomeDirectory()

			err := cmd.Run()
			if err != nil {
				log.Fatalf("Could not execute PreHook: %s", err)
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			os.Exit(1)
		}

		if err := t.Run(strings.Join(args, "")); err != nil {
			panic(err)
		}

	},
	PostRun: func(cmd *cobra.Command, args []string) {
		if context.Config.Twtxt.PostTweetHook != "" {
			cmd := exec.Command("/bin/sh", "-c", context.Config.Twtxt.PostTweetHook)
			cmd.Dir = config.HomeDirectory()

			err := cmd.Run()
			if err != nil {
				log.Fatalf("Could not execute PostHook: %s", err)
			}
		}
	},
}

var replyHash string

func Init(rootCmd *cobra.Command) {

	t = tweet{}

	tweetCmd.Flags().StringVarP(&replyHash, "reply", "r", "", "reply to a hash")

	rootCmd.AddCommand(tweetCmd)
}

type tweet struct {
}

func (t *tweet) Run(text string) error {
	payload := twtfile.Tweet{
		Nick:     context.Config.Twtxt.Nick,
		URL:      context.Config.Twtxt.TwtURL,
		Created:  time.Now(),
		Text:     text,
		Tweeting: true,
	}

	if replyHash != "" {
		payload.Hash = getReplyHash()
	}

	//append to file or create file
	if err := context.TwtFile.Append(payload.String()); err != nil {
		return err
	}

	fmt.Println(payload)
	return nil
}

func getReplyHash() string {
	if strings.HasPrefix(replyHash, "#") {
		return fmt.Sprintf("(%s)", replyHash)
	}
	return fmt.Sprintf("(#%s)", replyHash)
}
