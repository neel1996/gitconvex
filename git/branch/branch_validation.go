package branch

import (
	"errors"
	git2go "github.com/libgit2/git2go/v31"
)

type Validation interface {
	ValidateBranchFields() error
}

type validateBranch struct {
	repo       *git2go.Repository
	branchName []string
}

func (v validateBranch) ValidateBranchFields() error {
	if v.repo == nil {
		err := "repo is nil"
		return errors.New(err)
	}

	for _, branchName := range v.branchName {
		if branchName == "" {
			err := "branch name is empty"
			return errors.New(err)
		}
	}

	return nil
}

func NewBranchFieldsValidation(repo *git2go.Repository, branchName ...string) Validation {
	return validateBranch{
		repo:       repo,
		branchName: branchName,
	}
}
