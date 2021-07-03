package main

import (
	"fmt"
	"os"

	"github.com/ttacon/chalk"
	"github.com/urfave/cli/v2"

	"github.com/dmeijboom/pre-commit/pkg/config"
	"github.com/dmeijboom/pre-commit/pkg/runner"
	"github.com/dmeijboom/pre-commit/pkg/term"
)

func parseActions(checks []config.Check, ignoreCond bool) map[string]runner.Action {
	actions := map[string]runner.Action{}

	for _, check := range checks {
		action := runner.Cmd(check.Cmd)

		if !ignoreCond && len(check.When) > 0 {
			patterns := []string{}
			dirs := []string{}

			for _, is := range check.When {
				if is.Glob != "" {
					patterns = append(patterns, is.Glob)
				}

				if is.Dir != "" {
					dirs = append(dirs, is.Dir)
				}
			}

			if len(patterns) > 0 {
				action = runner.MatchGlob(action, patterns)
			}

			if len(dirs) > 0 {
				action = runner.MatchDir(action, dirs)
			}
		}

		actions[check.Name] = action
	}

	return actions
}

func runCmd(c *cli.Context) error {
	cfg, err := config.Load("pre-commit.json")
	if err != nil {
		return err
	}

	progress := term.NewProgressBar(len(cfg.Checks))
	progress.Status("running..").Render()

	workdir, err := os.Getwd()
	if err != nil {
		return err
	}

	ctx, err := runner.NewContext(workdir)
	if err != nil {
		return err
	}

	r := runner.NewRunner(ctx)

	exitOk := true
	actions := parseActions(cfg.Checks, c.Bool("all"))
	iter := r.RunAll(actions)

	var (
		done   bool
		result *runner.ActionResult
	)

	for !done {
		result, done = iter.Next()

		progress.Clear()

		if result.Skipped {
			fmt.Printf(
				"%s✓%s %s %s[skipped]%s\n",
				chalk.Green,
				chalk.Reset,
				result.ActionRef,
				chalk.White,
				chalk.Reset,
			)
		} else {
			if result.Ok() {
				fmt.Printf("%s✓ %s%s\n", chalk.Green, result.ActionRef, chalk.Reset)
			} else {
				exitOk = false

				fmt.Printf("%s✗ %s%s\n", chalk.Red, result.ActionRef, chalk.Reset)

				for _, message := range result.Messages {
					var color chalk.Color

					switch message.Type {
					case runner.NoticeMessage:
						color = chalk.Red
					case runner.WarningMessage:
						color = chalk.Yellow
					case runner.ErrorMessage:
						color = chalk.Red
					default:
						color = chalk.Black
					}

					style := chalk.Dim.NewStyle()
					style.Foreground(color)

					fmt.Printf("%s%s%s\n", style.String(), message.Body, chalk.Reset)
				}
			}
		}

		progress.Tick().Render()
	}

	progress.Done()

	fmt.Printf("%d checks finished\n", len(cfg.Checks))

	if !exitOk {
		os.Exit(1)
	}

	return nil
}
