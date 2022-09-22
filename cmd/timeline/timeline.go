package timeline

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tkanos/twx/cmd/context"
	"github.com/tkanos/twx/twtfile"
)

var t timeline

// followCmd represents the follow command
var timelineCmd = &cobra.Command{
	Use:   "timeline",
	Short: "Form a timeline",
	Long: `Form a Timeline
	
	Example: 
	timeline`,
	Run: func(cmd *cobra.Command, args []string) {

		t.f = twtfile.NewFetcher(context.Version, context.Config.Twtxt.DiscloseIdentity)

		if err := t.Run(strings.Join(args, "")); err != nil {
			panic(err)
		}

	},
}

var reverse bool

func Init(rootCmd *cobra.Command) {

	t = timeline{}

	timelineCmd.Flags().BoolVarP(&reverse, "reverse", "r", false, "reverse the timeline order")

	rootCmd.AddCommand(timelineCmd)
}

type timeline struct {
	f *twtfile.Fetcher
}

func (t *timeline) Run(text string) error {
	tweets := twtfile.Tweets{}
	for n, u := range context.Config.Following {
		f, _, _ := t.f.Fetch(n, u)
		tweets = append(tweets, f...)
	}

	// Adding my own tweets
	tweets = append(tweets, context.TwtFile.Tweets...)

	if reverse {
		if strings.ToLower(context.Config.Twtxt.Sorting) == "descending" {
			sort.Sort(sort.Reverse(tweets))
		} else {
			sort.Sort(tweets)
		}
	} else {
		if strings.ToLower(context.Config.Twtxt.Sorting) == "descending" {
			sort.Sort(tweets)
		} else {
			sort.Sort(sort.Reverse(tweets))
		}
	}

	fmt.Println(tweets)

	return nil
}
