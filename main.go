package main

import (
	"fmt"
	"os"

	"gitea.darkeli.com/yezi/git-bump/cli"
)

var (
	Version  = "unset"
	Revision = "unset"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	var opts cli.Options
	cli := cli.NewCLI(opts)

	args, err := cli.ParseArgs(args)
	if err != nil {
		return 1
	}

	if err := cli.Run(args); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %v\n", err)
		return 1
	}
	return 0
}
