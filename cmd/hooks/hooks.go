package hooks

import (
	"fmt"

	"github.com/tkanos/twx/cmd/hooks/command"
	"github.com/tkanos/twx/cmd/hooks/yarn"
)

func Execute(etype string, action string, cmd string, parameter map[string]string, config map[string]string) (configToSave map[string]string, parameters map[string]string, err error) {

	switch etype {
	case "cmd":
		return command.Execute(action, cmd, config)
	case "yarn":
		return yarn.Execute(action, parameter, config)
	}

	return nil, nil, fmt.Errorf("No action executed")
}
