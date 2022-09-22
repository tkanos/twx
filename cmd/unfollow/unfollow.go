package unfollow

import (
	"os"

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
			os.Exit(1)
		}

		if err := f.Run(args[0]); err != nil {
			panic(err)
		}
	},
}

func Init(rootCmd *cobra.Command) {

	f = unfollow{}

	rootCmd.AddCommand(unfollowCmd)
}

type unfollow struct {
}

func (f *unfollow) Run(nick string) error {

	// Trying to Remove following user in Config file
	err := f.removeFollowingInConfig(nick)

	if err == nil && context.Config.Twtxt.DiscloseIdentity {
		// Trying to Remove following user from twtxt.txt file
		err = f.removeFollowingInMetadataOfTwtxtFile(nick)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *unfollow) removeFollowingInConfig(nick string) error {

	// search Nick
	delete(context.Config.Following, nick)

	// write the configuration to the selected config file
	return context.Config.Save()
}

func (f *unfollow) removeFollowingInMetadataOfTwtxtFile(nick string) error {
	if _, ok := context.TwtFile.Meta.Follow[nick]; ok {
		context.TwtFile.Meta.Following = context.TwtFile.Meta.Following - 1
	}

	delete(context.TwtFile.Meta.Follow, nick)

	//Save
	context.TwtFile.SaveTwtxtFileWithMetadata(context.Config.Twtxt.DiscloseIdentity)

	return nil

}
