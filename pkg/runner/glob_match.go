package runner

import (
	"path"
	"path/filepath"
)

func MatchGlob(action Action, patterns []string) Action {
	return ActionFunc(func(ctx *Context) ([]Message, error) {
		match := false

		if ctx.Git != nil {
		filterLoop:
			for _, name := range ctx.Git.Dirty {
				for _, pattern := range patterns {
					filename := path.Base(name)

					if ok, _ := filepath.Match(pattern, filename); ok {
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
