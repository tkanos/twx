package command

import (
	"os/exec"

	"github.com/tkanos/twx/config"
)

func Execute(action string, command string, conf map[string]string) (configToSave map[string]string, parameters map[string]string, err error) {

	if action != "tweet" {
		return nil, nil, nil
	}

	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Dir = config.HomeDirectory()

	err = cmd.Run()
	if err != nil {
		return nil, nil, err
	}

	return nil, nil, nil
}
