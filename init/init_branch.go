package initialize

import (
	"github.com/neel1996/gitconvex/git/branch"
	"github.com/neel1996/gitconvex/git/branch/checkout"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/validator"
)

var (
	repoValidator   validator.Validator
	branchValidator validator.ValidatorWithStringFields
)

type InitBranch interface {
	BranchAdd() branch.Add
	BranchDelete() branch.Delete
	BranchCheckout() checkout.Factory
}

type initBranch struct {
	repo middleware.Repository
}

func init() {
	repoValidator = validator.NewRepoValidator()
	branchValidator = validator.NewBranchValidator()
}

func (i initBranch) BranchAdd() branch.Add {
	return branch.NewAddBranch(i.repo, branchValidator)
}

func (i initBranch) BranchDelete() branch.Delete {
	return branch.NewDeleteBranch(i.repo, branchValidator)
}

func (i initBranch) BranchCheckout() checkout.Factory {
	return checkout.NewCheckoutFactory(i.repo, repoValidator, branchValidator)
}

func NewInitBranch(repo middleware.Repository) InitBranch {
	return initBranch{repo: repo}
}
