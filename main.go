package main

import (
	"github.com/tkanos/twx/cmd"
	"github.com/tkanos/twx/cmd/context"
)

var version string

func main() {
	context.Version = version
	cmd.Execute()
}
