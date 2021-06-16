package branch

import (
	"errors"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/graph/model"
)

var logger global.Logger

type Branch interface {
	GitAddBranch() (string, error)
	GitCheckoutBranch() (string, error)
	GitCompareBranches() ([]*model.BranchCompareResults, error)
}

type branch struct {
	addBranch       AddBranch
	checkoutBranch  BranchCheckout
	compareBranches BranchCompare
}

func (b branch) GitAddBranch() (string, error) {
	addBranchResult := b.addBranch.AddBranch()

	if addBranchResult == global.BranchAddError {
		return "", errors.New(global.BranchAddError)
	}

	return addBranchResult, nil
}

func (b branch) GitCheckoutBranch() (string, error) {
	checkoutBranchResult := b.checkoutBranch.CheckoutBranch()

	if checkoutBranchResult == global.BranchCheckoutError {
		return "", errors.New(global.BranchCheckoutError)
	}

	return checkoutBranchResult, nil
}

func (b branch) GitCompareBranches() ([]*model.BranchCompareResults, error) {
	branchDiff := b.compareBranches.CompareBranch()

	if len(branchDiff) == 0 {
		return []*model.BranchCompareResults{}, errors.New("no difference between the two branches")
	}

	return branchDiff, nil
}

func NewBranchOperation(addBranch AddBranch, branchCheckout BranchCheckout, branchCompare BranchCompare) Branch {
	return branch{
		addBranch:       addBranch,
		checkoutBranch:  branchCheckout,
		compareBranches: branchCompare,
	}
}
