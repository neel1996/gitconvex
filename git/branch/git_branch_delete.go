package branch

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/global"
)

type Delete interface {
	DeleteBranch() error
}

type deleteBranch struct {
	repo       *git2go.Repository
	branchName string
}

// DeleteBranch deletes a branch from the repo
func (d deleteBranch) DeleteBranch() error {
	repo := d.repo
	branchName := d.branchName

	validationErr := NewBranchFieldsValidation(repo, branchName).ValidateBranchFields()
	if validationErr != nil {
		logger.Log(validationErr.Error(), global.StatusError)
		return validationErr
	}

	branch, deleteBranchErr := repo.LookupBranch(branchName, git2go.BranchLocal)
	if deleteBranchErr != nil {
		logger.Log(fmt.Sprintf("Failed to delete branch %s -> %v", branchName, deleteBranchErr.Error()), global.StatusError)
		return deleteBranchErr
	}

	deleteErr := branch.Delete()
	if deleteErr != nil {
		logger.Log(fmt.Sprintf("Failed to delete branch %s -> %v", branchName, deleteErr.Error()), global.StatusError)
		return deleteErr
	}

	logger.Log(fmt.Sprintf("Branch - %s has been removed from the repo", branchName), global.StatusInfo)
	return nil
}

func NewDeleteBranch(repo *git2go.Repository, branchName string) Delete {
	return deleteBranch{
		repo:       repo,
		branchName: branchName,
	}
}
