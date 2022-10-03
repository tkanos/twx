package unfollow

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/tkanos/twx/cmd/context"
	"github.com/tkanos/twx/cmd/hooks"
)

var f unfollow

// followCmd represents the follow command
var unfollowCmd = &cobra.Command{
	Use:   "unfollow <NICK>",
	Short: "Unfollow another user",
	Long: `Unfollow another user
	
	Example: unfollow Nick`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("Wrong arguments")
		}

		if context.Config.Twtxt.PreTweetHook != "" {
			var etype string
			switch context.Config.Twtxt.PreTweetHook {
			case "{{yarn}}":
				etype = "yarn"
			default:
				return
			}

			conf, _, err := hooks.Execute(etype, "unfollow", context.Config.Twtxt.PreTweetHook, map[string]string{"nick": args[0]}, context.Config.Hook)
			if err != nil {
				log.Fatalf("Could not execute PreHook: %s", err)
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

		if err := f.Run(args[0]); err != nil {
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
		context.TwtFile.Save(context.Config.Twtxt.DiscloseIdentity)

	},
}

func Init(rootCmd *cobra.Command) {

	f = unfollow{}

	rootCmd.AddCommand(unfollowCmd)
}

type unfollow struct {
}

func (f *unfollow) Run(nick string) error {

	// Remove Followed user in Config
	context.Config.Unfollow(nick)

	//Remove Followed user in Twtxt.txt
	context.TwtFile.Unfollow(nick)

	return nil
}
