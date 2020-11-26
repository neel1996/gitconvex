package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/neel1996/gitconvex-server/global"
	"strings"
)

// CheckoutBranch checks out the branchName received as argument
func CheckoutBranch(repo *git.Repository, branchName string) string {
	var isRemoteBranch bool
	logger := global.Logger{}
	w, _ := repo.Worktree()

	if strings.Contains(branchName, "remotes/") {
		splitRef := strings.Split(branchName, "/")
		branchName = "refs/heads/" + splitRef[len(splitRef)-1]
		isRemoteBranch = true
	} else {
		branchName = "refs/heads/" + branchName
	}

	// If the branch is a remote branch then a remote fetch will be performed and then the branch checkout will be initiated
	if isRemoteBranch {
		logger.Log(fmt.Sprintf("Branch - %s is a remote branch\nTrying with intermediate remote fetch!", branchName), global.StatusWarning)
		fetchErr := repo.Fetch(&git.FetchOptions{
			RefSpecs: []config.RefSpec{config.RefSpec(branchName + ":" + branchName)},
		})
		if fetchErr != nil {
			logger.Log("Remote fetch failed -> "+fetchErr.Error(), global.StatusWarning)
		}

		checkoutErr := w.Checkout(&git.CheckoutOptions{
			Branch: plumbing.ReferenceName(branchName),
			Force:  true,
		})
		if checkoutErr != nil {
			logger.Log(checkoutErr.Error(), global.StatusError)
			return "CHECKOUT_FAILED"
		}
	}

	checkoutErr := w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(branchName),
		Keep:   true,
	})
	if checkoutErr != nil {
		logger.Log(checkoutErr.Error(), global.StatusError)
		return "CHECKOUT_FAILED"
	}

	logger.Log(fmt.Sprintf("Current branch checked out to -> %s", branchName), global.StatusInfo)
	return fmt.Sprintf("Head checked out to branch - %v", branchName)
}
