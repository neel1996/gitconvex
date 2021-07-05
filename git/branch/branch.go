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
	GitListBranches() (model.ListOfBranches, error)
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
	err := b.Checkout.CheckoutBranch()

	if err != nil {
		return global.BranchCheckoutError, err
	}

	return global.BranchCheckoutSuccess, nil
}

func (b Operation) GitCompareBranches() ([]*model.BranchCompareResults, error) {
	branchDiff := b.Compare.CompareBranch()

	if len(branchDiff) == 0 {
		return []*model.BranchCompareResults{}, errors.New("no difference between the two branches")
	}

	return branchDiff, nil
}

func (b Operation) GitDeleteBranch() (*model.BranchDeleteStatus, error) {
	err := b.Delete.DeleteBranch()

	if err != nil {
		return &model.BranchDeleteStatus{Status: global.BranchDeleteError}, err
	}

	return &model.BranchDeleteStatus{Status: global.BranchDeleteSuccess}, nil
}

func (b Operation) GitListBranches() (model.ListOfBranches, error) {
	listOfBranches, err := b.List.ListBranches()
	if err != nil {
		return model.ListOfBranches{}, err
	}

	return listOfBranches, nil
}
