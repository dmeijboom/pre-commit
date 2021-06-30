package runner

import (
	"path"
)

func MatchDir(action Action, dirs []string) Action {
	return ActionFunc(func(ctx *Context) ([]Message, error) {
		match := false

		if ctx.Git != nil {
		filterLoop:
			for _, name := range ctx.Git.Dirty {
				for _, dirname := range dirs {
					if dirname == path.Dir(name) {
						match = true
						break filterLoop
					}
				}
			}
		}

		if !match {
			return nil, ErrSkipped
		}

		return action.Run(ctx)
	})
}
