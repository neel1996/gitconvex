package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/neel1996/gitconvex-server/global"
)

// AddBranch adds a new branch to the target repo

func AddBranch(repo *git.Repository, branchName string) string {
	logger := global.Logger{}
	headRef, headErr := repo.Head()

	if headErr != nil {
		logger.Log(fmt.Sprintf("Unable to fetch HEAD -> %s", headErr.Error()), global.StatusError)
		return "BRANCH_ADD_FAILED"
	} else {
		ref := plumbing.NewHashReference(plumbing.ReferenceName(fmt.Sprintf("refs/heads/%v", branchName)), headRef.Hash())
		branchErr := repo.Storer.SetReference(ref)

		if branchErr != nil {
			logger.Log(fmt.Sprintf("Failed to add branch - %s - %s", branchName, branchErr.Error()), global.StatusError)
			return "BRANCH_ADD_FAILED"
		}

		logger.Log(fmt.Sprintf("Added new branch - %s to the repo", branchName), global.StatusInfo)
		return "BRANCH_CREATION_SUCCESS"
	}
}
