package timeline

import (
	"fmt"
	"log"
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

		t.f = twtfile.NewFetcher(context.Version, context.Config.Twtxt.DiscloseIdentity, context.Config.TimeoutDuration())

		if err := t.Run(strings.Join(args, "")); err != nil {
			log.Fatal(err)
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

	if context.Config.Twtxt.IncludeYourself {
		// Adding my own tweets
		tweets = append(tweets, context.TwtFile.Tweets...)
	}

	// Sorting Tweets
	tweets = t.Sort(tweets, context.Config.Twtxt.Sorting, context.Config.Twtxt.LimitTimeline, reverse)

	for _, v := range tweets {
		fmt.Println(v)
	}

	return nil
}

func (t *timeline) Sort(tweets twtfile.Tweets, sorting string, limit int, reverse bool) twtfile.Tweets {
	sort.Sort(sort.Reverse(tweets))

	if limit > 0 && len(tweets) > limit {
		tweets = tweets[:limit]
	}

	if reverse {
		if strings.ToLower(sorting) == "ascending" {
			sort.Sort(tweets)
		}
	} else {
		if strings.ToLower(sorting) != "ascending" {
			sort.Sort(tweets)
		}
	}

	return tweets
}
