package branch

import (
	"github.com/neel1996/gitconvex/git/middleware"
)

type Validation interface {
	ValidateBranchFields(branchNames ...string) error
}

type validateBranch struct {
	repo middleware.Repository
}

func (v validateBranch) ValidateBranchFields(branchNames ...string) error {
	if v.repo == nil {
		return NilRepoError
	}

	if len(branchNames) == 0 {
		return EmptyBranchNameError
	}

	for _, branchName := range branchNames {
		if branchName == "" {
			return EmptyBranchNameError
		}
	}

	return nil
}

func NewBranchFieldsValidation(repo middleware.Repository) Validation {
	return validateBranch{
		repo: repo,
	}
}
