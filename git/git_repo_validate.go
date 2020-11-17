package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"go/types"
	"os"
)

// RepoValidator validates if the specified repoPath represents a valid git repo

func RepoValidator(repoPath string) (string, error) {
	_, err := os.Open(fmt.Sprintf("%v/.git/", repoPath))
	if err != nil {
		return "", types.Error{Msg: "The selected folder is not a git repo"}
	}

	repo, repoErr := git.PlainOpen(repoPath)
	if repoErr != nil {
		return "", types.Error{Msg: "The selected folder is not a valid git repo"}
	}

	_, headErr := repo.Head()
	if headErr != nil {
		return "", types.Error{Msg: "The selected folder is not a valid git repo"}
	}

	return "Repo is valid!", nil
}
