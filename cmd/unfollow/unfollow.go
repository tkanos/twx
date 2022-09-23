package unfollow

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/tkanos/twx/cmd/context"
)

var f unfollow

// followCmd represents the follow command
var unfollowCmd = &cobra.Command{
	Use:   "unfollow <NICK>",
	Short: "Unfollow another user",
	Long: `Unfollow another user
	
	Example: unfollow Nick`,
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
		context.TwtFile.SaveTwtxtFileWithMetadata(context.Config.Twtxt.DiscloseIdentity)

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
