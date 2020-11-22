package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/neel1996/gitconvex-server/global"
	"strings"
)

// CheckoutBranch checks out the branchName received as argument
func CheckoutBranch(repo *git.Repository, branchName string) string {
	logger := global.Logger{}

	if strings.Contains(branchName, "remotes/") {
		splitRef := strings.Split(branchName, "/")
		branchName = "refs/heads/" + splitRef[len(splitRef)-1]
	} else {
		branchName = "refs/heads/" + branchName
	}

	w, _ := repo.Worktree()
	checkoutErr := w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(branchName),
		Keep:   true,
	})
	if checkoutErr != nil {
		logger.Log(fmt.Sprintf("Failed to checkout branch - %s --> %v\nRetrying with new branch creation!", branchName, checkoutErr.Error()), global.StatusWarning)
		err := w.Checkout(&git.CheckoutOptions{
			Branch: plumbing.ReferenceName(branchName),
			Create: true,
			Keep:   true,
		})
		if err != nil {
			return "CHECKOUT_FAILED"
		}
	}
	logger.Log(fmt.Sprintf("Current branch checked out to -> %s", branchName), global.StatusInfo)
	return fmt.Sprintf("Head checked out to branch - %v", branchName)
}
