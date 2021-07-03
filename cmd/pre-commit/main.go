package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:  "pre-commit",
		Usage: "Make sure your code is OK before shipping",
		ExitErrHandler: func(context *cli.Context, err error) {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		},
		Commands: []*cli.Command{
			{
				Name:   "run",
				Usage:  "Run the checks",
				Action: runCmd,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "all",
						Usage: "Run all checks (ignore when conditions)",
					},
				},
			},
			{
				Name:  "install",
				Usage: "Install git-hook and example config",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "force",
						Usage: "Overwrite files if already exists",
					},
				},
				Action: installCmd,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
