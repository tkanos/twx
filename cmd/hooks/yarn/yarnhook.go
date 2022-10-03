package yarn

import (
	"bufio"
	"fmt"
	"os"
	"syscall"
	"time"

	"go.yarn.social/client"
	"golang.org/x/term"
)

type config struct {
	timeout time.Duration
	url     string
	token   string
}

func validateConfig(conf map[string]string) (c config, err error) {
	c = config{
		timeout: 5 * time.Second,
	}

	if v, ok := conf["yarn_url"]; ok {
		c.url = v
	} else {
		return c, fmt.Errorf("yarn_url missing in hook config")
	}

	if v, ok := conf["yarn_token"]; ok {
		c.token = v
	}

	return c, nil

}

func Execute(action string, parameter map[string]string, conf map[string]string) (configToSave map[string]string, parameters map[string]string, err error) {

	var c config
	if c, err = validateConfig(conf); err != nil {
		return nil, nil, err
	}

	fmt.Printf("yarn social pre hook plugin for %s\n", c.url)

	var save bool
	if c.token == "" {
		cli, err := client.NewClient(client.WithURI(c.url))
		if err != nil {
			return nil, nil, fmt.Errorf("error creating client, %s", err)
		}
		c.token, err = login(cli)
		if err != nil {
			return nil, nil, fmt.Errorf("could not connect to yarn social %s, %s", c.url, err)
		}
		conf["yarn_token"] = c.token
		save = true
	}

	cli, err := client.NewClient(
		client.WithURI(c.url),
		client.WithToken(c.token),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create a client to yarn social %s, %s", c.url, err)
	}

	output := map[string]string{}

	switch action {
	case "tweet":
		var created, hash, twt, reply string
		twt = parameter["tweet"]
		reply = parameter["reply"]
		created, hash, err = tweet(cli, twt, reply)
		output["created"] = created
		output["hash"] = hash
	case "follow":
		err = follow(cli, parameter["nick"], parameter["url"])
	case "unfollow":
		err = unfollow(cli, parameter["nick"])
	}
	if err != nil {
		return nil, nil, err
	}

	if save {
		return conf, output, nil
	}

	return nil, output, nil
}

func login(cli *client.Client) (token string, err error) {
	// get username and password
	username, password, err := readCredentials()
	if err != nil {
		return "", err
	}

	res, err := cli.Login(username, password)
	if err != nil {
		return "", err
	}

	fmt.Println("login successful")

	return res.Token, nil
}

func readCredentials() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}

	fmt.Print("Password: ")
	data, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}
	password := string(data)

	return username, password, nil
}

func tweet(cli *client.Client, tweet, reply string) (created, hash string, err error) {
	if reply != "" {
		tweet = fmt.Sprintf("(#%s) %s", reply, tweet)
	}

	p, err := cli.Post(tweet, "")
	if err != nil {
		return "", reply, err
	}
	return p.Created, p.Hash, err
}

func follow(cli *client.Client, nick, url string) error {
	return cli.Follow(nick, url)
}

func unfollow(cli *client.Client, nick string) error {
	return cli.Unfollow(nick)
}
