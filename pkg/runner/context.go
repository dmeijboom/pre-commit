package runner

import (
	"os"
	"os/exec"
	"strings"

	"github.com/go-git/go-git/v5"
)

type GitContext struct {
	Name  string
	Hash  string
	Dirty []string
}

type Context struct {
	Git *GitContext
}

func NewContext(root string) (*Context, error) {
	repo, err := git.PlainOpen(root)
	if os.IsNotExist(err) {
		return &Context{}, nil
	} else if err != nil {
		return nil, err
	}

	head, err := repo.Head()
	if err != nil {
		return nil, err
	}

	// workaround for worktree.Status() using go-git being extremely slow
	statusCmd := exec.Command("git", "-C", root, "status", "--untracked=no", "--porcelain=v2")

	output, err := statusCmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")

	dirtyFiles := []string{}

	for _, line := range lines {
		components := strings.Split(line, " ")

		if len(components) != 9 {
			continue
		}

		fileStatus := git.StatusCode(components[1][0])

		if fileStatus == git.Modified ||
			fileStatus == git.Added ||
			fileStatus == git.Renamed {
			dirtyFiles = append(dirtyFiles, components[8])
		}
	}

	return &Context{
		Git: &GitContext{
			Name:  head.Name().String(),
			Hash:  head.Hash().String(),
			Dirty: dirtyFiles,
		},
	}, nil
}
