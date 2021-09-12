package checkout

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/global"
)

type checkOutLocalBranch struct {
	repo       middleware.Repository
	branchName string
}

func (c checkOutLocalBranch) CheckoutBranch() error {
	repo := c.repo
	branchName := c.branchName

	logger.Log(fmt.Sprintf("Initiating local branch checkout for %s", branchName), global.StatusInfo)

	branch, branchErr := repo.LookupBranch(branchName, git2go.BranchLocal)
	if branchErr != nil {
		return c.LogAndReturnError(branchErr)
	}

	commit, lookupCommitErr := repo.LookupCommitV2(branch.Target())
	if lookupCommitErr != nil {
		return c.LogAndReturnError(lookupCommitErr)
	}

	tree, treeErr := commit.Tree()
	if treeErr != nil {
		return c.LogAndReturnError(treeErr)
	}

	checkoutErr := repo.CheckoutTree(tree, &git2go.CheckoutOptions{
		Strategy:       git2go.CheckoutSafe,
		DisableFilters: false,
	})

	if checkoutErr != nil {
		return c.LogAndReturnError(checkoutErr)
	}

	err := repo.SetHead(c.GenerateBranchFields().ReferenceBranchName)
	if err != nil {
		return c.LogAndReturnError(err)
	}

	logger.Log(fmt.Sprintf("Local branch %v has been checked out", c.branchName), global.StatusInfo)
	return nil
}

func (c checkOutLocalBranch) GenerateBranchFields() BranchDetails {
	referenceBranchName := "refs/heads/" + c.branchName

	return BranchDetails{
		BranchName:          c.branchName,
		ReferenceBranchName: referenceBranchName,
	}
}

func (c checkOutLocalBranch) LogAndReturnError(err error) error {
	logger.Log(err.Error(), global.StatusError)
	return err
}

func NewCheckOutLocalBranch(repo middleware.Repository, branchName string) Checkout {
	return checkOutLocalBranch{
		repo:       repo,
		branchName: branchName,
	}
}
