package term

import (
	"fmt"

	"github.com/ttacon/chalk"
)

const (
	divider = "█"
	width   = 20

	ansiClearLine = "\u001b[2K\r"
)

type ProgressBar struct {
	max    int
	cur    int
	dirty  bool
	status string
}

func NewProgressBar(max int) *ProgressBar {
	return &ProgressBar{
		max: max,
	}
}

func (bar *ProgressBar) Clear() {
	if !bar.dirty {
		return
	}

	fmt.Print(ansiClearLine)

	bar.dirty = false
}

func (bar *ProgressBar) Status(status string) *ProgressBar {
	bar.status = status

	return bar
}

func (bar *ProgressBar) Tick() *ProgressBar {
	bar.cur++

	return bar
}

func (bar *ProgressBar) Render() {
	bar.Clear()

	body := ""
	ratio := float32(bar.cur) / float32(bar.max)

	for i := 0; i < width; i++ {
		if i < int(width*ratio) {
			body += divider
			continue
		}

		body += " "
	}

	fmt.Printf(
		"%s%s%s %s%d/%d%s %s[%s%s%s%s%s]%s",
		chalk.Blue,
		bar.status,
		chalk.Reset,
		chalk.White,
		bar.cur,
		bar.max,
		chalk.Reset,

		chalk.Blue,
		chalk.Reset,
		chalk.White,
		body,
		chalk.Reset,
		chalk.Blue,
		chalk.Reset,
	)

	bar.dirty = true
}

func (bar *ProgressBar) Done() {
	bar.Clear()

	fmt.Print("\n")
}
