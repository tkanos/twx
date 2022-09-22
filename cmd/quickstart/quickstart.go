package quickstart

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tkanos/twx/twtfile"
	"github.com/tkanos/twx/utils"
	"gopkg.in/ini.v1"
)

var q quickstart

// quickstartCmd represents the quickstart command
var quickstartCmd = &cobra.Command{
	Use:   "quickstart",
	Short: "Quickstart wizard for setting up twx.",
	Long:  `Quickstart wizard for setting up twx.`,
	Run: func(cmd *cobra.Command, args []string) {
		configFilePath, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Fatal(err)
		}
		q.configFilePath = utils.ExpandTilde(configFilePath)
		if err := q.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

func Init(rootCmd *cobra.Command) {
	q = quickstart{reader: os.Stdin}

	var username string
	if nick, err := user.Current(); err == nil {
		username = nick.Username
	}

	quickstartCmd.Flags().BoolVarP(&(q.discloseIdentity), "disclose-identity", "", true, "Show your nickname and url in the User Agent.")
	quickstartCmd.Flags().StringVarP(&(q.twtxtFilePath), "file", "f", defaultTwtxtFileLocation(), "Show your nickname and url in the User Agent.")
	quickstartCmd.Flags().BoolVarP(&(q.news), "follow-news", "", true, "Follow the official twtxt and twtr news feeds.")
	quickstartCmd.Flags().StringVarP(&(q.nick), "nick", "n", username, "Specify the nickname for your feed.")
	quickstartCmd.Flags().StringVarP(&(q.url), "url", "u", "", "Specify the url that your feed will be hosted at.")

	rootCmd.AddCommand(quickstartCmd)
}

type quickstart struct {
	configFilePath   string
	nick             string
	url              string
	twtxtFilePath    string
	follow           map[string]string
	news             bool
	discloseIdentity bool
	reader           io.Reader
}

func (q *quickstart) Run() error {
	// Launch wizard
	q.wizard()

	//Create path if it doesn't exist
	if err := os.Mkdir(path.Dir(q.configFilePath), 0755); err != nil && !os.IsExist(err) {
		return err
	}

	//Create INI file
	if err := q.createIniConfig(); err != nil {
		return fmt.Errorf("could no create config file [%s], %w", q.configFilePath, err)
	}

	if err := q.createTwtxtFile(); err != nil {
		return fmt.Errorf("could no create twtxt.txt file [%s], %w", q.twtxtFilePath, err)
	}

	return nil
}

func (q *quickstart) wizard() {
	q.setNick()
	q.setConfigFile()
	q.setTwtxtFile()
	q.setUrl()
	q.setDiscloseIdentity()
	q.setFollowNews()
}

func (q *quickstart) setNick() {
	if v := scanString(fmt.Sprintf("Please enter your desired nick [%s]: ", q.nick), q.reader); v != "" {
		q.nick = v
	}
}

func (q *quickstart) setConfigFile() {
	if v := scanString(fmt.Sprintf("Please enter your desired location for you config file [%s]: ", q.configFilePath), q.reader); v != "" {
		q.configFilePath = v
	}
}

func (q *quickstart) setTwtxtFile() {
	q.twtxtFilePath = utils.ExpandTilde(q.twtxtFilePath)
	if v := scanString(fmt.Sprintf("Please enter the desired location for your twtxt file [%s]: ", q.twtxtFilePath), q.reader); v != "" {
		q.twtxtFilePath = utils.ExpandTilde(v)
	}
}

func (q *quickstart) setUrl() {
	if v := scanString("Please enter the URL your twtxt file will be accessible from (example : https://twtxt.net/user/<NICK>/twtxt.txt): ", q.reader); v != "" {
		q.url = v
	}
}

func (q *quickstart) setDiscloseIdentity() {
	if v := scanString("Do you want to disclose your identity? Your nick and URL will be shared when making HTTP requests (by the UserAgent) [Y/n]: ", q.reader); v != "" && strings.ToLower(v) == "n" {
		q.discloseIdentity = false
	}
}

func (q *quickstart) setFollowNews() {
	if v := scanString("Do you want to follow the twtxt news feeds? [Y/n]: ", q.reader); v != "" && strings.ToLower(v) == "n" {
		q.news = false
	}

	if q.news {
		q.follow = map[string]string{
			"twtxt_news": "https://buckket.org/twtxt_news.txt",
		}
	}
}

func (q *quickstart) createIniConfig() error {
	config := ini.Empty()
	config.Section("twtxt").Key("nick").SetValue(q.nick)
	config.Section("twtxt").Key("twtfile").SetValue(q.twtxtFilePath)
	config.Section("twtxt").Key("twturl").SetValue(q.url)
	config.Section("twtxt").Key("disclose_identity").SetValue(fmt.Sprint(q.discloseIdentity))

	// fillout following list
	if q.news {
		for nick, url := range q.follow {
			config.Section("following").Key(nick).SetValue(url)
		}
	}

	// write the configuration to the selected config file
	if err := config.SaveTo(q.configFilePath); err != nil {
		return err
	}

	return nil
}

func (q *quickstart) createTwtxtFile() error {
	// Get Template
	b, err := twtfile.TweetMetadata{
		Nick:      q.nick,
		URL:       q.url,
		Follow:    q.follow,
		Following: len(q.follow),
	}.CreateTwtxtMetaTemplate()

	if err != nil {
		return err
	}

	// Create Folder if it does not exist
	if err := os.Mkdir(path.Dir(q.twtxtFilePath), 0755); err != nil && !os.IsExist(err) {
		return err
	}

	//Create twtxt.txt file
	f, err := os.Create(q.twtxtFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	//Write in txt.txt file
	_, err = f.Write(b)
	if err != nil {
		return err
	}

	return nil

}

func defaultTwtxtFileLocation() string {
	home, err := os.UserHomeDir()
	if err != nil {
		//TODO : check what to do with different os
	}
	return filepath.Join(home, "twtxt.txt")
}

func scanString(text string, reader io.Reader) string {
	var value string
	fmt.Print(text)
	fmt.Fscanln(reader, &value)

	return strings.Trim(value, " ")
}
