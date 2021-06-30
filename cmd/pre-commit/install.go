package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ttacon/chalk"
	"github.com/urfave/cli/v2"
)

var installFiles = map[string]string{
	".git/hooks/pre-commit": "#!/bin/sh\npre-commit run\n",
	"pre-commit.json": `{"checks": [{
		"name": "Example",
		"cmd": "true",
		"when": [{
			"glob": "*.txt"
		}]
	}]}`,
}

func installCmd(ctx *cli.Context) error {
	for filename, content := range installFiles {
		_, err := os.Stat(filename)
		if ctx.Bool("force") || os.IsNotExist(err) {
			err := ioutil.WriteFile(
				filename,
				[]byte(content),
				0755,
			)
			if err != nil {
				return err
			}

			fmt.Printf("%s✓ wrote: %s%s\n", chalk.Green, filename, chalk.Reset)

			continue
		}

		fmt.Printf("%s✗ unable to write: %s%s\n", chalk.Red, filename, chalk.Reset)
	}

	return nil
}
