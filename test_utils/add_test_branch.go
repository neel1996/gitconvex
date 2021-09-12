package test_utils

import (
	"github.com/neel1996/gitconvex/git/branch"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/validator"
)

func AddNewTestLocalBranch(repo middleware.Repository, branchName string) {
	_ = branch.NewAddBranch(
		repo,
		validator.NewBranchValidator(),
	).AddBranch(branchName, false, nil)
}
