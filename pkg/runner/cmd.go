package runner

import (
	"bytes"
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

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cmd := exec.Command("sh", "-c", runner.cmd)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()

	if stdout.Len() > 0 {
		messages = append(messages, Message{
			Type: NoticeMessage,
			Body: stdout.String(),
		})
	}
	if stderr.Len() > 0 {
		messages = append(messages, Message{
			Type: WarningMessage,
			Body: stderr.String(),
		})
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			messages = append(messages, Message{
				Type: ErrorMessage,
				Body: fmt.Sprintf("command exited with status: %d", exitErr.ExitCode()),
			})
		}

		return messages, err
	}

	return messages, nil
}
