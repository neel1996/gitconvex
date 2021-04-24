package git

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex-server/global"
)

type AddBranchInterface interface {
	AddBranch() string
}

type AddBranchInput struct {
	Repo         *git2go.Repository
	BranchName   string
	RemoteSwitch bool
	TargetCommit *git2go.Commit
}

// AddBranch adds a new branch to the target repo
func (input AddBranchInput) AddBranch() string {
	targetCommit := input.TargetCommit
	repo := input.Repo
	branchName := input.BranchName
	head, headErr := repo.Head()

	logger.Log(fmt.Sprintf("Adding new branch -> %s", branchName), global.StatusInfo)
	if headErr != nil {
		logger.Log(fmt.Sprintf("Unable to fetch HEAD -> %s", headErr.Error()), global.StatusError)
		return global.BranchAddError
	} else {
		if targetCommit == nil {
			targetCommit, _ = repo.LookupCommit(head.Target())
			if targetCommit == nil {
				logger.Log("Target commit is nil", global.StatusError)
				return global.BranchAddError
			}
		}
		_, branchErr := repo.CreateBranch(branchName, targetCommit, false)

		if branchErr != nil {
			logger.Log(fmt.Sprintf("Failed to add branch - %s - %s", branchName, branchErr.Error()), global.StatusError)
			return global.BranchAddError
		}

		logger.Log(fmt.Sprintf("Added new branch - %s to the repo", branchName), global.StatusInfo)
		return global.BranchAddSuccess
	}
}
