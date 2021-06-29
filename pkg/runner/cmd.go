package runner

import (
	"os/exec"
)

type CmdAction struct {
	cmd string
}

func Cmd(cmd string) Action {
	return &CmdAction{cmd: cmd}
}

func (runner *CmdAction) Run(_ *Context) ([]Message, error) {
	cmd := exec.Command("sh", "-c", runner.cmd)

	_, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return []Message{}, nil
}
