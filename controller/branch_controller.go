package controller

import (
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/branch"
	"github.com/neel1996/gitconvex/git/branch/checkout"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/graph/model"
	"github.com/neel1996/gitconvex/models"
)

type BranchController interface {
	GitAddBranch(branchName string, isRemoteBranch bool, targetCommit *git2go.Commit) (string, error)
	GitCheckoutBranch(branchName string) (string, error)
	GitDeleteBranch(branchName string) (*model.BranchDeleteStatus, error)
	GitCompareBranches(baseBranch string, compareBranch string) ([]*model.BranchCompareResults, error)
	GitListBranches() (model.ListOfBranches, error)
}

type branchController struct {
	add             branch.Add
	checkoutFactory checkout.Factory
	compare         branch.Compare
	delete          branch.Delete
	list            branch.List
}

func (b branchController) GitAddBranch(branchName string, isRemoteBranch bool, targetCommit *git2go.Commit) (string, error) {
	err := b.add.AddBranch(branchName, isRemoteBranch, targetCommit)

	if err != nil {
		return global.BranchAddError, err
	}

	return global.BranchAddSuccess, nil
}

func (b branchController) GitCheckoutBranch(branchName string) (string, error) {
	err := b.checkoutFactory.GetCheckoutAction(branchName).CheckoutBranch()

	if err != nil {
		return global.BranchCheckoutError, err
	}

	return global.BranchCheckoutSuccess, nil
}

func (b branchController) GitCompareBranches(baseBranch string, compareBranch string) ([]*model.BranchCompareResults, error) {
	return b.compare.CompareBranch(baseBranch, compareBranch)
}

func (b branchController) GitDeleteBranch(branchName string) (*model.BranchDeleteStatus, error) {
	err := b.delete.DeleteBranch(branchName)

	if err != nil {
		return &model.BranchDeleteStatus{Status: global.BranchDeleteError}, err
	}

	return &model.BranchDeleteStatus{Status: global.BranchDeleteSuccess}, nil
}

func (b branchController) GitListBranches() (model.ListOfBranches, error) {
	listOfBranches, err := b.list.ListBranches()
	if err != nil {
		return model.ListOfBranches{}, err
	}

	return listOfBranches, nil
}

func NewBranchController(branchModel models.Branch) BranchController {
	return branchController{
		add:             branchModel.Add,
		checkoutFactory: branchModel.Checkout,
		compare:         branchModel.Compare,
		delete:          branchModel.Delete,
		list:            branchModel.List,
	}
}
