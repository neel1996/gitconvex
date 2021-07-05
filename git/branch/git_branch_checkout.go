package branch

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/global"
	"strings"
)

type Checkout interface {
	CheckoutBranch() error
}

type branchCheckout struct {
	repo       *git2go.Repository
	branchName string
}

type validatedBranchDetails struct {
	branchName          string
	referenceBranchName string
	remoteBranchName    string
	isRemoteBranch      bool
}

func printAndReturnError(err error) error {
	logger.Log(err.Error(), global.StatusError)
	return err
}

// CheckoutBranch checks out the branchName received as argument
func (b branchCheckout) CheckoutBranch() error {
	err := b.validateBranchFields()
	if err != nil {
		return err
	}

	branchDetails := b.validateAndSetBranchDetails()

	if branchDetails.isRemoteBranch {
		return b.checkoutRemoteBranch(branchDetails)
	}

	return b.checkoutLocalBranch(branchDetails)
}

func (b branchCheckout) validateBranchFields() error {
	if err := NewBranchFieldsValidation(b.repo, b.branchName).ValidateBranchFields(); err != nil {
		logger.Log(err.Error(), global.StatusError)
		return err
	}
	return nil
}

func (b branchCheckout) checkoutRemoteBranch(branchDetails validatedBranchDetails) error {
	branchName := b.branchName
	repo := b.repo

	logger.Log(fmt.Sprintf("Branch - %s is a remote branch. Trying with intermediate remote fetch!", branchName), global.StatusWarning)

	remoteBranch, remoteBranchErr := repo.LookupBranch(branchDetails.remoteBranchName, git2go.BranchRemote)
	if remoteBranchErr != nil {
		return printAndReturnError(remoteBranchErr)
	}

	remoteHead := remoteBranch.Target()
	remoteCommit, remoteCommitErr := repo.LookupCommit(remoteHead)
	if remoteCommitErr != nil {
		return printAndReturnError(remoteCommitErr)

	}

	remoteTree, remoteTreeErr := remoteCommit.Tree()
	if remoteTreeErr != nil {
		return printAndReturnError(remoteTreeErr)
	}

	if checkoutErr := repo.CheckoutTree(remoteTree, &git2go.CheckoutOptions{Strategy: git2go.CheckoutSafe}); checkoutErr != nil {
		return printAndReturnError(checkoutErr)
	}

	_, localLookupErr := repo.LookupBranch(branchName, git2go.BranchLocal)
	if localLookupErr != nil {
		logger.Log(localLookupErr.Error(), global.StatusError)
		return addAndCheckoutNewBranch(repo, branchName, remoteCommit, branchDetails)
	}

	if err := repo.SetHead(branchDetails.referenceBranchName); err != nil {
		return printAndReturnError(err)
	}

	logger.Log(fmt.Sprintf("Remote branch %v has been checked out", branchDetails.remoteBranchName), global.StatusInfo)
	return nil
}

func (b branchCheckout) checkoutLocalBranch(branchDetails validatedBranchDetails) error {
	repo := b.repo
	branchName := b.branchName

	branch, branchErr := repo.LookupBranch(branchName, git2go.BranchLocal)
	if branchErr != nil {
		fmt.Println("Lookup error")
		return printAndReturnError(branchErr)
	}

	lookupCommit, lookupCommitErr := repo.LookupCommit(branch.Target())
	if lookupCommit == nil {
		fmt.Println("Lookup commit error")
		return printAndReturnError(lookupCommitErr)
	}

	tree, treeErr := lookupCommit.Tree()
	if treeErr != nil {
		fmt.Println("Lookup tree error")
		return printAndReturnError(treeErr)
	}

	checkoutErr := repo.CheckoutTree(tree, &git2go.CheckoutOptions{
		Strategy:       git2go.CheckoutSafe,
		DisableFilters: false,
	})

	if checkoutErr != nil {
		return printAndReturnError(checkoutErr)
	}

	err := repo.SetHead(branchDetails.referenceBranchName)
	if err != nil {
		return printAndReturnError(err)
	}

	logger.Log(fmt.Sprintf("Local branch %v has been checked out", b.branchName), global.StatusInfo)
	return nil
}

func addAndCheckoutNewBranch(repo *git2go.Repository, branchName string, remoteCommit *git2go.Commit, branchDetails validatedBranchDetails) error {
	addBranch := NewAddBranch(repo, branchName, false, remoteCommit)

	branchAddError := addBranch.AddBranch()
	if branchAddError != nil {
		return branchAddError
	}

	if err := repo.SetHead(branchDetails.referenceBranchName); err != nil {
		return err
	}

	return nil
}

func (b branchCheckout) validateAndSetBranchDetails() validatedBranchDetails {
	var (
		branchName          = b.branchName
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

func NewBranchCheckout(repo *git2go.Repository, branchName string) Checkout {
	return branchCheckout{
		repo:       repo,
		branchName: branchName,
	}
}
