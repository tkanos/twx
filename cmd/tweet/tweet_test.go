package tweet

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tkanos/twx/cmd/context"
	"github.com/tkanos/twx/config"
)

func TestAddMentions(t *testing.T) {

	// Arrange
	test_data := []struct {
		input    string
		expected string
		message  string
	}{
		{
			input:    "Hello everyone",
			expected: "Hello everyone",
			message:  "1. when no mentions on twt, we should get it back",
		},
		{
			input:    "Hello @Nick",
			expected: "Hello @<Nick http://nick.com>",
			message:  "2. A mention should be replaced by it's twtxt counterpart",
		},
		{
			input:    "Hello @nick",
			expected: "Hello @<Nick http://nick.com>",
			message:  "3. A mention (even in different case) should be replaced by it's twtxt counterpart",
		},
		{
			input:    "Hello @Nick and @Rose",
			expected: "Hello @<Nick http://nick.com> and @<Rose http://rose.com>",
			message:  "4. Many mentions should be replaced",
		},
		{
			input:    "Hello @<Nick nick.com> and @Rose",
			expected: "Hello @<Nick nick.com> and @<Rose http://rose.com>",
			message:  "5. Mentions already in the good format shouldn't be replaced",
		},
		{
			input:    "Hello @<Nick nick.com>",
			expected: "Hello @<Nick nick.com>",
			message:  "6. Mentions already in the good format shouldn't be replaced",
		},

		{
			input:    "Hello @Bob",
			expected: "Hello @Bob",
			message:  "7. Mentions not followed cannot be replaced",
		},
	}

	context.Config = &config.Configuration{
		Following: map[string]string{
			"Nick": "http://nick.com",
			"Rose": "http://rose.com",
		},
	}

	for _, tt := range test_data {
		// Act
		actual := addMentions(tt.input)

		// Assert
		assert.Equal(t, tt.expected, actual, tt.message)
	}

	context.Config = nil

}
