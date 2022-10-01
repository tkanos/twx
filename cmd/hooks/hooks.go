package hooks

import "github.com/tkanos/twx/cmd/hooks/command"

func Execute(etype string, action string, cmd string, config map[string]string) (configToSave map[string]string, err error) {

	switch etype {
	case "cmd":
		return command.Execute(action, cmd, config)
	}

	return nil, nil
}
