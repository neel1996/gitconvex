package git

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/global"
	"go/types"
	"os"
)

// RepoValidator validates if the specified repoPath represents a valid git repo
func RepoValidator(repoPath string) (string, error) {
	_, err := os.Open(fmt.Sprintf("%v/.git/", repoPath))
	if err != nil {
		return "", types.Error{Msg: "The selected folder is not a git repo"}
	}

	repo, repoErr := git2go.OpenRepository(repoPath)
	if repoErr != nil {
		return "", types.Error{Msg: "The selected folder is not a valid git repo"}
	}

	_, headErr := repo.Head()
	if headErr != nil {
		logger.Log(fmt.Sprintf("Mind that the repo has no HEAD and a fresh commit is required -> %s", headErr.Error()), global.StatusWarning)
	}

	return "Repo is valid!", nil
}
