package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
)

type DeleteBranchInterface interface {
	DeleteBranch() *model.BranchDeleteStatus
}

type DeleteBranchInputs struct {
	Repo       *git.Repository
	BranchName string
	ForceFlag  bool
}

// DeleteBranch deleted a branch from the repo
// If forceFlag is true then it will forcefully delete a branch
// If forceFlag is false, then the branch status will be checked for unmerged changes and then it will be removed from the repo
func (inputs DeleteBranchInputs) DeleteBranch() *model.BranchDeleteStatus {
	var branchErr error
	logger := global.Logger{}

	repo := inputs.Repo
	branchName := inputs.BranchName
	forceFlag := inputs.ForceFlag

	headRef, _ := repo.Head()
	ref := plumbing.NewHashReference(plumbing.ReferenceName(fmt.Sprintf("refs/heads/%v", branchName)), headRef.Hash())

	if forceFlag {
		logger.Log("Deleting branch "+branchName+" forcefully", global.StatusInfo)
		branchErr = repo.Storer.RemoveReference(ref.Name())
	} else {
		b, bErr := repo.Branch(branchName)
		if bErr != nil {
			logger.Log(fmt.Sprintf("Failed to delete branch %s -> %v", branchName, bErr.Error()), global.StatusError)
			return &model.BranchDeleteStatus{Status: global.BranchDeleteError}
		} else {
			logger.Log("Deleting branch "+b.Name, global.StatusInfo)
			branchErr = repo.Storer.RemoveReference(ref.Name())
		}
	}

	if branchErr != nil {
		logger.Log(fmt.Sprintf("Failed to delete branch %s -> %v", branchName, branchErr.Error()), global.StatusError)
		return &model.BranchDeleteStatus{Status: global.BranchDeleteError}
	}

	return &model.BranchDeleteStatus{
		Status: global.BranchDeleteSuccess,
	}
}
