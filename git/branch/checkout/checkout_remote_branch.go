package checkout

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/branch"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/global"
	"strings"
)

type checkOutRemoteBranch struct {
	repo       middleware.Repository
	branchName string
	addBranch  branch.Add
}

func (c checkOutRemoteBranch) CheckoutBranch() error {
	repo := c.repo
	branchName := c.branchName

	logger.Log(fmt.Sprintf("Initiating remote branch checkout for %s", branchName), global.StatusInfo)

	branchFields := c.GenerateBranchFields()

	logger.Log(fmt.Sprintf("Branch - %s is a remote branch. Trying with intermediate remote fetch!", branchName), global.StatusWarning)
	remoteBranch, remoteBranchErr := repo.LookupBranch(branchFields.RemoteBranchName, git2go.BranchRemote)
	if remoteBranchErr != nil {
		return c.LogAndReturnError(remoteBranchErr)
	}

	remoteHead := remoteBranch.Target()
	remoteCommit, remoteCommitErr := repo.LookupCommitV2(remoteHead)
	if remoteCommitErr != nil {
		return c.LogAndReturnError(remoteCommitErr)

	}

	remoteTree, remoteTreeErr := remoteCommit.Tree()
	if remoteTreeErr != nil {
		return c.LogAndReturnError(remoteTreeErr)
	}

	if checkoutErr := repo.CheckoutTree(remoteTree, &git2go.CheckoutOptions{Strategy: git2go.CheckoutSafe}); checkoutErr != nil {
		return c.LogAndReturnError(checkoutErr)
	}

	_, localLookupErr := repo.LookupBranch(branchFields.BranchName, git2go.BranchLocal)
	if localLookupErr != nil {
		logger.Log(localLookupErr.Error(), global.StatusError)
		gitCommit := remoteCommit.GetGitCommit()
		return c.addAndCheckoutNewBranch(gitCommit, branchFields.ReferenceBranchName)
	}

	if err := repo.SetHead(branchFields.ReferenceBranchName); err != nil {
		return c.LogAndReturnError(err)
	}

	logger.Log(fmt.Sprintf("Remote branch %v has been checked out", branchFields.RemoteBranchName), global.StatusInfo)

	return nil
}

func (c checkOutRemoteBranch) addAndCheckoutNewBranch(remoteCommit *git2go.Commit, referenceBranchName string) error {
	logger.Log("Adding local branch for remote reference", global.StatusInfo)

	branchAddError := c.addBranch.AddBranch(c.branchName, false, remoteCommit)
	if branchAddError != nil {
		logger.Log("Error while adding branch "+branchAddError.Error(), global.StatusError)
		return branchAddError
	}

	if err := c.repo.SetHead(referenceBranchName); err != nil {
		return err
	}

	return nil
}

func (c checkOutRemoteBranch) GenerateBranchFields() BranchDetails {
	referenceBranchName := "refs/heads/" + c.splitAndReturnBranchName("/", c.branchName)
	remoteBranchName := c.splitAndReturnBranchName("remotes/", c.branchName)
	branchName := c.splitAndReturnBranchName("/", c.branchName)

	return BranchDetails{
		BranchName:          branchName,
		ReferenceBranchName: referenceBranchName,
		RemoteBranchName:    remoteBranchName,
	}
}

func (c checkOutRemoteBranch) LogAndReturnError(err error) error {
	logger.Log(err.Error(), global.StatusError)
	return err
}

func (c checkOutRemoteBranch) splitAndReturnBranchName(delimiter string, branchName string) string {
	splitString := strings.Split(branchName, delimiter)
	return splitString[len(splitString)-1]
}

func NewCheckoutRemoteBranch(repo middleware.Repository, branchName string, addBranch branch.Add) Checkout {
	return checkOutRemoteBranch{
		repo:       repo,
		branchName: branchName,
		addBranch:  addBranch,
	}
}
