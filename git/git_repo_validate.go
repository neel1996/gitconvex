package git

import (
	"errors"
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/global"
	"os"
)

type RepoValidator interface {
	ValidateRepoWithPath() error
	ValidateGitRepository()
}

type repoValidator struct {
	repoPath string
	repo     *git2go.Repository
}

// ValidateRepoWithPath validates if the specified repoPath represents a valid git repo
func (r repoValidator) ValidateRepoWithPath() error {
	repoPath := r.repoPath

	_, err := os.Open(fmt.Sprintf("%v/.git/", repoPath))
	if err != nil {
		return errors.New("the selected folder is not a git repo")
	}

	repo, repoErr := git2go.OpenRepository(repoPath)
	if repoErr != nil {
		return errors.New("the selected folder is not a valid git repo")
	}

	_, headErr := repo.Head()
	if headErr != nil {
		logger.Log(fmt.Sprintf("Mind that the repo has no HEAD and a fresh commit is required -> %s", headErr.Error()), global.StatusWarning)
	}

	return nil
}

func (r repoValidator) ValidateGitRepository() {
	repo := r.repo

	defer func() {
		if panicValue := recover(); panicValue != nil {
			logger.Log(fmt.Sprintf("%v", panicValue), global.StatusError)
		}
	}()

	if repo == nil {
		panic(errors.New("repository is nil and cannot be used for performing the desired operation"))
	}
}

func NewRepoValidator(repoPath string, repo *git2go.Repository) RepoValidator {
	return repoValidator{
		repoPath: repoPath,
		repo:     repo,
	}
}
