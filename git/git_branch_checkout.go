package git

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/global"
	"go/types"
	"strings"
)

type BranchCheckoutInterface interface {
	CheckoutBranch() string
}

type BranchCheckoutInputs struct {
	Repo       *git2go.Repository
	BranchName string
}

func returnCheckoutError(err error) string {
	logger.Log(err.Error(), global.StatusError)
	return global.BranchCheckoutError
}

func checkCheckoutError(err error) bool {
	if err != nil {
		logger.Log(err.Error(), global.StatusError)
		panic(err)
		return true
	}
	return false
}

// CheckoutBranch checks out the branchName received as argument
func (inputs BranchCheckoutInputs) CheckoutBranch() string {
	var isRemoteBranch bool
	var referenceBranchName string
	var remoteBranchName string
	var errStatus bool

	defer func() string {
		if r := recover(); r != nil {
			logger.Log(fmt.Sprintf("%v", r), global.StatusError)
			return global.BranchCheckoutError
		}
		return global.BranchCheckoutError
	}()

	repo := inputs.Repo
	branchName := inputs.BranchName

	if strings.Contains(branchName, "remotes/") {
		singleBranchName := strings.Split(branchName, "/")
		referenceBranchName = "refs/heads/" + singleBranchName[len(singleBranchName)-1]
		remoteBranchSplit := strings.Split(branchName, "remotes/")
		remoteBranchName = remoteBranchSplit[len(remoteBranchSplit)-1]
		branchName = singleBranchName[len(singleBranchName)-1]
		isRemoteBranch = true
	} else {
		referenceBranchName = "refs/heads/" + branchName
	}

	if isRemoteBranch {
		logger.Log(fmt.Sprintf("Branch - %s is a remote branch. Trying with intermediate remote fetch!", branchName), global.StatusWarning)

		remoteBranch, remoteBranchErr := repo.LookupBranch(remoteBranchName, git2go.BranchRemote)
		errStatus = checkCheckoutError(remoteBranchErr)

		remoteHead := remoteBranch.Target()

		remoteCommit, remoteCommitErr := repo.LookupCommit(remoteHead)
		errStatus = checkCheckoutError(remoteCommitErr)

		remoteTree, remoteTreeErr := remoteCommit.Tree()
		errStatus = checkCheckoutError(remoteTreeErr)

		checkoutErr := repo.CheckoutTree(remoteTree, &git2go.CheckoutOptions{Strategy: git2go.CheckoutSafe})
		errStatus = checkCheckoutError(checkoutErr)

		_, localLookupErr := repo.LookupBranch(branchName, git2go.BranchLocal)
		if localLookupErr != nil {
			logger.Log(localLookupErr.Error(), global.StatusError)

			var addBranch AddBranch
			addBranch = NewAddBranch(repo, branchName, false, remoteCommit)

			branchCreateStatus := addBranch.AddBranch()
			if branchCreateStatus != global.BranchAddError {
				err := repo.SetHead(referenceBranchName)
				if err != nil {
					return returnCheckoutError(err)
				} else {
					return fmt.Sprintf("Head checked out to branch - %v", branchName)
				}
			} else {
				return returnCheckoutError(types.Error{Msg: "Branch creation failed"})
			}
		} else {
			err := repo.SetHead(referenceBranchName)
			if err != nil {
				return returnCheckoutError(err)
			} else {
				return fmt.Sprintf("Head checked out to branch - %v", branchName)
			}
		}
	}

	branch, branchErr := repo.LookupBranch(branchName, git2go.BranchLocal)
	errStatus = checkCheckoutError(branchErr)

	topCommit, _ := repo.LookupCommit(branch.Target())
	if topCommit != nil {
		tree, treeErr := topCommit.Tree()
		errStatus = checkCheckoutError(treeErr)

		checkoutErr := repo.CheckoutTree(tree, &git2go.CheckoutOptions{
			Strategy:       git2go.CheckoutSafe,
			DisableFilters: false,
		})

		if checkoutErr != nil {
			return returnCheckoutError(checkoutErr)
		}

		err := repo.SetHead(referenceBranchName)
		if err != nil {
			return returnCheckoutError(err)
		}
	}

	if errStatus {
		return global.BranchCheckoutError
	} else {
		logger.Log(fmt.Sprintf("Current branch checked out to -> %s", branchName), global.StatusInfo)
		return fmt.Sprintf("Head checked out to branch - %v", branchName)
	}
}
