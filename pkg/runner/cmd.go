package runner

import (
	"fmt"
	"os/exec"
)

type CmdAction struct {
	cmd string
}

func Cmd(cmd string) Action {
	return &CmdAction{cmd: cmd}
}

func (runner *CmdAction) Run(_ *Context) ([]Message, error) {
	var messages []Message

	cmd := exec.Command("sh", "-c", runner.cmd)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if len(exitErr.Stderr) > 0 {
				messages = append(messages, Message{
					Type: WarningMessage,
					Body: string(exitErr.Stderr),
				})
			}

			messages = append(messages, Message{
				Type: ErrorMessage,
				Body: fmt.Sprintf("command exited with status: %d", exitErr.ExitCode()),
			})
		}

		return messages, err
	}

	if len(output) > 0 {
		messages = append(messages, Message{
			Type: NoticeMessage,
			Body: string(output),
		})
	}

	return messages, nil
}
