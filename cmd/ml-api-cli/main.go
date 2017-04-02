package main

import (
	"os"

	"github.com/melvin-laplanche/ml-api/cmd/ml-api-cli/generate"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:        "generate",
			Aliases:     []string{"n"},
			Usage:       "Generate new code",
			Subcommands: generate.SubCommands(),
		},
	}

	app.Run(os.Args)
}
