package git

import (
	"errors"
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/global"
	"strings"
)

type BranchCheckout interface {
	CheckoutBranch() string
}

type branchCheckout struct {
	Repo       *git2go.Repository
	BranchName string
}

type validatedBranchDetails struct {
	branchName          string
	referenceBranchName string
	remoteBranchName    string
	isRemoteBranch      bool
}

func branchCheckoutError(err error) string {
	logger.Log(err.Error(), global.StatusError)
	return global.BranchCheckoutError
}

func checkErrorAndPanic(err error) bool {
	if err != nil {
		logger.Log(err.Error(), global.StatusError)
		panic(err)
		return true
	}
	return false
}

// CheckoutBranch checks out the branchName received as argument
func (b branchCheckout) CheckoutBranch() string {
	var errStatus bool
	branchName := b.BranchName

	defer func() string {
		if r := recover(); r != nil {
			logger.Log(fmt.Sprintf("%v", r), global.StatusError)
			return global.BranchCheckoutError
		}
		return global.BranchCheckoutError
	}()

	branchDetails := b.validateAndSetBranchDetails(branchName)

	if branchDetails.isRemoteBranch {
		errStatus = b.checkoutRemoteBranch(branchDetails)
	}

	errStatus = b.checkoutLocalBranch(branchDetails)

	if errStatus {
		return global.BranchCheckoutError
	} else {
		logger.Log(fmt.Sprintf("Current branch checked out to -> %s", branchName), global.StatusInfo)
		return fmt.Sprintf("Head checked out to branch - %v", branchName)
	}
}

func (b branchCheckout) checkoutRemoteBranch(branchDetails validatedBranchDetails) bool {
	var errStatus bool
	branchName := b.BranchName
	repo := b.Repo

	logger.Log(fmt.Sprintf("Branch - %s is a remote branch. Trying with intermediate remote fetch!", branchName), global.StatusWarning)

	remoteBranch, remoteBranchErr := repo.LookupBranch(branchDetails.remoteBranchName, git2go.BranchRemote)
	errStatus = checkErrorAndPanic(remoteBranchErr)

	remoteHead := remoteBranch.Target()

	remoteCommit, remoteCommitErr := repo.LookupCommit(remoteHead)
	errStatus = checkErrorAndPanic(remoteCommitErr)

	remoteTree, remoteTreeErr := remoteCommit.Tree()
	errStatus = checkErrorAndPanic(remoteTreeErr)

	checkoutErr := repo.CheckoutTree(remoteTree, &git2go.CheckoutOptions{Strategy: git2go.CheckoutSafe})
	errStatus = checkErrorAndPanic(checkoutErr)

	_, localLookupErr := repo.LookupBranch(branchName, git2go.BranchLocal)
	if localLookupErr != nil {
		logger.Log(localLookupErr.Error(), global.StatusError)
		if err := addAndCheckoutNewBranch(repo, branchName, remoteCommit, branchDetails); err != nil {
			return true
		}
		return errStatus
	}

	if err := repo.SetHead(branchDetails.referenceBranchName); err != nil {
		branchCheckoutError(err)
		return errStatus
	}

	return errStatus
}

func (b branchCheckout) checkoutLocalBranch(branchDetails validatedBranchDetails) bool {
	errStatus := false
	repo := b.Repo
	branchName := b.BranchName

	branch, branchErr := repo.LookupBranch(branchName, git2go.BranchLocal)
	errStatus = checkErrorAndPanic(branchErr)

	topCommit, _ := repo.LookupCommit(branch.Target())
	if topCommit == nil {
		return true
	}

	tree, treeErr := topCommit.Tree()
	errStatus = checkErrorAndPanic(treeErr)

	checkoutErr := repo.CheckoutTree(tree, &git2go.CheckoutOptions{
		Strategy:       git2go.CheckoutSafe,
		DisableFilters: false,
	})

	if checkoutErr != nil {
		branchCheckoutError(checkoutErr)
		return true
	}

	err := repo.SetHead(branchDetails.referenceBranchName)
	if err != nil {
		branchCheckoutError(err)
		return true
	}
	return errStatus
}

func addAndCheckoutNewBranch(repo *git2go.Repository, branchName string, remoteCommit *git2go.Commit, branchDetails validatedBranchDetails) error {
	var addBranch AddBranch
	addBranch = NewAddBranch(repo, branchName, false, remoteCommit)

	branchCreateStatus := addBranch.AddBranch()
	if branchCreateStatus == global.BranchAddError {
		err := errors.New("branch creation failed")
		branchCheckoutError(err)
		return err
	}

	if err := repo.SetHead(branchDetails.referenceBranchName); err != nil {
		branchCheckoutError(err)
		return err
	}

	return nil
}

func (b branchCheckout) validateAndSetBranchDetails(branchName string) validatedBranchDetails {
	var (
		referenceBranchName string
		remoteBranchName    string
		isRemoteBranch      bool
	)

	if strings.Contains(branchName, "remotes/") {
		referenceBranchName = "refs/heads/" + splitAndReturnBranchName("/", branchName)
		remoteBranchName = splitAndReturnBranchName("remotes/", branchName)
		branchName = splitAndReturnBranchName("/", branchName)
		isRemoteBranch = true
	} else {
		referenceBranchName = "refs/heads/" + branchName
	}

	return validatedBranchDetails{
		branchName:          branchName,
		referenceBranchName: referenceBranchName,
		remoteBranchName:    remoteBranchName,
		isRemoteBranch:      isRemoteBranch,
	}
}

func splitAndReturnBranchName(delimiter string, branchName string) string {
	splitString := strings.Split(branchName, delimiter)
	return splitString[len(splitString)-1]
}

func NewBranchCheckout(repo *git2go.Repository, branchName string) BranchCheckout {
	return branchCheckout{
		Repo:       repo,
		BranchName: branchName,
	}
}
