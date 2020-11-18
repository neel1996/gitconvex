package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
)

// DeleteBranch deleted a branch from the repo
// If forceFlag is true then it will forcefully delete a branch
// If forceFlag is false, then the branch status will be checked for unmerged changes and then it will be removed from the repo
func DeleteBranch(repo *git.Repository, branchName string, forceFlag bool) *model.BranchDeleteStatus {
	var branchErr error
	logger := global.Logger{}

	headRef, _ := repo.Head()
	ref := plumbing.NewHashReference(plumbing.ReferenceName(fmt.Sprintf("refs/heads/%v", branchName)), headRef.Hash())

	if forceFlag {
		logger.Log("Deleting branch "+branchName+" forcefully", global.StatusInfo)
		branchErr = repo.Storer.RemoveReference(ref.Name())
	} else {
		b, bErr := repo.Branch(branchName)
		if bErr != nil {
			fmt.Println(bErr.Error())
		} else {
			logger.Log("Deleting branch "+b.Name, global.StatusInfo)
			branchErr = repo.Storer.RemoveReference(ref.Name())
		}
	}

	if branchErr != nil {
		logger.Log(fmt.Sprintf("Failed to delete branch %s -> %v", branchName, branchErr.Error()), global.StatusError)
		return &model.BranchDeleteStatus{Status: "BRANCH_DELETE_FAILED"}
	}

	return &model.BranchDeleteStatus{
		Status: "BRANCH_DELETE_SUCCESS",
	}
}
