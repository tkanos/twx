package timeline

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tkanos/twx/twtfile"
)

func TestSort(t *testing.T) {

	// Arrange
	test_data := []struct {
		tweets  twtfile.Tweets
		limit   int
		sorting string
		reverse bool
		check   func(twtfile.Tweets) bool
		message string
	}{
		{
			message: "1. Sort Descending",
			limit:   -1,
			sorting: "descending",
			tweets: twtfile.Tweets{
				twtfile.Tweet{Created: time.Now().Add(-48 * time.Hour)},
				twtfile.Tweet{Created: time.Now()},
				twtfile.Tweet{Created: time.Now().Add(-72 * time.Hour)},
				twtfile.Tweet{Created: time.Now().Add(-48 * time.Hour)},
			},
			check: func(tweets twtfile.Tweets) bool {
				return tweets[3].Created.After(tweets[0].Created)
			},
		},
		{
			message: "2. Sort Ascending",
			limit:   -1,
			sorting: "ascending",
			tweets: twtfile.Tweets{
				twtfile.Tweet{Created: time.Now().Add(-48 * time.Hour)},
				twtfile.Tweet{Created: time.Now()},
				twtfile.Tweet{Created: time.Now().Add(-72 * time.Hour)},
				twtfile.Tweet{Created: time.Now().Add(-48 * time.Hour)},
			},
			check: func(tweets twtfile.Tweets) bool {
				return tweets[0].Created.After(tweets[3].Created)
			},
		},
		{
			message: "3. Default Sorting should be Descending",
			limit:   -1,
			sorting: "",
			tweets: twtfile.Tweets{
				twtfile.Tweet{Created: time.Now().Add(-48 * time.Hour)},
				twtfile.Tweet{Created: time.Now()},
				twtfile.Tweet{Created: time.Now().Add(-72 * time.Hour)},
				twtfile.Tweet{Created: time.Now().Add(-48 * time.Hour)},
			},
			check: func(tweets twtfile.Tweets) bool {
				return tweets[3].Created.After(tweets[0].Created)
			},
		},
		{
			message: "4. Reverse Descending",
			limit:   -1,
			sorting: "descending",
			reverse: true,
			tweets: twtfile.Tweets{
				twtfile.Tweet{Created: time.Now().Add(-48 * time.Hour)},
				twtfile.Tweet{Created: time.Now()},
				twtfile.Tweet{Created: time.Now().Add(-72 * time.Hour)},
				twtfile.Tweet{Created: time.Now().Add(-48 * time.Hour)},
			},
			check: func(tweets twtfile.Tweets) bool {
				return tweets[0].Created.After(tweets[3].Created)
			},
		},
		{
			message: "5. Reverse Ascending",
			limit:   -1,
			sorting: "ascending",
			reverse: true,
			tweets: twtfile.Tweets{
				twtfile.Tweet{Created: time.Now().Add(-48 * time.Hour)},
				twtfile.Tweet{Created: time.Now()},
				twtfile.Tweet{Created: time.Now().Add(-72 * time.Hour)},
				twtfile.Tweet{Created: time.Now().Add(-48 * time.Hour)},
			},
			check: func(tweets twtfile.Tweets) bool {
				return tweets[3].Created.After(tweets[0].Created)
			},
		},
		{
			message: "6. Reverse Default Sorting should be Ascending",
			limit:   -1,
			sorting: "",
			reverse: true,
			tweets: twtfile.Tweets{
				twtfile.Tweet{Created: time.Now().Add(-48 * time.Hour)},
				twtfile.Tweet{Created: time.Now()},
				twtfile.Tweet{Created: time.Now().Add(-72 * time.Hour)},
				twtfile.Tweet{Created: time.Now().Add(-48 * time.Hour)},
			},
			check: func(tweets twtfile.Tweets) bool {
				return tweets[0].Created.After(tweets[3].Created)
			},
		},
		{
			message: "7. Limit Timeline when limit < tweets",
			limit:   2,
			tweets: twtfile.Tweets{
				twtfile.Tweet{Created: time.Now().Add(-48 * time.Hour)},
				twtfile.Tweet{Created: time.Now()},
				twtfile.Tweet{Created: time.Now().Add(-72 * time.Hour)},
				twtfile.Tweet{Created: time.Now().Add(-48 * time.Hour)},
			},
			check: func(tweets twtfile.Tweets) bool {
				return len(tweets) == 2
			},
		},
		{
			message: "8. Limit Timeline when limit > tweets",
			limit:   20,
			tweets: twtfile.Tweets{
				twtfile.Tweet{Created: time.Now().Add(-48 * time.Hour)},
				twtfile.Tweet{Created: time.Now()},
				twtfile.Tweet{Created: time.Now().Add(-72 * time.Hour)},
				twtfile.Tweet{Created: time.Now().Add(-48 * time.Hour)},
			},
			check: func(tweets twtfile.Tweets) bool {
				return len(tweets) == 4
			},
		},
		{
			message: "9. Limit Timeline when limit= 0",
			limit:   0,
			tweets: twtfile.Tweets{
				twtfile.Tweet{Created: time.Now().Add(-48 * time.Hour)},
				twtfile.Tweet{Created: time.Now()},
				twtfile.Tweet{Created: time.Now().Add(-72 * time.Hour)},
				twtfile.Tweet{Created: time.Now().Add(-48 * time.Hour)},
			},
			check: func(tweets twtfile.Tweets) bool {
				return len(tweets) == 4
			},
		},
	}

	for _, tt := range test_data {
		tl := timeline{}

		// Act
		actual := tl.Sort(tt.tweets, tt.sorting, tt.limit, tt.reverse)

		// Assert
		assert.True(t, tt.check(actual), tt.message)
	}
}
