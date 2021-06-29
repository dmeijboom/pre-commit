package runner

import "errors"

var ErrSkipped = errors.New("check was skipped")

type Runner struct {
	context *Context
}

func NewRunner(context *Context) *Runner {
	return &Runner{}
}

func (runner *Runner) RunAll(actions map[string]Action) *ActionResultIter {
	results := make(chan *ActionResult, len(actions))

	for name := range actions {
		go func(actionName string) {
			messages, err := actions[actionName].Run()
			if errors.Is(err, ErrSkipped) {
				results <- &ActionResult{
					Skipped:   true,
					ActionRef: actionName,
				}
				return
			}

			results <- &ActionResult{
				Err:       err,
				ActionRef: actionName,
				Messages:  messages,
			}
		}(name)
	}

	return &ActionResultIter{
		len:     len(actions),
		channel: results,
	}
}
