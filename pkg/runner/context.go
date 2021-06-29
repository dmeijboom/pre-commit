package runner

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
)

type GitContext struct {
	Name string
	Hash string
}

type Context struct {
	Git *GitContext
}

func NewContext(root string) (*Context, error) {
	repo, err := git.PlainOpen(fmt.Sprintf("%s/.git", root))
	if os.IsNotExist(err) {
		return &Context{}, nil
	} else if err != nil {
		return nil, err
	}

	head, err := repo.Head()
	if err != nil {
		return nil, err
	}

	return &Context{
		Git: &GitContext{
			Name: head.Name().String(),
			Hash: head.Hash().String(),
		},
	}, nil
}
