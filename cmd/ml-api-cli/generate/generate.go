package generate

import "github.com/urfave/cli"

// SubCommands returns a list of SubCommands for the "generate"" command
func SubCommands() []cli.Command {
	return []cli.Command{
		{
			Name:  "model",
			Usage: "generate ModelName",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "table, t",
					Usage: "Name of the SQL table used by this model",
				},
				cli.StringFlag{
					Name:   "file, f",
					Usage:  "Name of the file to parse",
					EnvVar: "GOFILE",
				},
				cli.StringFlag{
					Name:   "package, p",
					Usage:  "Name of the go package to parse",
					EnvVar: "GOPACKAGE",
				},
				cli.StringFlag{
					Name:  "exclude, e",
					Usage: "Methods to exclude from the generated file",
				},
			},
			Action: func(c *cli.Context) error {
				GenModel(c)
				return nil
			},
		},
	}
}
