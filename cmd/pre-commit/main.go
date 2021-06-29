package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
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
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
