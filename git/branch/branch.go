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
	GitListBranches(chan ListOfBranches)
}

type Operation struct {
	Add      Add
	Checkout Checkout
	Compare  Compare
	Delete   Delete
	List     List
}

func (b Operation) GitAddBranch() (string, error) {
	err := b.Add.AddBranch()

	if err != nil {
		return global.BranchAddError, err
	}

	return global.BranchAddSuccess, nil
}

func (b Operation) GitCheckoutBranch() (string, error) {
	checkoutBranchResult := b.Checkout.CheckoutBranch()

	if checkoutBranchResult == global.BranchCheckoutError {
		return "", errors.New(global.BranchCheckoutError)
	}

	return checkoutBranchResult, nil
}

func (b Operation) GitCompareBranches() ([]*model.BranchCompareResults, error) {
	branchDiff := b.Compare.CompareBranch()

	if len(branchDiff) == 0 {
		return []*model.BranchCompareResults{}, errors.New("no difference between the two branches")
	}

	return branchDiff, nil
}

func (b Operation) GitDeleteBranch() (*model.BranchDeleteStatus, error) {
	deleteStatus := b.Delete.DeleteBranch()

	if deleteStatus.Status == global.BranchDeleteError {
		return nil, errors.New("branch deletion failed")
	}

	return deleteStatus, nil
}

func (b Operation) GitListBranches(branchChannel chan ListOfBranches) {
	b.List.ListBranches(branchChannel)
}
