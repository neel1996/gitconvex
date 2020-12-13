package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/neel1996/gitconvex-server/global"
	"strings"
)

// intermediateFetch performs a remote fetch if the selected checkout branch is a remote branch
func intermediateFetch(repo *git.Repository, branchName string) {
	logger.Log("Fetching from remote for remote branch -> "+branchName, global.StatusInfo)
	remoteChan := make(chan RemoteDataModel)
	go RemoteData(repo, remoteChan)
	remoteData := <-remoteChan
	remoteURL := remoteData.RemoteURL
	FetchFromRemote(repo, *remoteURL[0], branchName)
}

// CheckoutBranch checks out the branchName received as argument
func CheckoutBranch(repo *git.Repository, branchName string) string {
	var isRemoteBranch bool
	var referenceBranchName string

	logger := global.Logger{}
	w, _ := repo.Worktree()

	if strings.Contains(branchName, "remotes/") {
		splitRef := strings.Split(branchName, "/")
		branchName = splitRef[len(splitRef)-1]
		referenceBranchName = "refs/heads/" + branchName
		isRemoteBranch = true
	} else {
		referenceBranchName = "refs/heads/" + branchName
	}

	// If the branch is a remote branch then a remote fetch will be performed and then the branch checkout will be initiated
	if isRemoteBranch {
		logger.Log(fmt.Sprintf("Branch - %s is a remote branch\nTrying with intermediate remote fetch!", branchName), global.StatusWarning)
		intermediateFetch(repo, branchName)

		checkoutErr := w.Checkout(&git.CheckoutOptions{
			Branch: plumbing.ReferenceName(referenceBranchName),
			Force:  true,
		})
		if checkoutErr != nil {
			logger.Log("Checkout failed - "+checkoutErr.Error(), global.StatusError)
			return global.BranchCheckoutError
		}
	}

	checkoutErr := w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(referenceBranchName),
		Keep:   true,
	})
	if checkoutErr != nil {
		logger.Log(checkoutErr.Error(), global.StatusError)
		return global.BranchCheckoutError
	}

	logger.Log(fmt.Sprintf("Current branch checked out to -> %s", branchName), global.StatusInfo)
	return fmt.Sprintf("Head checked out to branch - %v", branchName)
}
