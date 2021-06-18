package branch

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/graph/model"
)

type Delete interface {
	DeleteBranch() *model.BranchDeleteStatus
}

type deleteBranch struct {
	Repo       *git2go.Repository
	BranchName string
}

// DeleteBranch deletes a branch from the repo
func (inputs deleteBranch) DeleteBranch() *model.BranchDeleteStatus {
	repo := inputs.Repo
	branchName := inputs.BranchName
	branch, deleteBranchErr := repo.LookupBranch(branchName, git2go.BranchLocal)

	if deleteBranchErr != nil {
		logger.Log(fmt.Sprintf("Failed to delete branch %s -> %v", branchName, deleteBranchErr.Error()), global.StatusError)
		return &model.BranchDeleteStatus{Status: global.BranchDeleteError}
	}

	deleteErr := branch.Delete()
	if deleteErr != nil {
		logger.Log(fmt.Sprintf("Failed to delete branch %s -> %v", branchName, deleteErr.Error()), global.StatusError)
		return &model.BranchDeleteStatus{Status: global.BranchDeleteError}
	}

	logger.Log(fmt.Sprintf("Branch - %s has been removed from the repo", branchName), global.StatusInfo)
	return &model.BranchDeleteStatus{
		Status: global.BranchDeleteSuccess,
	}
}

func NewDeleteBranch(repo *git2go.Repository, branchName string) Delete {
	return deleteBranch{
		Repo:       repo,
		BranchName: branchName,
	}
}
