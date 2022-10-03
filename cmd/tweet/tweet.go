package tweet

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/tkanos/twx/cmd/context"
	"github.com/tkanos/twx/cmd/hooks"
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

		argstweet = strings.Join(args, "")

		if context.Config.Twtxt.PreTweetHook != "" {
			etype := "cmd"
			switch context.Config.Twtxt.PreTweetHook {
			case "{{yarn}}":
				etype = "yarn"
			}

			conf, output, err := hooks.Execute(etype, "tweet", context.Config.Twtxt.PreTweetHook, map[string]string{"tweet": argstweet, "reply": replyHash}, context.Config.Hook)
			if err != nil {
				log.Fatalf("Could not execute PostHook: %s", err)
			}
			if c, ok := output["created"]; ok {
				created = c
			}

			if conf != nil {
				context.Config.Hook = conf
				err = context.Config.Save()
				if err != nil {
					fmt.Printf("the tweet pre command executed successfully but the config could not be saved: %s\n", err)
				}
			}

		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("Wrong arguments")
		}

		if err := t.Run(argstweet); err != nil {
			log.Fatal(err)
		}

	},
	PostRun: func(cmd *cobra.Command, args []string) {
		if context.Config.Twtxt.PostTweetHook != "" {
			etype := "cmd"
			conf, _, err := hooks.Execute(etype, "tweet", context.Config.Twtxt.PostTweetHook, map[string]string{"tweet": argstweet, "reply": replyHash}, context.Config.Hook)
			if err != nil {
				log.Fatalf("Could not execute PostHook: %s", err)
			}

			if conf != nil {
				context.Config.Hook = conf
				err = context.Config.Save()
				if err != nil {
					log.Fatalf("the tweet post command executed successfully but the config could not be saved: %s", err)
				}
			}
		}
	},
}

var replyHash string
var created string

func Init(rootCmd *cobra.Command) {

	t = tweet{}

	tweetCmd.Flags().StringVarP(&replyHash, "reply", "r", "", "reply to a hash")

	rootCmd.AddCommand(tweetCmd)
}

var argstweet string

type tweet struct {
}

func (t *tweet) Run(text string) error {

	cTime := time.Now()
	if created != "" {
		cTime, _ = time.Parse(time.RFC3339, created)
	}

	text = addMentions(text)

	//append to file or create file
	if tweet, err := context.TwtFile.Tweet(context.Config.Twtxt.Nick, context.Config.Twtxt.TwtURL, text, replyHash, cTime); err != nil {
		return err
	} else {
		argstweet = tweet.String()
		fmt.Println(argstweet)
	}

	return nil
}

var re = regexp.MustCompile(`@[^< ]+`)

func addMentions(twt string) string {
	if !strings.Contains(twt, "@") {
		return twt
	}

	groups := re.FindAllString(twt, -1)
	if groups == nil {
		return twt
	}

	replaced := 0
	for k, v := range context.Config.Following {
		for _, nick := range groups {
			if strings.EqualFold("@"+k, nick) {
				replaced++
				twt = strings.Replace(twt, nick, fmt.Sprintf("@<%s %s>", k, v), 1)
			}
		}
		if replaced == len(groups) {
			break
		}
	}

	return twt

}
