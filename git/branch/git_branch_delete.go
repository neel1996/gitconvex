package branch

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/validator"
)

type Delete interface {
	DeleteBranch(branchName string) error
}

type deleteBranch struct {
	repo            middleware.Repository
	branchValidator validator.ValidatorWithStringFields
}

func (d deleteBranch) DeleteBranch(branchName string) error {
	repo := d.repo

	validationErr := d.branchValidator.ValidateWithFields(branchName)
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

func NewDeleteBranch(repo middleware.Repository, branchValidator validator.ValidatorWithStringFields) Delete {
	return deleteBranch{
		repo:            repo,
		branchValidator: branchValidator,
	}
}
