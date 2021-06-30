package main

import (
	"fmt"
	"github.com/superhawk610/terminal"
	"os"

	"github.com/superhawk610/bar"
	"github.com/ttacon/chalk"
	"github.com/urfave/cli/v2"

	"github.com/dmeijboom/pre-commit/pkg/config"
	"github.com/dmeijboom/pre-commit/pkg/runner"
)

func runCmd(c *cli.Context) error {
	cfg, err := config.Load("pre-commit.json")
	if err != nil {
		return err
	}

	term := terminal.New()

	progress := bar.NewWithOpts(
		bar.WithDimensions(len(cfg.Checks), 25),
		bar.WithFormat(
			fmt.Sprintf(
				" %s:state%s :percent :bar",
				chalk.Blue,
				chalk.Reset,
			),
		),
	)

	progress.Update(0, bar.Context{
		bar.Ctx("state", "init.."),
	})

	workdir, err := os.Getwd()
	if err != nil {
		return err
	}

	ctx, err := runner.NewContext(workdir)
	if err != nil {
		return err
	}

	r := runner.NewRunner(ctx)
	actions := map[string]runner.Action{}

	for _, check := range cfg.Checks {
		action := runner.Cmd(check.Cmd)

		if !c.Bool("all") && len(check.When) > 0 {
			patterns := make([]string, len(check.When))

			for i, is := range check.When {
				patterns[i] = is.Glob
			}

			action = runner.MatchGlob(action, patterns)
		}

		actions[check.Name] = action
	}

	exitOk := true
	iter := r.RunAll(actions)

	var (
		done   bool
		result *runner.ActionResult
	)

	for !done {
		result, done = iter.Next()

		// clear the entire line
		term.ClearLine()

		if result.Skipped {
			progress.Interruptf(
				"%s✓%s %s %s[skipped]%s",
				chalk.Green,
				chalk.Reset,
				result.ActionRef,
				chalk.White,
				chalk.Reset,
			)
		} else {
			if result.Ok() {
				progress.Interruptf("%s✓ %s%s", chalk.Green, result.ActionRef, chalk.Reset)
			} else {
				exitOk = false
				progress.Interruptf("%s✗ %s%s", chalk.Red, result.ActionRef, chalk.Reset)
			}
		}

		progress.TickAndUpdate(bar.Context{
			bar.Ctx("state", "running.."),
		})
	}

	progress.Done()

	// clear the progress bar
	fmt.Print("\033[F")
	term.ClearLine()

	fmt.Println("")
	fmt.Printf("%d checks finished\n", len(cfg.Checks))

	if !exitOk {
		os.Exit(1)
	}

	return nil
}
