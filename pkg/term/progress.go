package term

import (
	"fmt"
	"time"

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
	prev   int
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
	bar.prev = bar.cur
	bar.cur++

	return bar
}

func (bar *ProgressBar) render(barLen int) {
	bar.Clear()

	body := ""

	for i := 0; i < width; i++ {
		if i < barLen {
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

func (bar *ProgressBar) Render() {
	prev := int(width * float64(bar.prev) / float64(bar.max))
	cur := int(width * float64(bar.cur) / float64(bar.max))

	if prev == cur {
		bar.render(cur)

		return
	}

	for i := prev; i < cur; i++ {
		bar.render(i)

		time.Sleep(10 * time.Millisecond)
	}
}

func (bar *ProgressBar) Done() {
	bar.Clear()

	fmt.Print("\n")
}
