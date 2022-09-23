package tweet

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tkanos/twx/cmd/context"
	"github.com/tkanos/twx/config"
)

var t tweet

// followCmd represents the follow command
var tweetCmd = &cobra.Command{
	Use:     "tweet [-r hash] <tweet>",
	Aliases: []string{"post", "twt"},
	Short:   "Tweet or reply a  text",
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
			log.Fatal("Wrong arguments")
		}

		if err := t.Run(strings.Join(args, "")); err != nil {
			log.Fatal(err)
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

	//append to file or create file
	if tweet, err := context.TwtFile.Tweet(context.Config.Twtxt.Nick, context.Config.Twtxt.TwtURL, text, replyHash); err != nil {
		return err
	} else {
		fmt.Println(tweet.String())
	}

	return nil
}
