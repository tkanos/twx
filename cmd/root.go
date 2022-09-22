package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tkanos/twx/cmd/context"
	"github.com/tkanos/twx/cmd/follow"
	"github.com/tkanos/twx/cmd/quickstart"
	"github.com/tkanos/twx/cmd/timeline"
	"github.com/tkanos/twx/cmd/tweet"
	"github.com/tkanos/twx/cmd/unfollow"
	"github.com/tkanos/twx/config"
	"github.com/tkanos/twx/twtfile"
	"github.com/tkanos/twx/utils"
)

var configfile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "twx",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		configFilePath, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Fatalf("Cannot Get config file: %s", err)
		}

		configFilePath = utils.ExpandTilde(configFilePath)
		if context.Config, err = config.NewConfig(configFilePath); err != nil {
			log.Fatal(err)
		}

		if twtFile := context.Config.Twtxt.TwtFile; twtFile == "" {
			log.Fatalf("Cannot Get config file: %s", err)
		} else {
			twtFile = utils.ExpandTilde(twtFile)
			if _, err := os.Stat(twtFile); err != nil {
				//if it does not exist we should create it
				file, err := os.Create(twtFile)
				if err != nil {
					log.Fatalf("Cannot open twtxt file: %s", err)
				}
				defer file.Close()
			}
			context.TwtFile, err = twtfile.NewTwtFile(context.Config.Twtxt.Nick, twtFile, context.Config.Twtxt.TwtURL)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	configFolder, err := os.UserConfigDir()
	if err != nil {
		configFolder, err = os.UserHomeDir()
		if err != nil {
			//TODO : don't forget to manage apple/linux/windows/.....
		}
	}
	configDefaultFile := filepath.Join(configFolder, "twtxt", "config")
	rootCmd.PersistentFlags().StringVarP(&configfile, "config", "c", configDefaultFile, fmt.Sprintf("config file (default is %s)", configDefaultFile))
	quickstart.Init(rootCmd)
	follow.Init(rootCmd)
	unfollow.Init(rootCmd)
	tweet.Init(rootCmd)
	timeline.Init(rootCmd)
}
